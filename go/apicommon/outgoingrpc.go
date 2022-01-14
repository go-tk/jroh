package apicommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OutgoingRPC struct {
	Namespace      string
	ServiceName    string
	MethodName     string
	FullMethodName string
	MethodIndex    int
	Params         Model
	Results        Model

	Transport      http.RoundTripper
	URL            string
	OutboundHeader http.Header
	RawParams      []byte
	StatusCode     int
	InboundHeader  http.Header
	TraceID        string
	ErrorCode      ErrorCode
	RawResults     []byte

	handler         OutgoingRPCHandler
	filters         []OutgoingRPCHandler
	nextFilterIndex int
	readCloser      io.ReadCloser

	encodeParamsCache struct {
		err error
		has bool
	}

	readRawResultsCache struct {
		err error
		has bool
	}

	loadResultsCache struct {
		err error
		has bool
	}
}

func (or *OutgoingRPC) SetHandler(handler OutgoingRPCHandler)   { or.handler = handler }
func (or *OutgoingRPC) SetFilters(filters []OutgoingRPCHandler) { or.filters = filters }
func (or *OutgoingRPC) SetReadCloser(readCloser io.ReadCloser)  { or.readCloser = readCloser }

func (or *OutgoingRPC) Do(ctx context.Context) error {
	i := or.nextFilterIndex
	or.nextFilterIndex++
	n := len(or.filters)
	if i > n {
		panic("too many calls")
	}
	if i == 0 {
		or.OutboundHeader = outboundHeaderPool.Get().(http.Header)
		defer func() {
			if or.readCloser != nil {
				or.readCloser.Close()
			}
		}()
	}
	var err error
	if i < n {
		err = or.filters[i](ctx, or)
	} else {
		err = or.handler(ctx, or)
	}
	if i >= 1 || err != nil {
		return err
	}
	if err := or.LoadResults(ctx); err != nil {
		return fmt.Errorf("load results: %w", err)
	}
	return nil
}

func (or *OutgoingRPC) EncodeParams() error {
	cache := &or.encodeParamsCache
	if !cache.has {
		cache.err = or.doEncodeParams()
		cache.has = true
	}
	return cache.err
}

func (or *OutgoingRPC) doEncodeParams() error {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if DebugMode {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(or.Params); err != nil {
		return fmt.Errorf("marshal json: objectType=\"%T\": %w", or.Params, err)
	}
	rawParams := buffer.Bytes()
	if !DebugMode {
		// Remove trailing '\n'
		rawParams = rawParams[:len(rawParams)-1]
	}
	or.RawParams = rawParams
	return nil
}

func (or *OutgoingRPC) LoadResults(ctx context.Context) error {
	cache := &or.loadResultsCache
	if !cache.has {
		cache.err = or.doLoadResults(ctx)
		cache.has = true
	}
	return cache.err
}

func (or *OutgoingRPC) doLoadResults(ctx context.Context) error {
	if err := or.ReadRawResults(); err != nil {
		return err
	}
	if err := or.decodeRawResults(ctx); err != nil {
		return err
	}
	validationContext := NewValidationContext(ctx)
	if !or.Results.Validate(validationContext) {
		return fmt.Errorf("%w: %v", ErrInvalidResults, validationContext.ErrorDetails())
	}
	return nil
}

func (or *OutgoingRPC) decodeRawResults(ctx context.Context) error {
	if err := json.Unmarshal(or.RawResults, or.Results); err != nil {
		return fmt.Errorf("unmarshal json; objectType=\"%T\": %w", or.Results, err)
	}
	return nil
}

func (or *OutgoingRPC) ReadRawResults() error {
	cache := &or.readRawResultsCache
	if !cache.has {
		cache.err = or.doReadRawResults()
		cache.has = true
	}
	return cache.err
}

func (or *OutgoingRPC) doReadRawResults() error {
	var buffer bytes.Buffer
	n, err := buffer.ReadFrom(or.readCloser)
	if err != nil {
		return fmt.Errorf("read data: %w", err)
	}
	if n == 0 {
		return nil
	}
	or.RawResults = buffer.Bytes()
	return nil
}

type OutgoingRPCHandler func(ctx context.Context, outgoingRPC *OutgoingRPC) (err error)
