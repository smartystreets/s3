package s3

import (
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCommonFixture(t *testing.T) {
	gunit.Run(new(CommonFixture), t)
}

type CommonFixture struct {
	*gunit.Fixture
}

func (this *CommonFixture) TestHashFunctions() {
	this.So(hashMD5([]byte("Pretend this is a REALLY long byte array...")), should.Equal, "KbVTY8Vl6VccnzQf1AGOFw==")
	this.So(hashSHA256([]byte("This is... Sparta!!")), should.Equal,
		"5c81a4ef1172e89b1a9d575f4cd82f4ed20ea9137e61aa7f1ab936291d24e79a")

	key := []byte("asdf1234")
	contents := "SmartyStreets was here"

	expectedHMAC_SHA256 := []byte{
		65, 46, 186, 78, 2, 155, 71, 104, 49, 37, 5, 66, 195, 129, 159, 227,
		239, 53, 240, 107, 83, 21, 235, 198, 238, 216, 108, 149, 143, 222, 144, 94}
	this.So(hmacSHA256(key, contents), should.Resemble, expectedHMAC_SHA256)

	expectedHMAC_SHA1 := []byte{
		164, 77, 252, 0, 87, 109, 207, 110, 163, 75, 228, 122, 83, 255, 233, 237, 125, 206, 85, 70}
	this.So(hmacSHA1(key, contents), should.Resemble, expectedHMAC_SHA1)
}

func (this *CommonFixture) TestConcat() {
	this.So(join("\n", "Test1", "Test2"), should.Equal, "Test1\nTest2")
	this.So(join(".", "Test1"), should.Equal, "Test1")
	this.So(join("\t", "1", "2", "3", "4"), should.Equal, "1\t2\t3\t4")
}

func (this *CommonFixture) TestURINormalization() {
	this.So(
		normalizeURI("/-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"), should.Equal,
		"/-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	this.So(normalizeURI("/ /foo"), should.Equal, "/%20/foo")
	this.So(normalizeURI("/(foo)"), should.Equal, "/%28foo%29")

	this.So(
		normalizeQuery(url.Values{"p": []string{" +&;-=._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"}}),
		should.Equal,
		"p=%20%2B%26%3B-%3D._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}
