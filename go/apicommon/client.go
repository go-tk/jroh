package apicommon

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

type ClientOptions struct {
	RPCFilters  map[MethodIndex][]RPCHandler
	Middlewares map[MethodIndex][]ClientMiddleware
	Transport   http.RoundTripper
	Timeout     time.Duration
}

func (co *ClientOptions) Sanitize() {
	transport := co.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	co.Transport = TransportFunc(func(request *http.Request) (*http.Response, error) {
		response, err := transport.RoundTrip(request)
		if err != nil {
			return nil, err
		}
		var buffer bytes.Buffer
		_, err = buffer.ReadFrom(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, err
		}
		if buffer.Len() >= 1 {
			outgoingRPC := MustGetRPCFromContext(request.Context()).OutgoingRPC()
			outgoingRPC.rawResp = buffer.Bytes()
		}
		return response, nil
	})
}

type ClientMiddleware func(oldTransport http.RoundTripper) (newTransport http.RoundTripper)

type TransportFunc func(request *http.Request) (response *http.Response, err error)

var _ http.RoundTripper = TransportFunc(nil)

func (tf TransportFunc) RoundTrip(request *http.Request) (*http.Response, error) { return tf(request) }

type Client struct {
	c          http.Client
	rpcBaseURL string
}

func (c *Client) Init(rpcBaseURL string, timeout time.Duration) {
	c.rpcBaseURL = rpcBaseURL
	c.c.Transport = TransportFunc(func(request *http.Request) (*http.Response, error) {
		outgoingRPC := MustGetRPCFromContext(request.Context()).OutgoingRPC()
		return outgoingRPC.transport.RoundTrip(request)
	})
	c.c.Timeout = timeout
}

func (c *Client) DoRPC(ctx context.Context, outgoingRPC *OutgoingRPC, transport http.RoundTripper, rpcPath string) error {
	if rpc, ok := GetRPCFromContext(ctx); ok {
		outgoingRPC.traceID = rpc.traceID
		outgoingRPC.traceIDIsReceived = true
	}
	outgoingRPC.client = &c.c
	outgoingRPC.transport = transport
	outgoingRPC.url = c.rpcBaseURL + rpcPath
	return outgoingRPC.Do(makeContextWithRPC(ctx, &outgoingRPC.RPC))
}

type UnexpectedStatusCodeError struct {
	StatusCode int

	namespace   string
	serviceName string
	methodName  string
	traceID     string
}

func (usce *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("apicommon: unexpected status code; namespace=%q serviceName=%q methodName=%q traceID=%q statusCode=%q",
		usce.namespace, usce.serviceName, usce.methodName, usce.traceID, usce.StatusCode)
}

func HandleRPC(ctx context.Context, rpc *RPC) error {
	outgoingRPC := rpc.OutgoingRPC()
	if err := outgoingRPC.encodeParams(); err != nil {
		return fmt.Errorf("apicommon: params encoding failed; namespace=%q serviceName=%q methodName=%q traceID=%q: %v",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, err)
	}
	response, err := outgoingRPC.requestHTTP(ctx)
	if err != nil {
		return fmt.Errorf("apicommon: http request failed; namespace=%q serviceName=%q methodName=%q traceID=%q: %v",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, err)
	}
	outgoingRPC.statusCode = response.StatusCode
	if outgoingRPC.statusCode != http.StatusOK {
		return &UnexpectedStatusCodeError{
			StatusCode: outgoingRPC.statusCode,

			namespace:   rpc.namespace,
			serviceName: rpc.serviceName,
			methodName:  rpc.methodName,
			traceID:     rpc.traceID,
		}
	}
	if err := outgoingRPC.decodeResp(ctx); err != nil {
		return fmt.Errorf("apicommon: resp decoding failed; namespace=%q serviceName=%q methodName=%q traceID=%q: %v",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, err)
	}
	if outgoingRPC.error.Code != 0 {
		return fmt.Errorf("apicommon: rpc failed; namespace=%q serviceName=%q methodName=%q traceID=%q: %w",
			rpc.namespace, rpc.serviceName, rpc.methodName, rpc.traceID, &outgoingRPC.error)
	}
	return nil
}

var _ RPCHandler = HandleRPC
