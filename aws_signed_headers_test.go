package s3

import (
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestSignedHeadersFixture(t *testing.T) {
	gunit.Run(new(SignedHeadersFixture), t)
}

type SignedHeadersFixture struct {
	*gunit.Fixture
}

func (this *SignedHeadersFixture) TestCanonicalHeaders() {
	actualCanonicalHeaders, actualSignedHeaders := canonicalAndSignedHeaders(http.Header{
		"Host":         []string{"iam.amazonaws.com:1234"},
		"Content-Type": []string{"application/x-www-form-urlencoded; charset=utf-8"},
		"My-header1":   []string{"    a   b   c "},
		"X-Amz-Date":   []string{"20150830T123600Z"},
		"My-header2":   []string{"    \"a   b   c\" "},
	})

	expectedCanonicalHeaders := `content-type:application/x-www-form-urlencoded; charset=utf-8
host:iam.amazonaws.com
my-header1:a b c
my-header2:"a b c"
x-amz-date:20150830T123600Z
`
	expectedSignedHeaders := `content-type;host;my-header1;my-header2;x-amz-date`

	this.So(actualCanonicalHeaders, should.Equal, expectedCanonicalHeaders)
	this.So(actualSignedHeaders, should.Equal, expectedSignedHeaders)
}
