package s3

import (
	"net/http"
	"time"
)

type v4RequestData struct {
	service          string
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

func initializeRequestData(request *http.Request, service, region, bodyDigest string) *v4RequestData {
	requestTimestamp := request.Header.Get("X-Amz-Date")
	canonicalHeaders, signedHeaders := canonicalAndSignedHeaders(request.Header)
	return &v4RequestData{
		service:          service,
		region:           region,
		method:           request.Method,
		urlPath:          normalizeURI(request.URL.Path),
		urlQuery:         normalizeQuery(request.URL.Query()),
		bodyDigest:       bodyDigest,
		timestamp:        requestTimestamp,
		date:             timestampDateV4(requestTimestamp),
		canonicalHeaders: canonicalHeaders,
		signedHeaders:    signedHeaders,
	}
}

func (this v4RequestData) credentialScope() string {
	return join("/", this.date, this.region, this.service, awsV4CredentialScopeTerminationString)
}

var utcNow = func() time.Time { return time.Now().UTC() }

func timestampV4() string                     { return utcNow().Format(timeFormatV4) }
func timestampDateV4(timestamp string) string { return timestamp[:8] }

const timeFormatV4 = "20060102T150405Z"
