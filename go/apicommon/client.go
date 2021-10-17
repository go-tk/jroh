package apicommon

import (
	"bytes"
	"context"
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
	if co.Transport == nil {
		co.Transport = http.DefaultTransport
	}
}

type ClientMiddleware func(oldTransport http.RoundTripper) (newTransport http.RoundTripper)

func FillTransportTable(transportTable []http.RoundTripper, transport http.RoundTripper, clientMiddlewares map[MethodIndex][]ClientMiddleware) {
	transport = func(transport http.RoundTripper) http.RoundTripper {
		return TransportFunc(func(request *http.Request) (*http.Response, error) {
			outgoingRPC := MustGetRPCFromContext(request.Context()).OutgoingRPC()
			outgoingRPC.isRequested = true
			response, err := transport.RoundTrip(request)
			if err != nil {
				return nil, err
			}
			outgoingRPC.statusCode = response.StatusCode
			if response.StatusCode != http.StatusOK {
				response.Body.Close()
				return response, nil
			}
			var buffer bytes.Buffer
			_, err = buffer.ReadFrom(response.Body)
			response.Body.Close()
			if err != nil {
				return nil, err
			}
			if buffer.Len() >= 1 {
				outgoingRPC.rawResp = buffer.Bytes()
			}
			return response, nil
		})
	}(transport)
	for i := range transportTable {
		transport := transport
		methodIndex := MethodIndex(i)
		clientMiddlewares2 := clientMiddlewares[methodIndex]
		for i := len(clientMiddlewares2) - 1; i >= 0; i-- {
			clientMiddleware := clientMiddlewares2[i]
			transport = clientMiddleware(transport)
		}
		clientMiddlewares3 := clientMiddlewares[AnyMethod]
		for i := len(clientMiddlewares3) - 1; i >= 0; i-- {
			clientMiddleware := clientMiddlewares3[i]
			transport = clientMiddleware(transport)
		}
		transportTable[methodIndex] = transport
	}
}

type Client struct {
	c          http.Client
	rpcBaseURL string
}

func (c *Client) Init(timeout time.Duration, rpcBaseURL string) {
	c.c.Transport = TransportFunc(func(request *http.Request) (*http.Response, error) {
		outgoingRPC := MustGetRPCFromContext(request.Context()).OutgoingRPC()
		return outgoingRPC.transport.RoundTrip(request)
	})
	c.c.Timeout = timeout
	c.rpcBaseURL = rpcBaseURL
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

type TransportFunc func(request *http.Request) (response *http.Response, err error)

var _ http.RoundTripper = TransportFunc(nil)

func (tf TransportFunc) RoundTrip(request *http.Request) (*http.Response, error) { return tf(request) }
