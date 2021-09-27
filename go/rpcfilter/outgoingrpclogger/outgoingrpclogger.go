package outgoingrpclogger

import (
	"context"
	"unsafe"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/rs/zerolog"
)

type OptionsSetter func(options *options)

type options struct {
	MaxRawParamsSize int
	MaxRawRespSize   int
}

func (o *options) Init() *options {
	o.MaxRawParamsSize = 500
	o.MaxRawRespSize = 500
	return o
}

func MaxRawParamsSize(value int) OptionsSetter {
	return func(options *options) { options.MaxRawParamsSize = value }
}

func MaxRawRespSize(value int) OptionsSetter {
	return func(options *options) { options.MaxRawRespSize = value }
}

func NewForClient(logger zerolog.Logger, optionsSetters ...OptionsSetter) apicommon.RPCHandler {
	options := new(options).Init()
	for _, optionsSetter := range optionsSetters {
		optionsSetter(options)
	}
	return func(ctx context.Context, rpc *apicommon.RPC) error {
		subLogger := logger
		ctx = subLogger.WithContext(ctx)
		// Before
		err := rpc.Do(ctx)
		// After
		event := subLogger.Info()
		outgoingRPC := rpc.OutgoingRPC()
		if traceID := outgoingRPC.TraceID(); traceID != "" {
			event.Str("traceID", outgoingRPC.TraceID())
		}
		event.Str("url", outgoingRPC.URL())
		if rawParams := outgoingRPC.RawParams(); rawParams != nil {
			if apicommon.DebugMode || len(rawParams) <= options.MaxRawParamsSize {
				event.Str("params", bytesToString(rawParams))
			} else {
				event.Int("paramsSize", len(rawParams))
				event.Str("truncatedParams", bytesToString(rawParams[:options.MaxRawParamsSize]))
			}
		}
		if statusCode := outgoingRPC.StatusCode(); statusCode != 0 {
			event.Int("statusCode", statusCode)
			if rawResp := outgoingRPC.RawResp(); rawResp != nil {
				if errorCode := outgoingRPC.Error().Code; errorCode != 0 {
					event.Int("errorCode", int(errorCode))
				}
				if apicommon.DebugMode || len(rawResp) <= options.MaxRawRespSize {
					event.Str("resp", bytesToString(rawResp))
				} else {
					event.Int("respSize", len(rawResp))
					event.Str("truncatedResp", bytesToString(rawResp[:options.MaxRawRespSize]))
				}
			}
		}
		event.Msg("outgoing rpc")
		return err
	}
}

func bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
