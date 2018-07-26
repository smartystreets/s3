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

// ExpireTime specifies an expiration for the generated request.
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

// ServerSideEncryption specifies the server-side encryption algorithm to use.
// This option only applies to SignedPutRequest.
func ServerSideEncryption(algorithm s3.ServerSideEncryption) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.PutObjectRequest:
			t.Input.ServerSideEncryption = algorithm
		}
	}
}

const (
	ServerSideEncryptionAES256 = s3.ServerSideEncryptionAes256
	ServerSideEncryptionAWSKMS = s3.ServerSideEncryptionAwsKms
)

// ContentType specifies the Content Type of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentType(value string) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.PutObjectRequest:
			t.Input.ContentType = aws.String(value)
		}
	}
}

// ContentLength specifies the Content Length in bytes of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentLength(value int64) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.PutObjectRequest:
			t.Input.ContentLength = aws.Int64(value)
		}
	}
}

// ContentMD5 specifies the MD5 checksum of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentMD5(value string) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.PutObjectRequest:
			t.Input.ContentMD5 = aws.String(value)
		}
	}
}

// ContentEncoding specifies the content encoding of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentEncoding(value string) Option {
	return func(in interface{}) {
		switch t := in.(type) {
		case s3.PutObjectRequest:
			t.Input.ContentEncoding = aws.String(value)
		}
	}
}
