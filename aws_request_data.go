package s3

import (
	"net/http"
	"time"
)

type v4RequestData struct {
	region           string
	method           string
	urlPath          string
	urlQuery         string
	bodyDigest       string
	timestamp        string
	date             string
	canonicalHeaders string
	signedHeaders    string
}

func initializeRequestData(
	region, method, urlPath, urlQuery, bodyDigest, timestamp string,
	headers http.Header,
) *v4RequestData {
	canonicalHeaders, signedHeaders := canonicalAndSignedHeaders(headers)
	return &v4RequestData{
		region:           region,
		method:           method,
		urlPath:          normalizeURI(urlPath),
		urlQuery:         normalizeQuery(urlQuery),
		bodyDigest:       bodyDigest,
		timestamp:        timestamp,
		date:             timestampDateV4(timestamp),
		canonicalHeaders: canonicalHeaders,
		signedHeaders:    signedHeaders,
	}
}

func credentialScope(date, region string) string {
	return join("/", date, region, s3, awsV4CredentialScopeTerminationString)
}

func (this v4RequestData) credentialScope() string {
	return credentialScope(this.date, this.region)
}

var utcNow = func() time.Time { return time.Now().UTC() }

func timestampV4() string                     { return utcNow().Format(timeFormatV4) }
func timestampDateV4(timestamp string) string { return timestamp[:8] }

const timeFormatV4 = "20060102T150405Z"
