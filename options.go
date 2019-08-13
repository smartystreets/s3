package s3

import (
	"bytes"
	"io"
	"net/url"
	"path"
	"strings"
	"time"
)

// Option defines a callback for configuring the service and subsequent requests.
// It's important that this remain an exported name so that users can create slices
// of custom/conditional options.
type Option func(in *inputModel)

func method_(value string) Option {
	return func(in *inputModel) { in.method = value }
}

// Nop is a no-op. Useful as a placeholder in certain situations.
func Nop(_ *inputModel) {}

// CompositeOption allows multiple options to appear as one. This is handy
// when some options are unchanging (like credentials or a bucket name).
// Bundle them together in a single option to leave more room for the dynamic options.
func CompositeOption(options ...Option) Option {
	return func(in *inputModel) {
		for _, option := range options {
			if option != nil {
				option(in)
			}
		}
	}
}

// Region allows the user to specify the region for sending requests.
func Region(value string) Option {
	return func(in *inputModel) { in.region = value }
}

// Bucket allows the user to specify the bucket for sending requests.
func Bucket(value string) Option {
	return func(in *inputModel) { in.bucket = value }
}

// Key allows the user to specify the key for sending requests.
func Key(value string) Option {
	return func(in *inputModel) {
		if len(in.key) == 0 {
			in.key = TrimKey(value)
		} else {
			in.key = path.Join(in.key, value)
		}
	}
}

// StorageAddress allows the user to specify the region, bucket, and/or key
// for sending requests from the provided S3 URL.
func StorageAddress(value *url.URL) Option {
	endpoint, region, bucket, key := EndpointRegionBucketKey(value)
	if len(endpoint) > 0 && len(region) == 0 {
		region = "us-east-1"
	}

	return CompositeOption(
		ConditionalOption(Endpoint(endpoint), len(endpoint) > 0),
		ConditionalOption(Region(region), len(region) > 0),
		ConditionalOption(Bucket(bucket), len(bucket) > 0),
		ConditionalOption(Key(key), len(key) > 0),
	)
}

// Endpoint allows the user to specify an alternate s3-compatible endpoint/URL to use for signed requests.
func Endpoint(value string) Option {
	return func(in *inputModel) { in.endpoint = value }
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
		in.credentials = append(in.credentials, awsCredentials{
			AccessKeyID:     access,
			SecretAccessKey: secret,
		})
	}
}

// STSCredentials allows the user to specify hard-coded credential values from AWS STS for sending requests.
func STSCredentials(access, secret, token string, expiration time.Time) Option {
	return func(in *inputModel) {
		in.credentials = append(in.credentials, awsCredentials{
			AccessKeyID:     access,
			SecretAccessKey: secret,
			SecurityToken:   token,
			Expiration:      expiration,
		})
	}
}

// IAMRoleCredentials loads credentials from the EC2 instance's configured IAM role. Only applicable when running on EC2.
func IAMRoleCredentials() Option {
	return func(in *inputModel) {
		in.credentials = append(in.credentials, getIAMRoleCredentials())
	}
}

// EnvironmentCredentials loads credentials from common variations of environment variables.
func EnvironmentCredentials() Option {
	return func(in *inputModel) {
		in.credentials = append(in.credentials, loadCredentialsFromEnvironment())
	}
}

// AmbientCredentials loads credentials first from the environment, then from any configured IAM role (on EC2).
func AmbientCredentials() Option {
	return func(in *inputModel) {
		credentials := ambientCredentials()
		if credentials.AccessKeyID == "" || credentials.SecretAccessKey == "" {
			return
		}
		in.credentials = append(in.credentials, credentials)
	}
}

// IfNoneMatch specifies the "If-None-Match" header. See the docs for details:
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html#RESTObjectGET-requests-in-headers
// This option only applies to GET requests.
func IfNoneMatch(etag string) Option {
	return func(in *inputModel) { in.etag = etag }
}

// ExpireTime specifies an expiration for the generated input.
func ExpireTime(value time.Time) Option {
	return func(in *inputModel) { in.expireTime = value }
}

// ExpiresIn specifies a validity period for PresignedGet (only).
func ExpiresIn(value time.Duration) Option {
	return func(in *inputModel) { in.expiresIn = value }
}

// ContentString specifies the PUT request payload from a string.
func ContentString(value string) Option {
	return func(in *inputModel) {
		in.content = strings.NewReader(value)
		in.contentLength = int64(len(value))
	}
}

// ContentBytes specifies the PUT request payload from a slice of bytes.
func ContentBytes(value []byte) Option {
	return func(in *inputModel) {
		in.content = bytes.NewReader(value)
		in.contentLength = int64(len(value))
	}
}

// Content specifies the PUT request payload from an io.ReadSeeker.
func Content(value io.ReadSeeker) Option {
	return func(in *inputModel) { in.content = value }
}

// ContentType specifies the Content Type of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentType(value string) Option {
	return func(in *inputModel) { in.contentType = value }
}

// ContentLength specifies the Content Length in bytes of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentLength(value int64) Option {
	return func(in *inputModel) { in.contentLength = value }
}

// ContentMD5 specifies the MD5 checksum of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentMD5(value string) Option {
	return func(in *inputModel) { in.contentMD5 = value }
}

// ContentEncoding specifies the content encoding of the payload/blob.
// This option only applies to SignedPutRequest.
func ContentEncoding(value string) Option {
	return func(in *inputModel) { in.contentEncoding = value }
}

// ServerSideEncryption specifies the server-side encryption algorithm to use.
// This option only applies to SignedPutRequest.
func ServerSideEncryption(value ServerSideEncryptionValue) Option {
	return func(in *inputModel) { in.serverSideEncryption = value }
}

// Timestamp specifies the timestamp to be included as the X-Amz-Date as well
// as for use in time based calculations. Helpful for testing.
func Timestamp(value time.Time) Option {
	return func(in *inputModel) { in.now = value }
}

type ServerSideEncryptionValue string

const (
	ServerSideEncryptionAES256 ServerSideEncryptionValue = "AES256"
	ServerSideEncryptionAWSKMS ServerSideEncryptionValue = "aws:kms"
)
