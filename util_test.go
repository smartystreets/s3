package s3

import (
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestParsingFixture(t *testing.T) {
	gunit.Run(new(ParsingFixture), t)
}

type ParsingFixture struct {
	*gunit.Fixture
}

func URL(urlWithoutScheme string) (parsed *url.URL) {
	if len(urlWithoutScheme) > 0 {
		parsed, _ = url.Parse("https:" + urlWithoutScheme)
	}
	return parsed
}

func (this *ParsingFixture) assertFields(input *url.URL, expectedEndpoint, expectedRegion, expectedBucket, expectedKey string) {
	endpoint, region, bucket, key := EndpointRegionBucketKey(input)
	this.So(endpoint, should.Equal, expectedEndpoint)
	this.So(region, should.Equal, expectedRegion)
	this.So(bucket, should.Equal, expectedBucket)
	this.So(key, should.Equal, expectedKey)
}

func (this *ParsingFixture) Test() {
	this.assertFields(nil, "", "", "", "")
	this.assertFields(URL(""), "", "", "", "")
	this.assertFields(URL("//s3.amazonaws.com"), "", "", "", "")
	this.assertFields(URL("//s3.amazonaws.com/bucket"), "", "", "bucket", "")
	this.assertFields(URL("//s3.amazonaws.com/bucket/key"), "", "", "bucket", "key")
	this.assertFields(URL("//s3.amazonaws.com/bucket/k/e/y"), "", "", "bucket", "k/e/y")
	this.assertFields(URL("//s3-region.amazonaws.com/bucket/key"), "", "region", "bucket", "key")
	this.assertFields(URL("//bucket.s3.amazonaws.com/key"), "", "", "bucket", "key")
	this.assertFields(URL("//bucket.s3-region.amazonaws.com/key"), "", "region", "bucket", "key")
	this.assertFields(URL("//s3/bucket/key"), "https://s3", "", "bucket", "key")
	this.assertFields(URL("//localhost/bucket/key"), "https://localhost", "", "bucket", "key")
	this.assertFields(URL("//1.2.3.4:5678/bucket/key"), "https://1.2.3.4:5678", "", "bucket", "key")
}
