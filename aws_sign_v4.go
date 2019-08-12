package s3

import (
	"encoding/hex"
	"fmt"
)

func calculateAWSv4Signature(data *v4RequestData, credentials awsCredentials) string {
	signer := newV4Signer(data, credentials)
	signature := signer.calculateSignature()
	return signature.task4_AuthorizationHeader
}

type v4Signer struct {
	keys awsCredentials
	data *v4RequestData
}

func newV4Signer(data *v4RequestData, credentials awsCredentials) *v4Signer {
	return &v4Signer{
		keys: credentials,
		data: data,
	}
}

type v4Signature struct {
	task1_CanonicalRequest      string
	task2_StringToSign          string
	task3_IntermediateSignature string
	task4_AuthorizationHeader   string
}

func (this *v4Signer) calculateSignature() v4Signature {
	task1 := this.task1_CanonicalRequest()
	task2 := this.task2_StringToSign(task1)
	task3 := this.task3_IntermediateSignature(task2)
	task4 := this.task4_AuthorizationHeader(task3)
	return v4Signature{
		task1_CanonicalRequest:      task1,
		task2_StringToSign:          task2,
		task3_IntermediateSignature: task3,
		task4_AuthorizationHeader:   task4,
	}
}

// TASK 1: https://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func (this *v4Signer) task1_CanonicalRequest() string {
	return join("\n",
		this.data.method,
		this.data.urlPath,
		this.data.urlQuery,
		this.data.canonicalHeaders,
		this.data.signedHeaders,
		this.data.bodyDigest,
	)
}

// TASK 2: https://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html
func (this *v4Signer) task2_StringToSign(canonicalRequest string) string {
	return join("\n",
		awsV4SignatureAlgorithm,
		this.data.timestamp,
		this.data.credentialScope(),
		hashSHA256([]byte(canonicalRequest)),
	)
}

// TASK 3: https://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html
func (this *v4Signer) task3_IntermediateSignature(stringToSign string) string {
	signingKey := []byte(awsV4SignatureInitializationString + this.keys.SecretAccessKey)
	signingKey = hmacSHA256(signingKey, this.data.date)
	signingKey = hmacSHA256(signingKey, this.data.region)
	signingKey = hmacSHA256(signingKey, s3)
	signingKey = hmacSHA256(signingKey, awsV4CredentialScopeTerminationString)
	signingKey = hmacSHA256(signingKey, stringToSign)
	return hex.EncodeToString(signingKey)
}

// TASK 4: https://docs.aws.amazon.com/general/latest/gr/sigv4-add-signature-to-request.html
func (this *v4Signer) task4_AuthorizationHeader(signature string) string {
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		awsV4SignatureAlgorithm,
		this.keys.AccessKeyID,
		this.data.credentialScope(),
		this.data.signedHeaders,
		signature,
	)
}

const (
	s3 = "s3"

	awsV4SignatureInitializationString    = "AWS4"
	awsV4CredentialScopeTerminationString = "aws4_request"
	awsV4SignatureAlgorithm               = "AWS4-HMAC-SHA256"
)
