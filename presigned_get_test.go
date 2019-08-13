package s3

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestPresignedGetFixture(t *testing.T) {
	gunit.Run(new(PresignedGetFixture), t)
}

type PresignedGetFixture struct {
	*gunit.Fixture
}

func (this *PresignedGetFixture) TestSignature() {
	address, err := NewPresignedGet(
		Credentials("AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
		Region("us-east-1"),
		Bucket("examplebucket"),
		Key("test.txt"),
		Timestamp(time.Date(2013, 5, 24, 0, 0, 0, 0, time.UTC)),
		ExpiresIn(time.Hour*24),
	)

	this.So(err, should.BeNil)
	this.So(address, should.Equal,
		"https://examplebucket.s3.amazonaws.com/test.txt"+
			"?X-Amz-Algorithm=AWS4-HMAC-SHA256"+
			"&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20130524%2Fus-east-1%2Fs3%2Faws4_request"+
			"&X-Amz-Date=20130524T000000Z"+
			"&X-Amz-Expires=86400"+
			"&X-Amz-SignedHeaders=host"+
			"&X-Amz-Signature=aeeed9bbccd4d02ee5c0109b86d86835f995330da4c265957d157751f604d404",
	)
}
