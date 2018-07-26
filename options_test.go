package s3

import (
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