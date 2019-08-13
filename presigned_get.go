package s3

import (
	"encoding/hex"
	"net/url"
)

type Presigner struct {
	input *inputModel
}

func NewPresigner(input *inputModel) *Presigner {
	return &Presigner{input: input}
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html
func (this *Presigner) GenerateURL() (string, error) {
	var (
		canonicalQuery   = this.task0_buildCanonicalQueryString()
		canonicalRequest = this.task1_composeCanonicalRequest(canonicalQuery)
		stringToSign     = this.task2_composeStringToSign(canonicalRequest)
		signature        = this.task3_calculateSignature(stringToSign)
	)
	return this.assembleURL(canonicalQuery, signature)
}

func (this *Presigner) task0_buildCanonicalQueryString() url.Values {
	query := make(url.Values)
	query.Set(HeaderAlgorithm, awsV4SignatureAlgorithm)
	query.Set(HeaderCredential, this.input.fullCredentialScope())
	query.Set(HeaderDate, this.input.timestampV4())
	query.Set(HeaderExpires, this.input.expiresInSeconds())
	query.Set(HeaderSignedHeaders, "host")
	return query
}

func (this *Presigner) task1_composeCanonicalRequest(query url.Values) string {
	return join("\n",
		"GET",
		"/"+this.input.key,
		normalizeQuery(query),
		"host:"+this.input.buildVirtualHostname(),
		"",     // blank line
		"host", // signed header names
		"UNSIGNED-PAYLOAD",
	)
}

func (this *Presigner) task2_composeStringToSign(canonicalRequest string) string {
	return join("\n",
		awsV4SignatureAlgorithm,
		this.input.timestampV4(),
		this.input.credentialScope(),
		hashSHA256([]byte(canonicalRequest)),
	)
}

func (this *Presigner) task3_calculateSignature(stringToSign string) string {
	signingKey := []byte(awsV4SignatureInitializationString + this.input.credential().SecretAccessKey)
	signingKey = hmacSHA256(signingKey, timestampDateV4(this.input.timestampV4()))
	signingKey = hmacSHA256(signingKey, this.input.region)
	signingKey = hmacSHA256(signingKey, "s3")
	signingKey = hmacSHA256(signingKey, awsV4CredentialScopeTerminationString)
	signingKey = hmacSHA256(signingKey, stringToSign)
	return hex.EncodeToString(signingKey)
}

func (this *Presigner) assembleURL(canonicalQuery url.Values, signature string) (string, error) {
	raw := this.input.buildVirtualHostingURL() +
		"?" + normalizeQuery(canonicalQuery) +
		"&" + HeaderSignature +
		"=" + signature
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

const (
	HeaderAlgorithm     = "X-Amz-Algorithm"
	HeaderCredential    = "X-Amz-Credential"
	HeaderDate          = "X-Amz-Date"
	HeaderExpires       = "X-Amz-Expires"
	HeaderSignedHeaders = "X-Amz-SignedHeaders"
	HeaderSignature     = "X-Amz-Signature"
)
