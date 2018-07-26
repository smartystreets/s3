package s3

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 is a wrapper over a subset of functionality implemented by:
// github.com/aws/aws-sdk-go-v2/service/s3.S3
type S3 struct{ service *s3.S3 }

// New builds *S3 with the provided options.
func New(options ...Option) *S3 {
	this := new(S3)
	apply(append(defaults, options...), this)
	return this
}

// SignedGetRequest creates a GET *http.Request for the specified resource, using any specified options.
func (this *S3) SignedGetRequest(bucket, key string, options ...Option) (*http.Request, error) {
	input := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
	request := this.service.GetObjectRequest(input)
	apply(options, request)
	return request.HTTPRequest, request.Sign()
}

// SignedGetRequest creates a PUT *http.Request for the specified resource and blob, using any specified options.
func (this *S3) SignedPutRequest(bucket, key string, blob io.ReadSeeker, options ...Option) (*http.Request, error) {
	input := s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: blob}
	request := this.service.PutObjectRequest(&input)
	apply(options, request)
	return request.HTTPRequest, request.Sign()
}
