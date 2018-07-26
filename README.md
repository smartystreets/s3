# s3
--
    import "github.com/smartystreets/s3"


## Usage

#### func  BucketKey

```go
func BucketKey(address *url.URL) (bucket, key string)
```
BucketKey returns the S3 bucket and key embedded in an S3 URL. For details on
how s3 urls are formed, please see the S3 docs:
https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html

#### func  SignedGetRequest

```go
func SignedGetRequest(bucket, key string, options ...Option) (*http.Request, error)
```
SignedGetRequest creates a GET *http.Request for the specified resource, using
any specified options.

#### func  SignedPutRequest

```go
func SignedPutRequest(bucket, key string, blob io.ReadSeeker, options ...Option) (*http.Request, error)
```
SignedGetRequest creates a PUT *http.Request for the specified resource and
blob, using any specified options.

#### type Option

```go
type Option func(in interface{})
```

Option defines a callback for configuring the service and subsequent requests.

#### func  Credentials

```go
func Credentials(access, secret string) Option
```
Credentials allows the user to specify hard-coded credential values for sending
requests. Only applicable when supplied to New().

#### func  ExpireTime

```go
func ExpireTime(validity time.Duration) Option
```
ExpireTime specifies an expiration for the generated request: This option
applies to functions/methods that generate *http.Request.

#### func  IfNoneMatch

```go
func IfNoneMatch(etag string) Option
```
IfNoneMatch specifies the "If-None-Match" header. See the docs for details:
https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectGET.html#RESTObjectGET-requests-request-headers
This option only applies to SignedGetRequest().

#### type S3

```go
type S3 struct {
}
```

S3 is a wrapper over a subset of functionality implemented by:
github.com/aws/aws-sdk-go-v2/service/s3.S3

#### func  New

```go
func New(options ...Option) *S3
```
New builds *S3 with the provided options.

#### func (*S3) SignedGetRequest

```go
func (this *S3) SignedGetRequest(bucket, key string, options ...Option) (*http.Request, error)
```
SignedGetRequest creates a GET *http.Request for the specified resource, using
any specified options.

#### func (*S3) SignedPutRequest

```go
func (this *S3) SignedPutRequest(bucket, key string, blob io.ReadSeeker, options ...Option) (*http.Request, error)
```
SignedGetRequest creates a PUT *http.Request for the specified resource and
blob, using any specified options.
