package crypto

import (
	"log"
	"testing"
)

/*
=== RUN   TestHmac256

	crypto_test.go:19: dd14fe1a8967127104b65ef52c009422
	crypto_test.go:21: key:  ba10fd13343f53e4
	crypto_test.go:22: 127092b52a8c76b9f35f40ffc33ec3cf

--- PASS: TestHmac256 (0.00s)
PASS
*/
func TestHmac256(t *testing.T) {
	t.Log(Hmac256("123456", ""))
	key := GetIteratorStr(16)
	t.Log("key: ", key)
	t.Log(Hmac256("123456", key))
}

/*
=== RUN   TestSha256

	crypto_test.go:26: test Sha256
	crypto_test.go:27: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92

--- PASS: TestSha256 (0.00s)
PASS
*/
func TestSha256(t *testing.T) {
	t.Log("test Sha256")
	t.Log(Sha256("123456"))
}

// === RUN   TestGetIteratorStr
// 2025/03/24 15:39:27 current str:  a58ede
// 2025/03/24 15:39:27 current str:  4ce165
// 2025/03/24 15:39:27 current str:  88c159
// --- PASS: TestGetIteratorStr (0.00s)
// PASS
func TestGetIteratorStr(t *testing.T) {
	for i := 0; i < 1000; i++ {
		log.Println("current str: ", GetIteratorStr(6))
	}
}

var (
	k  = GetIteratorStr(16)
	iv = GetIteratorStr(16)
)

/*
验证 AesEncrypt CBC 加密
=== RUN   TestCbc256

	crypto_test.go:57: nhqF9cDUOQcOI9Bj4BU4Ng== <nil>

--- PASS: TestCbc256 (0.00s)
PASS
*/
func TestCbc256(t *testing.T) {
	t.Log(AesEncrypt("123456", k, iv))
}

/*
=== RUN   TestDecodeCbc256

	crypto_test.go:73: 123456

--- PASS: TestDecodeCbc256 (0.00s)
PASS
*/
func TestDecodeCbc256(t *testing.T) {
	b, _ := AesEncrypt("123456", k, iv)
	str, _ := AesDecrypt(b, k, iv)
	t.Log(str)
}

/*
* 验证加解密
=== RUN   TestAesEbc
crypto_test.go:54: ebc加密后: 3e75cb8bcd9d5e08
crypto_test.go:57: ebc解密: 123456
--- PASS: TestAesEbc (0.00s)
PASS
*/
func TestAesEbc(t *testing.T) {
	k := GetIteratorStr(8)
	b, _ := EncryptEcb("123456", k)
	t.Log("ebc加密后:", b)

	s, _ := DecryptEcb(b, k)
	t.Log("ebc解密:", s)
}

/*
*
测试aes-256-cbc加密
$ go test -v -test.run=TestAesCbc
=== RUN   TestAesCbc
2025/03/24 15:43:38 /fxQRPGIHJ9AFsG67MSVDvLFSDp+/ZFGkHT+Y46h4jln9IzORfsEhR6L2qh5mDDQ
2025/03/24 15:43:38 HRHtimkjsJktwu6AzH2ji9MP9OLpRBRf35Xcm7zFNmr5Lj8X1rxxJiCIQJqnLC8r
2025/03/24 15:43:38 Sj1ENtUBam7C6PglPZgLZGy9lC8bppce7NS8RExuVa+xWow04Trnlc+kJh+Wz9LL
2025/03/24 15:43:38 中文数字123字母ABC符号!@#$%^&*() <nil>
--- PASS: TestAesCbc (0.00s)
PASS
*/
func TestAesCbc(t *testing.T) {
	str := `中文数字123字母ABC符号!@#$%^&*()`
	k2 := `abcdefghijklmnop`
	iv2 := `1234567890123456`
	s, _ := AesEncrypt(str, k2, iv2)
	log.Println(s)

	// log.Println(AesDecrypt(`/fxQRPGIHJ9AFsG67MSVDvLFSDp+/ZFGkHT+Y46h4jln9IzORfsEhR6L2qh5mDDQ`, k2, iv2))

	k2 = `abcdefghijklmnop1234567890123456`
	iv2 = `1234567890123456`
	s, _ = AesEncrypt(str, k2, iv2)
	log.Println(s)

	k2 = `abcdefghijklmnop12345678`
	iv2 = `1234567890123456`
	s, _ = AesEncrypt(str, k2, iv2)
	log.Println(s)

	log.Println(AesDecrypt(s, k2, iv2))
}
