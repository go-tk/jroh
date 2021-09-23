package outgoingrpclogger

import (
	"net/http"
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
	o.MaxRawParamsSize = 512
	o.MaxRawRespSize = 512
	return o
}

func MaxRawParamsSize(value int) OptionsSetter {
	return func(options *options) { options.MaxRawParamsSize = value }
}

func MaxRawRespSize(value int) OptionsSetter {
	return func(options *options) { options.MaxRawRespSize = value }
}

func New(logger zerolog.Logger, optionsSetters ...OptionsSetter) apicommon.ClientMiddleware {
	options := new(options).Init()
	for _, optionsSetter := range optionsSetters {
		optionsSetter(options)
	}
	return func(transport http.RoundTripper) http.RoundTripper {
		return apicommon.TransportFunc(func(request *http.Request) (*http.Response, error) {
			// Before
			response, err := transport.RoundTrip(request)
			if err != nil {
				return nil, err
			}
			// After
			event := logger.Info()
			outgoingRPC := apicommon.MustGetRPCFromContext(request.Context()).OutgoingRPC()
			if traceID := outgoingRPC.TraceID(); traceID != "" {
				event.Str("traceID", outgoingRPC.TraceID())
			}
			event.Str("url", outgoingRPC.URL())
			if rawParams := outgoingRPC.RawParams(); rawParams != nil {
				if apicommon.DebugMode || len(rawParams) <= options.MaxRawParamsSize {
					event.Str("params", bytesToString(rawParams))
				} else {
					event.Str("truncatedParams", bytesToString(rawParams[:options.MaxRawParamsSize]))
				}
			}
			if statusCode := outgoingRPC.StatusCode(); statusCode != 0 {
				event.Int("statusCode", statusCode)
				if rawResp := outgoingRPC.RawResp(); rawResp != nil {
					if apicommon.DebugMode || len(rawResp) <= options.MaxRawRespSize {
						event.Str("resp", bytesToString(rawResp))
					} else {
						event.Str("truncatedResp", bytesToString(rawResp[:options.MaxRawRespSize]))
					}
				}
			}
			event.Msg("outgoing rpc")
			return response, nil
		})
	}
}

func bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
