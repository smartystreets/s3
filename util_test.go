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

func (this *BucketKeyFixture) assertRegionBucketKey(input, expectedRegion, expectedBucket, expectedKey string) {
	address, _ := url.Parse(input)
	region, bucket, key := RegionBucketKey(address)
	this.So(region, should.Equal, expectedRegion)
	this.So(bucket, should.Equal, expectedBucket)
	this.So(key, should.Equal, expectedKey)
}

func (this *BucketKeyFixture) Test() {
	this.assertRegionBucketKey("", "", "", "")
	this.assertRegionBucketKey("https://s3.amazonaws.com", "", "", "")
	this.assertRegionBucketKey("https://s3.amazonaws.com/bucket", "", "bucket", "")
	this.assertRegionBucketKey("https://s3.amazonaws.com/bucket/key", "", "bucket", "key")
	this.assertRegionBucketKey("https://s3.amazonaws.com/bucket/k/e/y", "", "bucket", "k/e/y")
	this.assertRegionBucketKey("https://s3-region.amazonaws.com/bucket/key", "region", "bucket", "key")
	this.assertRegionBucketKey("https://bucket.s3.amazonaws.com/key", "", "bucket", "key")
	this.assertRegionBucketKey("https://bucket.s3-region.amazonaws.com/key", "region", "bucket", "key")
}
