package s3

import (
	"io"
	"net/http"
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

func (this *inputModel) buildRequest() (request *http.Request, err error) {
	request, err = http.NewRequest(this.method, this.buildURL(), this.content)
	if err != nil {
		return nil, err
	}

	// TODO: add headers to request...
	// TODO: add query string parameters to request...
	// TODO: Sign request using aws v4 signature...

	return request, nil
}

func (this *inputModel) buildURL() string {
	return "" // TODO: build initial url
}
