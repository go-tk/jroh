package signaturverifier

import (
	"bytes"
	"context"
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

type KeyFetcher func(ctx context.Context, senderID string) (key []byte, ok bool, err error)

type OptionsBuilder func(options *options)

type options struct {
	MaxHeaderStrLength int
	TimestampGetter    func() int64
	MaxTimestampSkew   int64
}

func (o *options) Init() *options {
	o.MaxHeaderStrLength = 256
	o.TimestampGetter = func() int64 { return time.Now().Unix() }
	o.MaxTimestampSkew = 10
	return o
}

func MaxHeaderStrLength(value int) OptionsBuilder {
	return func(options *options) { options.MaxHeaderStrLength = value }
}

func TimestampGetter(value func() int64) OptionsBuilder {
	return func(options *options) { options.TimestampGetter = value }
}

func MaxTimestampSkew(value int64) OptionsBuilder {
	return func(options *options) { options.MaxTimestampSkew = value }
}

func NewForServer(keyFetcher KeyFetcher, optionsBuilders ...OptionsBuilder) apicommon.ServerMiddleware {
	options := new(options).Init()
	for _, optionsBuilder := range optionsBuilders {
		optionsBuilder(options)
	}
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			incomingRPC := apicommon.MustGetRPCFromContext(ctx).IncomingRPC()
			headerStr := r.Header.Get("Authorization")
			if headerStr == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if headerStrLength := len(headerStr); headerStrLength > options.MaxHeaderStrLength {
				err := fmt.Errorf("signaturverifier: header str too long; headerStrLength=%v maxHeaderStrLength=%v",
					headerStrLength, options.MaxHeaderStrLength)
				incomingRPC.RespondHTTPWithErr(w, http.StatusUnprocessableEntity, err, "")
				return
			}
			header, err := parseHeaderStr(headerStr)
			if err != nil {
				incomingRPC.RespondHTTPWithErr(w, http.StatusBadRequest, err, "")
				return
			}
			if now := options.TimestampGetter(); !checkTimestamp(header.Timestamp, now, options.MaxTimestampSkew) {
				err := fmt.Errorf("signaturverifier: unexpected timestamp; timestamp=%v maxTimestampSkew=%v",
					header.Timestamp, options.MaxTimestampSkew)
				incomingRPC.RespondHTTPWithErr(w, http.StatusUnprocessableEntity, err, "")
				return
			}
			key, ok, err := keyFetcher(ctx, header.SenderID)
			if err != nil {
				err := fmt.Errorf("signaturverifier: key fetching failed: %v", err)
				incomingRPC.RespondHTTPWithErr(w, http.StatusInternalServerError, err, "")
				return
			}
			if !ok {
				err := fmt.Errorf("signaturverifier: key not found; senderID=%q", header.SenderID)
				incomingRPC.RespondHTTPWithErr(w, http.StatusUnprocessableEntity, err, "")
				return
			}
			signature := makeSignature(
				header.Timestamp,
				header.SenderID,
				incomingRPC.FullMethodName(),
				incomingRPC.RawParams(),
				header.AlgorithmType,
				key,
			)
			if header.Signature != signature {
				err := fmt.Errorf("signaturverifier: unexpected signature; signature=%q", header.Signature)
				incomingRPC.RespondHTTPWithErr(w, http.StatusUnprocessableEntity, err, "")
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}

type header struct {
	Timestamp     int64
	SenderID      string
	AlgorithmType algorithmType
	Signature     string
}

func parseHeaderStr(headerStr string) (header, error) {
	var (
		timestamp        int64
		senderID         string
		algorithmTypeStr string
		signature        string
	)
	n, _ := fmt.Sscanf(headerStr, "Signature t=%d,sid=%q,at=%q,s=%q", &timestamp, &senderID, &algorithmTypeStr, &signature)
	if n != 4 {
		return header{}, fmt.Errorf("signaturverifier: invalid header str; headerStr=%q", headerStr)
	}
	var algorithmType algorithmType
	switch algorithmTypeStr {
	case "md5":
		algorithmType = algorithmMD5
	case "sha1":
		algorithmType = algorithmSHA1
	case "sha256":
		algorithmType = algorithmSHA256
	case "sha512":
		algorithmType = algorithmSHA512
	default:
		return header{}, fmt.Errorf("signaturverifier: unknown algorithm type; algorithmTypeStr=%q", algorithmTypeStr)
	}
	return header{
		Timestamp:     timestamp,
		SenderID:      senderID,
		AlgorithmType: algorithmType,
		Signature:     signature,
	}, nil
}

func checkTimestamp(timestamp, now, maxTimestampSkew int64) bool {
	if minTimestamp := now - maxTimestampSkew; timestamp < minTimestamp {
		return false
	}
	if maxTimestamp := now + maxTimestampSkew; timestamp > maxTimestamp {
		return false
	}
	return true
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
