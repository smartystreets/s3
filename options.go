package s3

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Option defines a callback for configuring the service and subsequent requests.
type Option func(in interface{})

var defaults = []Option{defaultCredentials}

func apply(options []Option, thing interface{}) {
	for _, option := range options {
		option(thing)
	}
}

func defaultCredentials(in interface{}) {
	switch t := in.(type) {
	case *S3:
		config, _ := external.LoadDefaultAWSConfig()
		t.service = s3.New(config)
	}
}

// Credentials allows the user to specify hard-coded credential values for sending requests.
// Only applicable when supplied to New().
func Credentials(access, secret string) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case *S3:
			config, _ := external.LoadDefaultAWSConfig(external.WithCredentialsValue(aws.Credentials{
				AccessKeyID:     access,
				SecretAccessKey: secret,
			}))
			t.service = s3.New(config)
		}
	}
}

// IfNoneMatch specifies the "If-None-Match" header. See the docs for details:
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html#RESTObjectGET-requests-request-headers
// This option only applies to SignedGetRequest().
func IfNoneMatch(etag string) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.GetObjectRequest:
			t.Input.IfNoneMatch = aws.String(etag)
		}
	}
}

// ExpireTime specifies an expiration for the generated request:
// This option applies to functions/methods that generate *http.Request.
func ExpireTime(validity time.Duration) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.GetObjectRequest:
			t.ExpireTime = validity
		case s3.PutObjectRequest:
			t.ExpireTime = validity
		}
	}
}
