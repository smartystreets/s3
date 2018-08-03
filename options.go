package s3

import (
	"bytes"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Option defines a callback for configuring the service and subsequent requests.
// It's important that this remain an exported name so that users can create slices
// of custom/conditional options.
type Option func(in *inputModel)

// Nop is a no-op. Useful as a placeholder in certain situations.
func Nop(_ *inputModel) {}

// CompositeOption allows multiple options to appear as one. This is handy
// when some options are unchanging (like credentials or a bucket name).
// Bundle them together in a single option to leave more room for the dynamic options.
func CompositeOption(options ...Option) Option {
	return func(in *inputModel) {
		in.applyOptions(options)
	}
}

// Region allows the user to specify the region for sending requests.
func Region(value string) Option {
	return func(in *inputModel) { in.region = external.WithRegion(value) }
}

// Bucket allows the user to specify the bucket for sending requests.
func Bucket(value string) Option {
	return func(in *inputModel) { in.bucket = &value }
}

// Key allows the user to specify the key for sending requests.
func Key(value string) Option {
	return func(in *inputModel) {
		if in.key == nil {
			in.key = aws.String(TrimKey(value))
		} else {
			in.key = aws.String(path.Join(*in.key, value))
		}
	}
}

// StorageAddress allows the user to specify the region, bucket, and/or key
// for sending requests from the provided S3 URL.
func StorageAddress(value *url.URL) Option {
	region, bucket, key := RegionBucketKey(value)
	return CompositeOption(
		ConditionalOption(Region(region), len(region) > 0),
		ConditionalOption(Bucket(bucket), len(bucket) > 0),
		ConditionalOption(Key(key), len(key) > 0),
	)
}

// ConditionalOption returns the option if condition == true, otherwise returns nil (nop).
func ConditionalOption(option Option, condition bool) Option {
	if condition {
		return option
	} else {
		return nil
	}
}

// Credentials allows the user to specify hard-coded credential values for sending requests.
func Credentials(access, secret string) Option {
	return func(in *inputModel) {
		in.credentials = external.WithCredentialsValue{
			AccessKeyID:     access,
			SecretAccessKey: secret,
		}
	}
}

// IfNoneMatch specifies the "If-None-Match" header. See the docs for details:
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html#RESTObjectGET-requests-in-headers
// This option only applies to GET requests.
func IfNoneMatch(etag string) Option {
	return func(in *inputModel) { in.etag = &etag }
}

// ExpireTime specifies an expiration for the generated input.
func ExpireTime(value time.Duration) Option {
	return func(in *inputModel) { in.expireTime = value }
}

// ContentString specifies the PUT request payload from a string.
func ContentString(value string) Option {
	return func(in *inputModel) {
		in.content = strings.NewReader(value)
		in.contentLength = aws.Int64(int64(len(value)))
	}
}

// ContentBytes specifies the PUT request payload from a slice of bytes.
func ContentBytes(value []byte) Option {
	return func(in *inputModel) {
		in.content = bytes.NewReader(value)
		in.contentLength = aws.Int64(int64(len(value)))
	}
}

// Content specifies the PUT request payload from an io.ReadSeeker.
func Content(value io.ReadSeeker) Option {
	return func(in *inputModel) { in.content = value }
}

// ContentType specifies the Content Type of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentType(value string) Option {
	return func(in *inputModel) { in.contentType = &value }
}

// ContentLength specifies the Content Length in bytes of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentLength(value int64) Option {
	return func(in *inputModel) { in.contentLength = &value }
}

// ContentMD5 specifies the MD5 checksum of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentMD5(value string) Option {
	return func(in *inputModel) { in.contentMD5 = &value }
}

// ContentEncoding specifies the content encoding of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentEncoding(value string) Option {
	return func(in *inputModel) { in.contentEncoding = &value }
}

// ServerSideEncryption specifies the server-side encryption algorithm to use.
// This option only applies to SignedPutRequest.
func ServerSideEncryption(value s3.ServerSideEncryption) Option {
	return func(in *inputModel) { in.serverSideEncryption = value }
}

const (
	ServerSideEncryptionAES256 = s3.ServerSideEncryptionAes256
	ServerSideEncryptionAWSKMS = s3.ServerSideEncryptionAwsKms
)
