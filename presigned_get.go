package s3

import (
	"encoding/hex"
	"net/url"
)

type presigner struct {
	input *inputModel
}

func newPresigner(input *inputModel) *presigner {
	return &presigner{input: input}
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html
func (this *presigner) GenerateURL() (string, error) {
	var (
		canonicalQuery   = this.task0_buildCanonicalQueryString()
		canonicalRequest = this.task1_composeCanonicalRequest(canonicalQuery)
		stringToSign     = this.task2_composeStringToSign(canonicalRequest)
		signature        = this.task3_calculateSignature(stringToSign)
	)
	return this.assembleURL(canonicalQuery, signature)
}

func (this *presigner) task0_buildCanonicalQueryString() url.Values {
	query := make(url.Values)
	query.Set(headerAlgorithm, awsV4SignatureAlgorithm)
	query.Set(headerCredential, this.input.fullCredentialScope())
	query.Set(headerDate, this.input.timestampV4())
	query.Set(headerExpires, this.input.expiresInSeconds())
	query.Set(headerSignedHeaders, "host")
	return query
}

func (this *presigner) task1_composeCanonicalRequest(query url.Values) string {
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

func (this *presigner) task2_composeStringToSign(canonicalRequest string) string {
	return join("\n",
		awsV4SignatureAlgorithm,
		this.input.timestampV4(),
		this.input.credentialScope(),
		hashSHA256([]byte(canonicalRequest)),
	)
}

func (this *presigner) task3_calculateSignature(stringToSign string) string {
	signingKey := []byte(awsV4SignatureInitializationString + this.input.credential().SecretAccessKey)
	signingKey = hmacSHA256(signingKey, timestampDateV4(this.input.timestampV4()))
	signingKey = hmacSHA256(signingKey, this.input.region)
	signingKey = hmacSHA256(signingKey, "s3")
	signingKey = hmacSHA256(signingKey, awsV4CredentialScopeTerminationString)
	signingKey = hmacSHA256(signingKey, stringToSign)
	return hex.EncodeToString(signingKey)
}

func (this *presigner) assembleURL(canonicalQuery url.Values, signature string) (string, error) {
	raw := this.input.buildVirtualHostingURL() +
		"?" + normalizeQuery(canonicalQuery) +
		"&" + headerSignature +
		"=" + signature
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

const (
	headerAlgorithm     = "X-Amz-Algorithm"
	headerCredential    = "X-Amz-Credential"
	headerDate          = "X-Amz-Date"
	headerExpires       = "X-Amz-Expires"
	headerSignedHeaders = "X-Amz-SignedHeaders"
	headerSignature     = "X-Amz-Signature"
)
