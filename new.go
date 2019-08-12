package s3

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

func NewPresignedURL(options ...Option) (*url.URL, error) {
	input := newInput("GET", options)
	if input.expireSeconds <= 0 {
		return nil, errors.New("use of ExpireSeconds with positive value is required")
	}
	if input.bucket == "" {
		return nil, errors.New("use of Bucket option with non-blank value is required")
	}
	if input.key == "" {
		return nil, errors.New("use of Key option with non-blank value is required")
	}
	if len(input.credentials) == 0 {
		return nil, errors.New("use of Credentials is required")
	}
	now := timestampV4()
	query := make(url.Values)
	query.Set("X-Amz-Algorithm", awsV4SignatureAlgorithm)
	query.Set("X-Amz-Credential", input.credentials[0].AccessKeyID+"/"+credentialScope(timestampDateV4(now), input.region))
	query.Set("X-Amz-Date", now)
	query.Set("X-Amz-Expires", strconv.Itoa(input.expireSeconds))
	query.Set("X-Amz-SignedHeaders", "host")
	query.Set("X-Amz-Signature", "")

	initial := input.buildURL()
	data := initializeRequestData(
		input.region, GET, initial.Path, query.Encode(),
		"UNSIGNED-PAYLOAD",
		now, http.Header{"host": []string{initial.Host}},
	)
	signer := newV4Signer(data, input.credentials[0])
	signature := signer.calculateSignature()
	query.Set("X-Amz-Signature", signature.task3_IntermediateSignature)
	initial.RawQuery = query.Encode()

	return initial, nil
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
