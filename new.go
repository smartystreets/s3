package s3

import (
	"errors"
	"net/http"
)

func NewPresignedGet(options ...Option) (string, error) {
	input := newInput(GET, options)

	if err := input.validate(); err != nil {
		return "", err
	}

	return NewPresigner(input).GenerateURL()
}

func NewRequest(method string, options ...Option) (*http.Request, error) {
	input := newInput(method, options)

	if err := input.validate(); err != nil {
		return nil, err
	}

	return input.buildAndSignRequest()
}

const (
	GET = "GET"
	PUT = "PUT"
)

var (
	ErrInvalidRequestMethod = errors.New("invalid method")
	ErrBucketMissing        = errors.New("bucket is required")
	ErrKeyMissing           = errors.New("key is required")
	ErrContentMissing       = errors.New("content is required")
)
