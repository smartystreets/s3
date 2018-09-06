package s3

import (
	"net/url"
	"strings"
)

// TrimKey removes leading and trailing slashes from the value.
// Generally, S3 keys don't have leading and trailing slashes so
// this serves as a convenient sanitization function.
func TrimKey(key string) string {
	return strings.Trim(key, "/")
}

// EndpointRegionBucketKey returns the S3 endpoint, region, bucket, and key embedded in a URL.
// For details on how S3 urls are formed, please see the S3 docs:
// https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html
//
// S3 URL examples showing optional placement of bucket and region (whitespace added for alignment):
//
// virtual-style bucket, no region:    http://bucket.s3           .amazonaws.com
// virtual-style bucket, with region:  http://bucket.s3-aws-region.amazonaws.com
// path-style bucket, no region:       http://       s3           .amazonaws.com/bucket
// path-style bucket, with region:     http://       s3-aws-region.amazonaws.com/bucket
// path-style bucket, custom endpoint: http://                       42.43.44.45/bucket
func EndpointRegionBucketKey(address *url.URL) (endpoint, region, bucket, key string) {
	bucket, key = BucketKey(address)
	if address != nil {
		region = extractRegion(address.Host)
		if !strings.Contains(address.Host, "s3") {
			endpoint = address.Scheme + "://" + address.Host
		}
	}
	return endpoint, region, bucket, key
}

// BucketKey returns the S3 bucket and key embedded in an S3 URL.
// For details on how s3 urls are formed, please see the S3 docs:
// https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html
func BucketKey(address *url.URL) (bucket, key string) {
	if address == nil {
		return "", ""
	}
	if isPathStyleAddress(address.Host) {
		path := TrimKey(address.Path)
		elements := strings.Split(path, "/")
		bucket = elements[0]
		if len(elements) > 1 {
			key = strings.Join(elements[1:], "/")
		}
	} else {
		bucket = extractVirtualBucket(address.Host)
		key = TrimKey(address.Path)
	}
	return bucket, key
}

func isPathStyleAddress(host string) bool {
	return !strings.Contains(host, "s3") || strings.HasPrefix(host, "s3.") || strings.HasPrefix(host, "s3-")
}

func extractVirtualBucket(host string) string {
	bucketEnd := strings.Index(host, ".s3")
	if bucketEnd <= 0 {
		return ""
	}
	return host[:bucketEnd]
}

func extractRegion(host string) string {
	regionBegin := strings.Index(host, "s3") + 3
	regionEnd := strings.Index(host, ".amazonaws.com")
	if regionBegin < 0 || regionEnd <= regionBegin {
		return ""
	}
	return host[regionBegin:regionEnd]
}
