package loggerrf

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/rs/zerolog"
)

func NewIncomingRPCFilter(logger zerolog.Logger, optionsBuilders ...OptionsBuilder) apicommon.IncomingRPCHandler {
	options := new(options).Init()
	for _, optionsBuilder := range optionsBuilders {
		optionsBuilder(options)
	}
	return func(ctx context.Context, incomingRPC *apicommon.IncomingRPC) (returnedErr error) {
		subLogger := logger.With().Str("traceID", incomingRPC.TraceID).Logger()
		ctx = subLogger.WithContext(ctx)
		defer func() {
			if v := recover(); v != nil {
				incomingRPC.StatusCode = http.StatusInternalServerError
				incomingRPC.ErrorCode = -1
				buffer := make([]byte, 4096)
				i := copy(buffer, fmt.Sprintf("panic: %v\n\n", v))
				i += runtime.Stack(buffer[i:], false)
				returnedErr = errors.New(apicommon.BytesToString(buffer[:i]))
			}
			if returnedErr == nil {
				if err := incomingRPC.EncodeResults(); err != nil {
					logger.Error().Str("traceID", incomingRPC.TraceID).Msg(err.Error())
				}
			}
			var event *zerolog.Event
			if apicommon.ServerShouldReportError(returnedErr, incomingRPC.StatusCode) {
				event = subLogger.Error()
			} else {
				event = subLogger.Info()
			}
			event.Str("fullMethodName", incomingRPC.FullMethodName)
			event.Int("statusCode", incomingRPC.StatusCode)
			event.Int("errorCode", int(incomingRPC.ErrorCode))
			if returnedErr != nil {
				event.Str("err", returnedErr.Error())
			}
			event.Str("remoteIP", incomingRPC.RemoteIP)
			paramsStr := apicommon.BytesToString(incomingRPC.RawParams)
			if apicommon.DebugMode || len(paramsStr) < options.MaxParamsSize {
				event.Str("params", paramsStr)
			} else {
				event.Int("paramsSize", len(paramsStr))
				event.Str("truncatedParams", paramsStr[:options.MaxParamsSize])
			}
			resultsStr := apicommon.BytesToString(incomingRPC.RawResults)
			if apicommon.DebugMode || len(resultsStr) < options.MaxResultsSize {
				event.Str("results", resultsStr)
			} else {
				event.Int("resultsSize", len(resultsStr))
				event.Str("truncatedResults", resultsStr[:options.MaxResultsSize])
			}
			event.Msg("incoming rpc")
		}()
		return incomingRPC.Do(ctx)
	}
}
