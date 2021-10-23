package outgoingrpclogger

import (
	"context"
	"unsafe"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/rs/zerolog"
)

type OptionsBuilder func(options *options)

type options struct {
	MaxRawParamsSize int
	MaxRawRespSize   int
}

func (o *options) Init() *options {
	o.MaxRawParamsSize = 500
	o.MaxRawRespSize = 500
	return o
}

func MaxRawParamsSize(value int) OptionsBuilder {
	return func(options *options) { options.MaxRawParamsSize = value }
}

func MaxRawRespSize(value int) OptionsBuilder {
	return func(options *options) { options.MaxRawRespSize = value }
}

func NewForClient(logger zerolog.Logger, optionsBuilders ...OptionsBuilder) apicommon.RPCHandler {
	options := new(options).Init()
	for _, optionsBuilder := range optionsBuilders {
		optionsBuilder(options)
	}
	return func(ctx context.Context, rpc *apicommon.RPC) (returnedErr error) {
		outgoingRPC := rpc.OutgoingRPC()
		subLogger := logger
		ctx = subLogger.WithContext(ctx)
		// Before
		returnedErr = outgoingRPC.Do(ctx)
		// After
		var event *zerolog.Event
		if !outgoingRPC.IsRequested() && returnedErr != nil {
			event = subLogger.Error()
		} else {
			event = subLogger.Info()
		}
		if traceID := outgoingRPC.TraceID(); traceID != "" {
			event.Str("traceID", outgoingRPC.TraceID())
		}
		event.Str("fullMethodName", outgoingRPC.FullMethodName())
		event.Str("url", outgoingRPC.URL())
		if rawParams := outgoingRPC.RawParams(); rawParams != nil {
			if apicommon.DebugMode || len(rawParams) <= options.MaxRawParamsSize {
				event.Str("params", bytesToString(rawParams))
			} else {
				event.Int("paramsSize", len(rawParams))
				event.Str("truncatedParams", bytesToString(rawParams[:options.MaxRawParamsSize]))
			}
		}
		event.Bool("isRequested", outgoingRPC.IsRequested())
		if outgoingRPC.IsRequested() {
			if statusCode := outgoingRPC.StatusCode(); statusCode != 0 {
				event.Int("statusCode", statusCode)
				if rawResp := outgoingRPC.RawResp(); rawResp != nil {
					if apicommon.DebugMode || len(rawResp) <= options.MaxRawRespSize {
						event.Str("resp", bytesToString(rawResp))
					} else {
						event.Int("respSize", len(rawResp))
						event.Str("truncatedResp", bytesToString(rawResp[:options.MaxRawRespSize]))
					}
					if errorCode := outgoingRPC.Error().Code; errorCode != 0 {
						event.Int("errorCode", int(errorCode))
					}
				}
			}
		}
		if returnedErr != nil {
			if _, ok := returnedErr.(*apicommon.Error); !ok {
				event.AnErr("err", returnedErr)
			}
		}
		event.Msg("outgoing rpc")
		return
	}
}

func bytesToString(bytes []byte) string { return *(*string)(unsafe.Pointer(&bytes)) }
