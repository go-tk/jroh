package apicommon

import (
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

func HandleRPC(ctx context.Context, rpc *RPC) error {
	outgoingRPC := rpc.OutgoingRPC()
	if err := outgoingRPC.encodeParams(); err != nil {
		return err
	}
	if err := outgoingRPC.requestHTTP(ctx); err != nil {
		return err
	}
	if err := outgoingRPC.decodeResp(ctx); err != nil {
		return err
	}
	if outgoingRPC.error.Code != 0 {
		return &outgoingRPC.error
	}
	return nil
}

var _ RPCHandler = HandleRPC
