package s3

import (
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestBucketKeyFixture(t *testing.T) {
	gunit.Run(new(BucketKeyFixture), t)
}

type BucketKeyFixture struct {
	*gunit.Fixture
}

func parseURL(input string) *url.URL {
	parsed, _ := url.Parse(input)
	return parsed
}

func (this *BucketKeyFixture) assertRegionBucketKey(input *url.URL, expectedRegion, expectedBucket, expectedKey string) {
	region, bucket, key := RegionBucketKey(input)
	this.So(region, should.Equal, expectedRegion)
	this.So(bucket, should.Equal, expectedBucket)
	this.So(key, should.Equal, expectedKey)
}

func (this *BucketKeyFixture) Test() {
	this.assertRegionBucketKey(nil, "", "", "")
	this.assertRegionBucketKey(parseURL(""), "", "", "")
	this.assertRegionBucketKey(parseURL("https://s3.amazonaws.com"), "", "", "")
	this.assertRegionBucketKey(parseURL("https://s3.amazonaws.com/bucket"), "", "bucket", "")
	this.assertRegionBucketKey(parseURL("https://s3.amazonaws.com/bucket/key"), "", "bucket", "key")
	this.assertRegionBucketKey(parseURL("https://s3.amazonaws.com/bucket/k/e/y"), "", "bucket", "k/e/y")
	this.assertRegionBucketKey(parseURL("https://s3-region.amazonaws.com/bucket/key"), "region", "bucket", "key")
	this.assertRegionBucketKey(parseURL("https://bucket.s3.amazonaws.com/key"), "", "bucket", "key")
	this.assertRegionBucketKey(parseURL("https://bucket.s3-region.amazonaws.com/key"), "region", "bucket", "key")
}
