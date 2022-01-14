package zerologmw

import (
	"context"

	"github.com/go-tk/jroh/go/apicommon"
	"github.com/rs/zerolog"
)

func NewOutgoingRPCFilter(logger zerolog.Logger, optionsBuilders ...OptionsBuilder) apicommon.OutgoingRPCHandler {
	options := new(options).Init()
	for _, optionsBuilder := range optionsBuilders {
		optionsBuilder(options)
	}
	return func(ctx context.Context, outgoingRPC *apicommon.OutgoingRPC) (returnedErr error) {
		subLogger := logger
		ctx = subLogger.WithContext(ctx)

		returnedErr = outgoingRPC.Do(ctx)

		if returnedErr == nil {
			if err := outgoingRPC.LoadResults(ctx); err != nil {
				logger.Err(err).Str("traceID", outgoingRPC.TraceID).
					Msg("results loading failed")
			}
		}
		var event *zerolog.Event
		if apicommon.ClientShouldReportError(returnedErr, outgoingRPC.StatusCode) {
			event = subLogger.Error()
		} else {
			event = subLogger.Info()
		}
		event.Str("traceID", outgoingRPC.TraceID)
		event.Str("fullMethodName", outgoingRPC.FullMethodName)
		event.Int("statusCode", outgoingRPC.StatusCode)
		event.Int("errorCode", int(outgoingRPC.ErrorCode))
		if returnedErr != nil {
			event.Str("err", returnedErr.Error())
		}
		event.Str("url", outgoingRPC.URL)
		paramsStr := apicommon.BytesToString(outgoingRPC.RawParams)
		if apicommon.DebugMode || len(paramsStr) < options.MaxParamsSize {
			event.Str("params", paramsStr)
		} else {
			event.Int("paramsSize", len(paramsStr))
			event.Str("truncatedParams", paramsStr[:options.MaxParamsSize])
		}
		resultsStr := apicommon.BytesToString(outgoingRPC.RawResults)
		if apicommon.DebugMode || len(resultsStr) < options.MaxResultsSize {
			event.Str("results", resultsStr)
		} else {
			event.Int("resultsSize", len(resultsStr))
			event.Str("truncatedResults", resultsStr[:options.MaxResultsSize])
		}
		event.Msg("outgoing rpc")
		return
	}
}
