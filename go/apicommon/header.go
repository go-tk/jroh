package apicommon

import (
	"net/http"
	"strconv"
	"time"
)

const traceIDHeaderKey = "X-Jroh-Trace-Id"

func setTraceID(traceID string, header http.Header) {
	header[traceIDHeaderKey] = []string{traceID}
}

func getTraceID(header http.Header) (string, bool) {
	if headerValues := header[traceIDHeaderKey]; len(headerValues) >= 1 {
		traceID := headerValues[0]
		return traceID, true
	}
	return "", false
}

const deadlineHeaderKey = "X-Jroh-Deadline"

func setDeadline(deadline time.Time, header http.Header) {
	timestamp := float64(deadline.UnixNano()) / float64(time.Second/time.Nanosecond)
	timestampStr := strconv.FormatFloat(timestamp, 'f', 3, 64)
	header[deadlineHeaderKey] = []string{timestampStr}
}

func getDeadline(header http.Header) (time.Time, bool) {
	if headerValues := header[deadlineHeaderKey]; len(headerValues) >= 1 {
		timestampStr := headerValues[0]
		timestamp, _ := strconv.ParseFloat(timestampStr, 64)
		deadline := time.Unix(0, int64(timestamp*float64(time.Second/time.Nanosecond)))
		return deadline, true
	}
	return time.Time{}, false
}

const errorCodeHeaderKey = "X-Jroh-Error-Code"

func setErrorCode(errorCode ErrorCode, header http.Header) {
	errorCodeStr := strconv.FormatInt(int64(errorCode), 10)
	header[errorCodeHeaderKey] = []string{errorCodeStr}
}

func getErrorCode(header http.Header) (ErrorCode, bool) {
	if headerValues := header[errorCodeHeaderKey]; len(headerValues) >= 1 {
		errorCodeStr := headerValues[0]
		n, _ := strconv.ParseInt(errorCodeStr, 10, 32)
		errorCode := ErrorCode(n)
		return errorCode, true
	}
	return 0, false
}

func setContentTypeAsJSON(header http.Header) {
	header["Content-Type"] = []string{"application/json; charset=utf-8"}
}
