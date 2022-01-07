package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IncomingRPC struct {
	Namespace      string
	ServiceName    string
	MethodName     string
	FullMethodName string
	MethodIndex    int
	Params         Model
	Results        Model

	RemoteIP       string
	InboundHeader  http.Header
	OutboundHeader http.Header
	TraceID        string
	RawParams      []byte
	StatusCode     int
	ErrorCode      ErrorCode
	RawResults     []byte

	handler         IncomingRPCHandler
	filters         []IncomingRPCHandler
	nextFilterIndex int
	reader          io.Reader

	readRawParamsCache struct {
		err error
		has bool
	}

	loadParamsCache struct {
		err error
		has bool
	}

	encodeResultsCache struct {
		err error
		has bool
	}
}

func (ir *IncomingRPC) SetHandler(handler IncomingRPCHandler)   { ir.handler = handler }
func (ir *IncomingRPC) SetFilters(filters []IncomingRPCHandler) { ir.filters = filters }
func (ir *IncomingRPC) SetReader(reader io.Reader)              { ir.reader = reader }

func (ir *IncomingRPC) Do(ctx context.Context) (returnedErr error) {
	i := ir.nextFilterIndex
	ir.nextFilterIndex++
	n := len(ir.filters)
	if i > n {
		panic("too many calls")
	}
	if i == 0 {
		ctx = MakeContextWithIncomingRPC(ctx, ir)
	}
	if i < n {
		return ir.filters[i](ctx, ir)
	}
	if err := ir.LoadParams(ctx); err != nil {
		if error, ok := err.(*Error); ok {
			ir.StatusCode = error.StatusCode
			ir.ErrorCode = error.Code
		} else {
			ir.StatusCode = http.StatusBadRequest
			ir.ErrorCode = -1
		}
		return err
	}
	if err := ir.handler(ctx, ir); err != nil {
		if error, ok := err.(*Error); ok {
			ir.StatusCode = error.StatusCode
			ir.ErrorCode = error.Code
		} else {
			if errIsTemporary(err) {
				ir.StatusCode = http.StatusServiceUnavailable
			} else {
				ir.StatusCode = http.StatusInternalServerError
			}
			ir.ErrorCode = -1
		}
		return err
	}
	return nil
}

func (ir *IncomingRPC) LoadParams(ctx context.Context) error {
	cache := &ir.loadParamsCache
	if !cache.has {
		cache.err = ir.doLoadParams(ctx)
		cache.has = true
	}
	return cache.err
}

func (ir *IncomingRPC) doLoadParams(ctx context.Context) error {
	if err := ir.ReadRawParams(); err != nil {
		return err
	}
	if err := ir.decodeRawParams(ctx); err != nil {
		return err
	}
	validationContext := NewValidationContext(ctx)
	if !ir.Params.Validate(validationContext) {
		error := NewInvalidParamsError()
		error.Details = validationContext.ErrorDetails()
		return error
	}
	return nil
}

func (ir *IncomingRPC) decodeRawParams(ctx context.Context) error {
	if err := json.Unmarshal(ir.RawParams, ir.Params); err != nil {
		return fmt.Errorf("raw params decoding failed: %w", err)
	}
	return nil
}

func (ir *IncomingRPC) ReadRawParams() error {
	cache := &ir.readRawParamsCache
	if !cache.has {
		cache.err = ir.doReadRawParams()
		cache.has = true
	}
	return cache.err
}

func (ir *IncomingRPC) doReadRawParams() error {
	var buffer bytes.Buffer
	n, err := buffer.ReadFrom(ir.reader)
	if err != nil {
		return fmt.Errorf("raw params read failed: %w", err)
	}
	if n == 0 {
		return nil
	}
	ir.RawParams = buffer.Bytes()
	return nil
}

func (ir *IncomingRPC) EncodeResults() error {
	cache := &ir.encodeResultsCache
	if !cache.has {
		cache.err = ir.doEncodeResults()
		cache.has = true
	}
	return cache.err
}

func (ir *IncomingRPC) doEncodeResults() error {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if DebugMode {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(ir.Results); err != nil {
		return fmt.Errorf("results encoding failed: %w", err)
	}
	rawResults := buffer.Bytes()
	if !DebugMode {
		// Remove trailing '\n'
		rawResults = rawResults[:len(rawResults)-1]
	}
	ir.RawResults = rawResults
	return nil
}

type IncomingRPCHandler func(ctx context.Context, incomingRPC *IncomingRPC) (err error)

type contextValueIncomingRPC struct{}

func MakeContextWithIncomingRPC(ctx context.Context, incomingRPC *IncomingRPC) context.Context {
	return context.WithValue(ctx, contextValueIncomingRPC{}, incomingRPC)
}

func GetIncomingRPCFromContext(ctx context.Context) (*IncomingRPC, bool) {
	incomingRPC, ok := ctx.Value(contextValueIncomingRPC{}).(*IncomingRPC)
	return incomingRPC, ok
}
