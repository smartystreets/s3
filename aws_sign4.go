package s3

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func calculateAWSv4Signature(request *http.Request, credentials ...awsCredentials) string {
	signer := newV4Signer(request, credentials...)
	signature := signer.calculateSignature()
	return signature.task4_AuthHeader
}

type v4Signer struct {
	request *http.Request
	keys    awsCredentials
	meta    *metadata
}

func newV4Signer(request *http.Request, credentials ...awsCredentials) *v4Signer {
	return &v4Signer{
		request: request,
		keys:    chooseKeys(credentials),
		meta:    new(metadata),
	}
}

type v4Signature struct {
	task1_HashedCanonicalRequest string
	task2_StringToSign           string
	task3_Signature              string
	task4_AuthHeader             string
}

func (this *v4Signer) calculateSignature() v4Signature {
	var (
		task1 = this.task1_HashedCanonicalRequest()
		task2 = this.task2_StringToSign(task1)
		task3 = this.task3_Signature(task2)
		task4 = this.task4_AuthHeader(task3)
	)
	return v4Signature{
		task1_HashedCanonicalRequest: task1,
		task2_StringToSign:           task2,
		task3_Signature:              task3,
		task4_AuthHeader:             task4,
	}
}

// TASK 1: https://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func (this *v4Signer) task1_HashedCanonicalRequest() string {
	var sortedHeaderKeys []string
	for key, _ := range this.request.Header {
		switch key {
		case "Content-Type", "Content-Md5", "Host":
		default:
			if !strings.HasPrefix(key, "X-Amz-") {
				continue
			}
		}
		sortedHeaderKeys = append(sortedHeaderKeys, strings.ToLower(key))
	}
	sort.Strings(sortedHeaderKeys)

	var headersToSign string
	for _, key := range sortedHeaderKeys {
		value := strings.TrimSpace(this.request.Header.Get(key))
		if key == "host" {
			//AWS does not include port in signing request.
			if strings.Contains(value, ":") {
				split := strings.Split(value, ":")
				port := split[1]
				if port == "80" || port == "443" {
					value = split[0]
				}
			}
		}
		headersToSign += key + ":" + value + "\n"
	}
	this.meta.signedHeaders = concat(";", sortedHeaderKeys...)
	canonicalRequest := concat("\n",
		this.request.Method,
		normuri(this.request.URL.Path),
		normquery(this.request.URL.Query()),
		headersToSign,
		this.meta.signedHeaders,
		this.request.Header.Get("X-Amz-Content-Sha256"),
	)

	return hashSHA256([]byte(canonicalRequest))
}

// TASK 2: https://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html
func (this *v4Signer) task2_StringToSign(hashedCanonReq string) string {
	requestTimestamp := this.request.Header.Get("X-Amz-Date")
	this.meta.algorithm = "AWS4-HMAC-SHA256"
	this.meta.service, this.meta.region = serviceAndRegion(this.request.Host)
	this.meta.date = timestampDateV4(requestTimestamp)
	this.meta.credentialScope = concat("/", this.meta.date, this.meta.region, this.meta.service, "aws4_request")
	return concat("\n", this.meta.algorithm, requestTimestamp, this.meta.credentialScope, hashedCanonReq)
}

// TASK 3: https://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html
func (this *v4Signer) task3_Signature(stringToSign string) string {
	signingKey := []byte("AWS4" + this.keys.SecretAccessKey)
	fmt.Println("1", hex.EncodeToString(signingKey))

	signingKey = hmacSHA256(signingKey, this.meta.date)
	fmt.Println("2", hex.EncodeToString(signingKey))

	signingKey = hmacSHA256(signingKey, this.meta.region)
	fmt.Println("3", hex.EncodeToString(signingKey))

	signingKey = hmacSHA256(signingKey, this.meta.service)
	fmt.Println("4", hex.EncodeToString(signingKey))

	signingKey = hmacSHA256(signingKey, "aws4_request")
	fmt.Println("5", hex.EncodeToString(signingKey))

	signingKey = hmacSHA256(signingKey, stringToSign)
	fmt.Println("6", hex.EncodeToString(signingKey))

	return hex.EncodeToString(signingKey)
}

// Task 4: https://docs.aws.amazon.com/general/latest/gr/sigv4-add-signature-to-request.html
func (this *v4Signer) task4_AuthHeader(signature string) string {
	return this.meta.algorithm +
		" Credential=" + this.keys.AccessKeyID + "/" + this.meta.credentialScope +
		", SignedHeaders=" + this.meta.signedHeaders +
		", Signature=" + signature
}

func timestampV4() string {
	return utcNow().Format(timeFormatV4)
}

func timestampDateV4(timestamp string) string {
	return timestamp[:8]
}

const timeFormatV4 = "20060102T150405Z"
