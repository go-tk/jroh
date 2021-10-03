package incomingrpclogger

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

func NewForServer(logger zerolog.Logger, optionsSetters ...OptionsSetter) apicommon.ServerMiddleware {
	options := new(options).Init()
	for _, optionsSetter := range optionsSetters {
		optionsSetter(options)
	}
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			incomingRPC := apicommon.MustGetRPCFromContext(ctx).IncomingRPC()
			subLogger := logger.With().Str("traceID", incomingRPC.TraceID()).Logger()
			ctx = subLogger.WithContext(ctx)
			r = r.WithContext(ctx)
			// Before
			handler.ServeHTTP(w, r)
			// After
			var event *zerolog.Event
			if incomingRPC.StatusCode()/100 == 5 || incomingRPC.InternalErr() != nil {
				event = subLogger.Error()
			} else {
				event = subLogger.Info()
			}
			event.Str("fullMethodName", incomingRPC.FullMethodName())
			if remoteAddr := r.RemoteAddr; remoteAddr != "" {
				event.Str("remoteAddr", remoteAddr)
			}
			event.Str("rpcPath", r.URL.Path)
			if rawParams := incomingRPC.RawParams(); rawParams != nil {
				if apicommon.DebugMode || len(rawParams) <= options.MaxRawParamsSize {
					event.Str("params", bytesToString(rawParams))
				} else {
					event.Int("paramsSize", len(rawParams))
					event.Str("truncatedParams", bytesToString(rawParams[:options.MaxRawParamsSize]))
				}
			}
			event.Int("statusCode", incomingRPC.StatusCode())
			if errorCode := incomingRPC.Error().Code; errorCode != 0 {
				event.Int("errorCode", int(errorCode))
			}
			if internalErr := incomingRPC.InternalErr(); internalErr != nil {
				event.AnErr("internalErr", internalErr)
				if stackTrace := incomingRPC.StackTrace(); stackTrace != "" {
					event.Str("stackTrace", stackTrace)
				}
			}
			if rawResp := incomingRPC.RawResp(); rawResp != nil {
				if apicommon.DebugMode || len(rawResp) <= options.MaxRawRespSize {
					event.Str("resp", bytesToString(rawResp))
				} else {
					event.Int("respSize", len(rawResp))
					event.Str("truncatedResp", bytesToString(rawResp[:options.MaxRawRespSize]))
				}
			}
			if responseWriteErr := incomingRPC.ResponseWriteErr(); responseWriteErr != nil {
				event.AnErr("responseWriteErr", responseWriteErr)
			}
			event.Msg("incoming rpc")
		})
	}
}

func bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
