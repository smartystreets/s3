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

func (this *BucketKeyFixture) TestPathStyle() {
	this.assertBucketKey("", "", "")
	this.assertBucketKey("https://s3.amazonaws.com/", "", "")
	this.assertBucketKey("https://s3.amazonaws.com/bucket/k/e/y", "bucket", "k/e/y")
	this.assertBucketKey("https://s3-us-west-1.amazonaws.com/bucket/k/e/y", "bucket", "k/e/y")
	this.assertBucketKey("https://bucket.s3.amazonaws.com/k/e/y", "bucket", "k/e/y")
}
func (this *BucketKeyFixture) assertBucketKey(input, expectedBucket, expectedKey string) {
	address, _ := url.Parse(input)
	bucket, key := BucketKey(address)
	this.So(bucket, should.Equal, expectedBucket)
	this.So(key, should.Equal, expectedKey)
}
