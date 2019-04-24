package s3

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestV4SignerFixture(t *testing.T) {
	gunit.Run(new(V4SignerFixture), t)
}

// See: https://docs.aws.amazon.com/general/latest/gr/signature-v4-test-suite.html
type V4SignerFixture struct {
	*gunit.Fixture
	credentials awsCredentials
}

func (this *V4SignerFixture) Setup() {
	this.credentials = awsCredentials{
		AccessKeyID:     "AKIDEXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
	}
}

func (this *V4SignerFixture) TestAWSSignatureVersion4() {
	root := "testdata/aws-sig-v4-test-suite/aws-sig-v4-test-suite/"
	listing, err := ioutil.ReadDir(root)
	if !this.So(err, should.BeNil) {
		return
	}
	for _, item := range listing {
		name := item.Name()
		test := LoadTestCase(filepath.Join(root, name))
		if !test.IsValid() {
			this.Println("** Skipping folder:", name)
		} else {
			this.Println("** Running folder:", name)
			this.runTestCase(test)
		}
	}
}

func (this *V4SignerFixture) runTestCase(test *V4SignerTestCase) {
	request, bodyDigest := test.buildRequest()
	signature := newV4Signer("service", "us-east-1", bodyDigest, request, this.credentials).calculateSignature()

	this.So(signature.task1_CanonicalRequest, should.Equal, test.ExpectedCanonicalRequest)
	this.So(signature.task2_StringToSign, should.Equal, test.ExpectedStringToSign)
	this.So(signature.task4_AuthorizationHeader, should.Equal, test.ExpectedAuthorizationHeaderValue)
}

//////////////////////////////////////////////////////////////////

type V4SignerTestCase struct {
	RawRequestToSign                 string
	ExpectedCanonicalRequest         string
	ExpectedStringToSign             string
	ExpectedAuthorizationHeaderValue string
	ExpectedSignedRequest            string
}

func LoadTestCase(folder string) *V4SignerTestCase {
	base := filepath.Base(folder)
	return &V4SignerTestCase{
		RawRequestToSign:                 loadFileContents(filepath.Join(folder, base+".req")),
		ExpectedCanonicalRequest:         loadFileContents(filepath.Join(folder, base+".creq")),
		ExpectedStringToSign:             loadFileContents(filepath.Join(folder, base+".sts")),
		ExpectedAuthorizationHeaderValue: loadFileContents(filepath.Join(folder, base+".authz")),
		ExpectedSignedRequest:            loadFileContents(filepath.Join(folder, base+".sreq")),
	}
}
func loadFileContents(path string) string {
	all, _ := ioutil.ReadFile(path)
	return string(all)
}

func (this *V4SignerTestCase) IsValid() bool {
	return len(this.RawRequestToSign) > 0 &&
		len(this.ExpectedCanonicalRequest) > 0 &&
		len(this.ExpectedStringToSign) > 0 &&
		len(this.ExpectedAuthorizationHeaderValue) > 0 &&
		len(this.ExpectedSignedRequest) > 0
}

func (this *V4SignerTestCase) buildRequest() (request *http.Request, bodyDigest string) {
	lines := strings.Split(this.RawRequestToSign, "\n")
	firstLine := strings.TrimSuffix(lines[0], " HTTP/1.1")
	space := strings.Index(firstLine, " ")
	method := firstLine[:space]
	address, _ := url.Parse(firstLine[space+1:])
	body := gatherBody(lines[1:])
	bodyReader := strings.NewReader(body)
	request = httptest.NewRequest(method, address.String(), bodyReader)
	request.Header = gatherHeaders(lines[1:])
	request.URL.Host = request.Header.Get("Host")
	return request, hashSHA256([]byte(body))
}

func gatherBody(lines []string) string {
	endOfHeaders := 0
	for _, line := range lines {
		endOfHeaders++
		if len(line) == 0 {
			break
		}
	}
	return strings.Join(lines[endOfHeaders:], "\n")
}

func gatherHeaders(lines []string) http.Header {
	header := make(http.Header)
	for _, line := range lines {
		if len(line) == 0 {
			break
		}
		colon := strings.Index(line, ":")
		key := line[:colon]
		value := line[colon+1:]
		header.Add(key, value)
	}
	return header
}
