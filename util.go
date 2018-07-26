package s3

import (
	"net/url"
	"strings"
)

// BucketKey returns the S3 bucket and key embedded in an S3 URL.
// For details on how s3 urls are formed, please see the S3 docs:
// https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html
func BucketKey(address *url.URL) (bucket, key string) {
	if strings.HasPrefix(address.Host, "s3.") || strings.HasPrefix(address.Host, "s3-") { // path-style
		path := strings.Trim(address.Path, "/")
		elements := strings.Split(path, "/")
		if len(elements) > 0 {
			bucket = elements[0]
		}
		if len(elements) > 1 {
			key = strings.Join(elements[1:], "/")
		}
	} else { // virtual style
		bucket = strings.Split(address.Host, ".")[0]
		key = strings.Trim(address.Path, "/")
	}

	return bucket, key
}
