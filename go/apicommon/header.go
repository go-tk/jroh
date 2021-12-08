package apicommon

import (
	"net/http"
	"strconv"
)

const traceIDHeaderKey = "X-Jroh-Trace-Id"

func setTraceID(traceID string, header http.Header) {
	header[traceIDHeaderKey] = []string{traceID}
}

func getTraceID(header http.Header) (string, bool) {
	if headerValues := header[traceIDHeaderKey]; len(headerValues) >= 1 {
		return headerValues[0], true
	}
	return "", false
}

const errorCodeHeaderKey = "X-Jroh-Error-Code"

func setErrorCode(errorCode ErrorCode, header http.Header) {
	header[errorCodeHeaderKey] = []string{strconv.FormatInt(int64(errorCode), 10)}
}

func getErrorCode(header http.Header) (ErrorCode, bool) {
	if headerValues := header[errorCodeHeaderKey]; len(headerValues) >= 1 {
		i, _ := strconv.ParseInt(headerValues[0], 10, 32)
		return ErrorCode(i), true
	}
	return 0, false
}

func setContentTypeAsJSON(header http.Header) {
	header["Content-Type"] = []string{"application/json; charset=utf-8"}
}
