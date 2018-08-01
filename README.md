# s3
--
    import "github.com/smartystreets/s3"


## Usage

```go
const (
	GET = "GET"
	PUT = "PUT"
)
```

```go
const (
	ServerSideEncryptionAES256 = s3.ServerSideEncryptionAes256
	ServerSideEncryptionAWSKMS = s3.ServerSideEncryptionAwsKms
)
```

```go
var (
	ErrInvalidRequestMethod = errors.New("Invalid method.")
	ErrBucketMissing        = errors.New("Bucket is required.")
	ErrKeyMissing           = errors.New("Key is required.")
	ErrContentMissing       = errors.New("Content is required.")
)
```

#### func  BucketKey

```go
func BucketKey(address *url.URL) (bucket, key string)
```
BucketKey returns the S3 bucket and key embedded in an S3 URL. For details on
how s3 urls are formed, please see the S3 docs:
https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html

#### func  NewRequest

```go
func NewRequest(method string, options ...Option) (*http.Request, error)
```

#### func  Nop

```go
func Nop(_ *inputModel)
```
Nop is a no-op. Useful as a placeholder in certain situations.

#### func  RegionBucketKey

```go
func RegionBucketKey(address *url.URL) (region, bucket, key string)
```
RegionBucketKey returns the S3 region, bucket, and key embedded in an S3 URL.
For details on how S3 urls are formed, please see the S3 docs:
https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html

S3 URL examples showing optional placement of bucket and region (whitespace
added for alignment):

virtual-style bucket, no region: http://bucket.s3 .amazonaws.com virtual-style
bucket, with region: http://bucket.s3-aws-region.amazonaws.com path-style
bucket, no region: http:// s3 .amazonaws.com/bucket path-style bucket, with
region: http:// s3-aws-region.amazonaws.com/bucket

#### func  TrimKey

```go
func TrimKey(key string) string
```
TrimKey removes leading and trailing slashes from the value. Generally, S3 keys
don't have leading and trailing slashes so this serves as a convenient
sanitization function.

#### type Option

```go
type Option func(in *inputModel)
```

Option defines a callback for configuring the service and subsequent requests.
It's important that this remain an exported name so that users can create slices
of custom/conditional options.

#### func  Bucket

```go
func Bucket(value string) Option
```
Bucket allows the user to specify the bucket for sending requests.

#### func  CompositeOption

```go
func CompositeOption(options ...Option) Option
```
CompositeOption allows multiple options to appear as one. This is handy when
some options are unchanging (like credentials or a bucket name). Bundle them
together in a single option to leave more room for the dynamic options.

#### func  Content

```go
func Content(value io.ReadSeeker) Option
```
Content specifies the PUT request payload from an io.ReadSeeker.

#### func  ContentBytes

```go
func ContentBytes(value []byte) Option
```
ContentBytes specifies the PUT request payload from a slice of bytes.

#### func  ContentEncoding

```go
func ContentEncoding(value string) Option
```
ContentEncoding specifies the content encoding of the payload/blob. This option
only applies to SignedPutRequest.

#### func  ContentLength

```go
func ContentLength(value int64) Option
```
ContentLength specifies the Content Length in bytes of the payload/blob. This
option only applies to SignedPutRequest.

#### func  ContentMD5

```go
func ContentMD5(value string) Option
```
ContentMD5 specifies the MD5 checksum of the payload/blob. This option only
applies to SignedPutRequest.

#### func  ContentString

```go
func ContentString(value string) Option
```
ContentString specifies the PUT request payload from a string.

#### func  ContentType

```go
func ContentType(value string) Option
```
ContentType specifies the Content Type of the payload/blob. This option only
applies to SignedPutRequest.

#### func  Credentials

```go
func Credentials(access, secret string) Option
```
Credentials allows the user to specify hard-coded credential values for sending
requests.

#### func  ExpireTime

```go
func ExpireTime(value time.Duration) Option
```
ExpireTime specifies an expiration for the generated input.

#### func  IfNoneMatch

```go
func IfNoneMatch(etag string) Option
```
IfNoneMatch specifies the "If-None-Match" header. See the docs for details:
https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html#RESTObjectGET-requests-in-headers
This option only applies to GET requests.

#### func  Key

```go
func Key(value string) Option
```
Key allows the user to specify the key for sending requests.

#### func  Region

```go
func Region(value string) Option
```
Region allows the user to specify the region for sending requests.

#### func  ServerSideEncryption

```go
func ServerSideEncryption(value s3.ServerSideEncryption) Option
```
ServerSideEncryption specifies the server-side encryption algorithm to use. This
option only applies to SignedPutRequest.

#### func  StorageAddress

```go
func StorageAddress(value *url.URL) Option
```
StorageAddress allows the user to specify the region, bucket, and/or key for
sending requests from the provided S3 URL.
