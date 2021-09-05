package apicommon

import (
	"context"
	"net/http"
	"time"
)

type Client struct {
	c          http.Client
	rpcBaseURL string
}

type ClientOptions struct {
	Transport http.RoundTripper
	Timeout   time.Duration

	RPCInterceptors map[MethodIndex][]RPCHandler
}

func (co *ClientOptions) sanitize() {
	if co.Transport == nil {
		co.Transport = http.DefaultTransport
	}
}

func (c *Client) Init(rpcBaseURL string, options ClientOptions) {
	c.rpcBaseURL = rpcBaseURL
	options.sanitize()
	c.c.Transport = options.Transport
	c.c.Timeout = options.Timeout
}

func (c *Client) DoRPC(ctx context.Context, outgoingRPC *OutgoingRPC, rpcPath string) error {
	if rpc, ok := GetRPCFromContext(ctx); ok {
		outgoingRPC.traceID = rpc.traceID
		outgoingRPC.traceIDIsReceived = true
	}
	outgoingRPC.client = &c.c
	outgoingRPC.url = c.rpcBaseURL + rpcPath
	return outgoingRPC.Do(makeContextWithRPC(ctx, &outgoingRPC.RPC))
}
