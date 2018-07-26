package s3

import (
	"io"
	"net/http"
)

// SignedGetRequest creates a GET *http.Request for the specified resource, using any specified options.
func SignedGetRequest(bucket, key string, options ...Option) (*http.Request, error) {
	return New(options...).SignedGetRequest(bucket, key, options...)
}

// SignedGetRequest creates a PUT *http.Request for the specified resource and blob, using any specified options.
func SignedPutRequest(bucket, key string, blob io.ReadSeeker, options ...Option) (*http.Request, error) {
	return New(options...).SignedPutRequest(bucket, key, blob, options...)
}
