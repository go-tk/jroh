package apicommon

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"unsafe"
)

func BytesToString(bytes []byte) string { return *(*string)(unsafe.Pointer(&bytes)) }

func ServerShouldReportError(err error, statusCode int) bool {
	return err != nil && (statusCode/100 == 5 && !errors.Is(err, context.Canceled))
}

func ClientShouldReportError(err error, statusCode int) bool {
	return err != nil && (statusCode/100 == 4 || !errIsTemporary(err))
}

func errIsTemporary(err error) bool {
	for {
		if v, ok := err.(interface{ Temporary() bool }); ok && v.Temporary() {
			return true
		}
		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
	}
}

type TransportFunc func(request *http.Request) (response *http.Response, err error)

var _ http.RoundTripper = TransportFunc(nil)

func (tf TransportFunc) RoundTrip(request *http.Request) (*http.Response, error) { return tf(request) }

func MakeInMemoryTransport(handler http.Handler) http.RoundTripper {
	return TransportFunc(func(request *http.Request) (*http.Response, error) {
		newCtx, cancel := context.WithCancel(context.Background())
		responseCh := make(chan *http.Response, 1)
		go func() {
			responseRecorder := httptest.NewRecorder()
			newRequest := request.Clone(newCtx)
			if newRequest.Body == nil {
				newRequest.Body = ioutil.NopCloser(bytes.NewReader(nil))
			}
			handler.ServeHTTP(responseRecorder, newRequest)
			response := responseRecorder.Result()
			responseCh <- response
		}()
		ctx := request.Context()
		select {
		case <-ctx.Done():
			cancel()
			<-responseCh
			return nil, ctx.Err()
		case response := <-responseCh:
			cancel()
			return response, nil
		}
	})
}

type ReaderFunc func(buffer []byte) (dataSize int, err error)

var _ io.Reader = (ReaderFunc)(nil)

func (rf ReaderFunc) Read(buffer []byte) (int, error) { return rf(buffer) }
