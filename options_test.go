package s3

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestOptionsFixture(t *testing.T) {
	gunit.Run(new(OptionsFixture), t)
}

type OptionsFixture struct {
	*gunit.Fixture
}

func (this *OptionsFixture) TestMissingMethod() {
	request, err := NewRequest("", Bucket("bucket"), Key("key"))
	this.So(err, should.Equal, ErrInvalidRequestMethod)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) TestInvalidMethod() {
	request, err := NewRequest("POST", Bucket("bucket"), Key("key"))
	this.So(err, should.Equal, ErrInvalidRequestMethod)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) Test_GET_MissingBucket() {
	request, err := NewRequest(GET, Key("key"))
	this.So(err, should.Equal, ErrBucketMissing)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) Test_GET_MissingKey() {
	request, err := NewRequest(GET, Bucket("bucket"))
	this.So(err, should.Equal, ErrKeyMissing)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) Test_PUT_MissingContent() {
	request, err := NewRequest(PUT, Bucket("bucket"), Key("key"))
	this.So(err, should.Equal, ErrContentMissing)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) TestZeroLengthKey() {
	request, err := NewRequest(GET, Bucket("bucket"), Key(""))
	this.So(err, should.Equal, ErrKeyMissing)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) TestZeroLengthBucket() {
	request, err := NewRequest(GET, Bucket(""), Key("key"))
	this.So(err, should.Equal, ErrBucketMissing)
	this.So(request, should.BeNil)
}

func (this *OptionsFixture) TestHardCodedRegionOption() {
	request, err := NewRequest(GET, Region("us-west-1"), Bucket("bucket"), Key("key"))
	this.So(err, should.BeNil)
	this.So(request.URL.Host, should.ContainSubstring, "us-west-1")
}

func (this *OptionsFixture) TestHardCodedCredentials() {
	request, err := NewRequest(GET, Bucket("bucket"), Key("key"), Credentials("access", "secret"))
	this.So(err, should.BeNil)
	this.So(request.Header.Get("Authorization"), should.ContainSubstring, "Credential=access")
}

func (this *OptionsFixture) TestBucketAndKey() {
	request, err := NewRequest(GET, Bucket("bucket"), Key("/key/"))
	this.So(err, should.BeNil)
	this.So(request.URL.Path, should.ContainSubstring, "bucket")
	this.So(request.URL.Path, should.ContainSubstring, "key")
}

func (this *OptionsFixture) TestEndpoint() {
	request, err := NewRequest(GET, Endpoint("http://localhost:9000"), Bucket("bucket"), Key("key"))
	this.So(err, should.BeNil)
	this.So(request.URL.Scheme, should.Equal, "http")
	this.So(request.URL.Host, should.Equal, "localhost:9000")
	this.So(request.URL.Path, should.Equal, "/bucket/key")
}

func (this *OptionsFixture) TestExpireTimeAppearsInHeader() {
	now := time.Now()
	requestWithExpiration, _ := NewRequest(GET, Bucket("bucket"), Key("key"), ExpireTime(now))
	requestWithoutExpiration, _ := NewRequest(GET, Bucket("bucket"), Key("key"))

	this.So(requestWithExpiration.Header.Get("x-amz-expires"), should.Equal, fmt.Sprint(now.Unix()))
	this.So(requestWithoutExpiration.Header.Get("x-amz-expires"), should.BeBlank)
}

func (this *OptionsFixture) TestIfNoneMatchAddHeader() {
	request, _ := NewRequest(GET, Bucket("bucket"), Key("key"), IfNoneMatch("etag"))
	this.So(request.Header.Get("If-None-Match"), should.Equal, "etag")
}

func (this *OptionsFixture) TestServerSideEncryption() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentString("hi"), ServerSideEncryption(ServerSideEncryptionAES256))
	this.So(put.Header.Get("x-amz-server-side-encryption"), should.Equal, ServerSideEncryptionAES256)
}

func (this *OptionsFixture) TestContentType() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentString("hi"), ContentType("application/boink"))
	this.So(put.Header.Get("Content-Type"), should.Equal, "application/boink")
}

func (this *OptionsFixture) TestContentEncoding() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentString("hi"), ContentEncoding("utf-8"))
	this.So(put.Header.Get("Content-Encoding"), should.Equal, "utf-8")
}

func (this *OptionsFixture) TestContentMD5() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentString("hi"), ContentMD5("abcdef01"))
	this.So(put.Header.Get("Content-MD5"), should.Equal, "abcdef01")
}

func (this *OptionsFixture) TestContentLength() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentString("hi"), ContentLength(42))
	this.So(put.ContentLength, should.Equal, 42)
	this.So(put.Header.Get("Content-Length"), should.Equal, "42")
}

func (this *OptionsFixture) TestPUT_ContentBytes() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentBytes([]byte("hi")))
	all, _ := ioutil.ReadAll(put.Body)
	this.So(string(all), should.Equal, "hi")
	this.So(put.Header.Get("Content-Length"), should.Equal, "2")
}

func (this *OptionsFixture) TestPUT_ContentString() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), ContentString("hi"))
	all, _ := ioutil.ReadAll(put.Body)
	this.So(string(all), should.Equal, "hi")
	this.So(put.Header.Get("Content-Length"), should.Equal, "2")
}

func (this *OptionsFixture) TestPUT_Content() {
	put, _ := NewRequest(PUT, Bucket("bucket"), Key("key"), Content(strings.NewReader("hi")))
	all, _ := ioutil.ReadAll(put.Body)
	this.So(string(all), should.Equal, "hi")
}

func (this *OptionsFixture) TestStorageAddress() {
	address := &url.URL{Scheme: "https", Host: "bucket.s3.us-west-1.amazonaws.com", Path: "/key", RawPath: "/key"}
	request, _ := NewRequest(GET, StorageAddress(address))
	this.So(request.URL.String(), should.Equal, "https://s3-us-west-1.amazonaws.com/bucket/key")
}

func (this *OptionsFixture) TestStorageAddress_AlternateEndpoint() {
	address := &url.URL{Scheme: "https", Host: "1.2.3.4", Path: "/bucket/key", RawPath: "/bucket/key"}
	request, _ := NewRequest(GET, StorageAddress(address))
	this.So(request.URL.String(), should.Resemble, address.String())
}

func (this *OptionsFixture) TestStorageAddressWithKeyAsSeparateOptions() {
	address := &url.URL{Scheme: "https", Host: "bucket.s3.us-west-1.amazonaws.com"}
	request, _ := NewRequest(GET, StorageAddress(address), Key("key"))
	this.So(request.URL.String(), should.Equal, "https://s3-us-west-1.amazonaws.com/bucket/key")
}

func (this *OptionsFixture) TestMultipleKeysAreCombinedAsPathElements() {
	address := &url.URL{Scheme: "https", Host: "bucket.s3.us-west-1.amazonaws.com", Path: "/a/"}
	request, _ := NewRequest(GET,
		StorageAddress(address), // This option will include Key("/a/").
		Key("/b/"),
		Key("/c/"),
	)
	this.So(request.URL.Path, should.EndWith, "/a/b/c")
}
