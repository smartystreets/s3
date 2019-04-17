package s3

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type inputModel struct {
	credentials []awsCredentials

	method   string
	endpoint string
	region   string
	bucket   string
	key      string

	expireTime time.Time
	etag       string

	content         io.ReadSeeker
	contentType     string
	contentEncoding string
	contentMD5      string
	contentLength   int64

	serverSideEncryption ServerSideEncryptionValue
}

func newInput(method string, options []Option) *inputModel {
	return new(inputModel).applyOptions(append(options, method_(method)))
}

func (this *inputModel) applyOptions(options []Option) *inputModel {
	for _, option := range options {
		if option != nil {
			option(this)
		}
	}
	if len(this.credentials) == 0 {
		AmbientCredentials()(this)
	}
	return this
}

func (this *inputModel) validate() error {
	if this.method != GET && this.method != PUT {
		return ErrInvalidRequestMethod
	}
	if len(this.bucket) == 0 {
		return ErrBucketMissing
	}
	if len(this.key) == 0 {
		return ErrKeyMissing
	}
	if this.method == PUT && this.content == nil {
		return ErrContentMissing
	}
	return nil
}

func (this *inputModel) buildAndSignRequest() (request *http.Request, err error) {
	request, err = http.NewRequest(this.method, this.buildURL(), this.content)
	if err != nil {
		return nil, err
	}

	this.prepareRequestForSigning(request)
	signature := calculateAWSv4Signature(this.region, request, this.credentials[0])
	request.Header.Set("Authorization", signature)
	return request, nil
}

func (this *inputModel) prepareRequestForSigning(request *http.Request) {
	if request.URL.Path == "" {
		request.URL.Path += "/"
	}
	if this.contentLength > 0 {
		request.ContentLength = this.contentLength
	}
	if len(this.contentType) == 0 {
		this.contentType = "application/x-www-form-urlencoded; charset=utf-8"
	}
	setHeader(request, "Host", request.Host) // This must be included in range of headers to sign
	setHeader(request, "Content-Length", formatInt64(this.contentLength))
	setHeader(request, "Content-Encoding", this.contentEncoding)
	setHeader(request, "Content-Type", this.contentType)
	setHeader(request, "Content-MD5", this.contentMD5)
	setHeader(request, "If-None-Match", this.etag)
	setHeader(request, "X-Amz-Server-Side-Encryption", string(this.serverSideEncryption))
	setHeader(request, "X-Amz-Security-Token", this.credentials[0].SecurityToken)
	setHeader(request, "X-Amz-Content-Sha256", hashSHA256(readAndReplaceBody(request)))
	setHeader(request, "X-Amz-Expires", formatUnixTimeStamp(this.expireTime))
	setHeader(request, "X-Amz-Date", timestampV4())
}
func (this *inputModel) buildURL() string {
	builder := new(strings.Builder)

	if len(this.endpoint) > 0 {
		builder.WriteString(this.endpoint)
	} else {
		builder.WriteString("https://s3")
		if len(this.region) > 0 {
			builder.WriteString("-")
			builder.WriteString(this.region)
		}
		builder.WriteString(".amazonaws.com")
	}

	if !strings.HasSuffix(builder.String(), "/") {
		builder.WriteString("/")
	}
	builder.WriteString(this.bucket)
	builder.WriteString("/")
	builder.WriteString(this.key)
	return builder.String()
}
