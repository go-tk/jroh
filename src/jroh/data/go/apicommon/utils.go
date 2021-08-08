package apicommon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
)

var DebugMode bool

type Validator interface {
	Validate() (err error)
}

func ReadParams(reader io.Reader, rpcInfo *RPCInfo) bool {
	rawParams, err := ioutil.ReadAll(reader)
	if err != nil {
		convertErr(err, rpcInfo)
		return false
	}
	rpcInfo.SetRawParams(rawParams)
	params := rpcInfo.Params()
	error := rpcInfo.Error()
	if err := json.Unmarshal(rawParams, params); err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			*error = *ErrParse.WithData(ErrorData{
				"details": err.Error(),
			})
		case *json.UnmarshalTypeError:
			*error = *ErrInvalidParams.WithData(ErrorData{
				"details": err.Error(),
			})
		default:
			convertErr(err, rpcInfo)
		}
		return false
	}
	if validator, ok := params.(Validator); ok {
		if err := validator.Validate(); err != nil {
			*error = *ErrInvalidParams.WithData(ErrorData{
				"details": err.Error(),
			})
			return false
		}
	}
	return true
}

func SaveErr(err error, rpcInfo *RPCInfo) {
	if err == nil {
		return
	}
	convertErr(err, rpcInfo)
}

func SavePanicValue(panicValue interface{}, rpcInfo *RPCInfo) {
	errStr := fmt.Sprintf("%v", panicValue)
	rpcInfo.SetInternalErr(errors.New(errStr))
	buffer := make([]byte, 4096)
	n := runtime.Stack(buffer, false)
	stackTrace := string(buffer[:n])
	rpcInfo.SetStackTrace(stackTrace)
	error := rpcInfo.Error()
	if DebugMode {
		*error = *ErrInternal.WithData(ErrorData{
			"details":    errStr,
			"stackTrace": stackTrace,
		})
	} else {
		*error = *ErrInternal
	}
}

func WriteResp(resp interface{}, responseWriter http.ResponseWriter, rpcInfo *RPCInfo) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(responseWriter).Encode(resp); err != nil {
		rpcInfo.SetRespWriteErr(err)
	}
}

func convertErr(err error, rpcInfo *RPCInfo) {
	error := rpcInfo.Error()
	if error2, ok := err.(*Error); ok {
		*error = *error2
		return
	}
	rpcInfo.SetInternalErr(err)
	if DebugMode {
		*error = *ErrInternal.WithData(ErrorData{
			"details": err.Error(),
		})
	} else {
		*error = *ErrInternal
	}
}
