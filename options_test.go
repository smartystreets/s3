package s3

import (
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

func (this *OptionsFixture) TestHardCodedCredentials() {
	request, err := SignedGetRequest("bucket", "key", Credentials("access", "secret"))
	this.So(err, should.BeNil)
	this.So(request.Header.Get("Authorization"), should.ContainSubstring, "Credential=access")
}

func (this *OptionsFixture) TestSignedGet_ExpireTimeForcesCreationOfSignatureInQueryString() {
	requestWithExpiration, _ := SignedGetRequest("bucket", "key", ExpireTime(time.Second*30))
	requestWithoutExpiration, _ := SignedGetRequest("bucket", "key")

	this.So(requestWithExpiration.URL.Query(), should.NotBeEmpty)
	this.So(requestWithoutExpiration.URL.Query(), should.BeEmpty)
}

func (this *OptionsFixture) TestSignedPut_ExpireTimeForcesCreationOfSignatureInQueryString() {
	requestWithExpiration, _ := SignedPutRequest("bucket", "key", nil, ExpireTime(time.Second*30))
	requestWithoutExpiration, _ := SignedPutRequest("bucket", "key", nil)

	this.So(requestWithExpiration.URL.Query(), should.NotBeEmpty)
	this.So(requestWithoutExpiration.URL.Query(), should.BeEmpty)
}

func (this *OptionsFixture) TestIfNoneMatchAddHeader() {
	request, _ := SignedGetRequest("bucket", "key", IfNoneMatch("etag"))
	this.So(request.Header.Get("If-None-Match"), should.Equal, "etag")
}

func (this *OptionsFixture) TestServerSideEncryption() {
	put, _ := SignedPutRequest("bucket", "key", strings.NewReader("hi"), ServerSideEncryption(ServerSideEncryptionAES256))
	this.So(put.Header.Get("x-amz-server-side-encryption"), should.Equal, ServerSideEncryptionAES256)
}

func (this *OptionsFixture) TestContentType() {
	put, _ := SignedPutRequest("bucket", "key", strings.NewReader("hi"), ContentType("application/boink"))
	this.So(put.Header.Get("Content-Type"), should.Equal, "application/boink")
}

func (this *OptionsFixture) TestContentEncoding() {
	put, _ := SignedPutRequest("bucket", "key", strings.NewReader("hi"), ContentEncoding("utf-8"))
	this.So(put.Header.Get("Content-Encoding"), should.Equal, "utf-8")
}

func (this *OptionsFixture) TestContentMD5() {
	put, _ := SignedPutRequest("bucket", "key", strings.NewReader("hi"), ContentMD5("abcdef01"))
	this.So(put.Header.Get("Content-MD5"), should.Equal, "abcdef01")
}

func (this *OptionsFixture) TestContentLength() {
	put, _ := SignedPutRequest("bucket", "key", strings.NewReader("hi"), ContentLength(42))
	this.So(put.ContentLength, should.Equal, 42)
	this.So(put.Header.Get("Content-Length"), should.Equal, "42")
}