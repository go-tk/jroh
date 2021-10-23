package signatureattacher

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"time"

	"github.com/go-tk/jroh/go/apicommon"
)

type OptionsBuilder func(options *options)

type options struct {
	AlgorithmType algorithmType
}

func (o *options) Init() *options {
	o.AlgorithmType = algorithmMD5
	return o
}

func AlgorithmMD5() OptionsBuilder {
	return func(options *options) { options.AlgorithmType = algorithmMD5 }
}

func AlgorithmSHA1() OptionsBuilder {
	return func(options *options) { options.AlgorithmType = algorithmSHA1 }
}

func AlgorithmSHA256() OptionsBuilder {
	return func(options *options) { options.AlgorithmType = algorithmSHA256 }
}

func AlgorithmSHA512() OptionsBuilder {
	return func(options *options) { options.AlgorithmType = algorithmSHA512 }
}

func NewForClient(senderID string, key []byte, optionsBuilders ...OptionsBuilder) apicommon.ClientMiddleware {
	options := new(options).Init()
	for _, optionsBuilder := range optionsBuilders {
		optionsBuilder(options)
	}
	return func(transport http.RoundTripper) http.RoundTripper {
		return apicommon.TransportFunc(func(request *http.Request) (returnedResponse *http.Response, returnedErr error) {
			outgoingRPC := apicommon.MustGetRPCFromContext(request.Context()).OutgoingRPC()
			timestamp := time.Now().Unix()
			signature := makeSignature(
				timestamp,
				senderID,
				outgoingRPC.FullMethodName(),
				outgoingRPC.RawParams(),
				options.AlgorithmType,
				key,
			)
			headerStr := dumpHeader(header{
				Timestamp:     timestamp,
				SenderID:      senderID,
				AlgorithmType: options.AlgorithmType,
				Signature:     signature,
			})
			request.Header.Set("Authorization", headerStr)
			returnedResponse, returnedErr = transport.RoundTrip(request)
			return
		})
	}
}

type algorithmType int

const (
	algorithmMD5 = 1 + iota
	algorithmSHA1
	algorithmSHA256
	algorithmSHA512
)

func makeSignature(
	timestamp int64,
	senderID string,
	recipientID string,
	message []byte,
	algorithmType algorithmType,
	key []byte,
) string {
	var f func() hash.Hash
	switch algorithmType {
	case algorithmMD5:
		f = md5.New
	case algorithmSHA1:
		f = sha1.New
	case algorithmSHA256:
		f = sha256.New
	case algorithmSHA512:
		f = sha512.New
	default:
		panic("unreachable code")
	}
	hash := hmac.New(f, key)
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "t=%d,sid=%s,rid=%s,m=", timestamp, senderID, recipientID)
	hash.Write(buffer.Bytes())
	hash.Write(message)
	rawSignature := hash.Sum(nil)
	buffer.Reset()
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	encoder.Write(rawSignature)
	encoder.Close()
	signature := buffer.String()
	return signature
}

type header struct {
	Timestamp     int64
	SenderID      string
	AlgorithmType algorithmType
	Signature     string
}

func dumpHeader(header header) string {
	var algorithmTypeStr string
	switch header.AlgorithmType {
	case algorithmMD5:
		algorithmTypeStr = "md5"
	case algorithmSHA1:
		algorithmTypeStr = "sha1"
	case algorithmSHA256:
		algorithmTypeStr = "sha256"
	case algorithmSHA512:
		algorithmTypeStr = "sha512"
	default:
		panic("unreachable code")
	}
	return fmt.Sprintf(
		"Signature t=%d,sid=%q,at=%q,s=%q",
		header.Timestamp,
		header.SenderID,
		algorithmTypeStr,
		header.Signature,
	)
}
