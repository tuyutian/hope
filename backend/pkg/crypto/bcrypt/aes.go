package bcrypt

import (
	"backend/pkg/crypto"
)

var _ BCrypto = (*aesImpl)(nil)

// aesImpl aes-256-cbc
type aesImpl struct {
	key string
	iv  string
}

// NewAesBCrypto return BCrypto from aes
// iv必须是16位
// 当key 16位的时候 相当于php openssl_decrypt(base64_decode($strEncode), 'aes-128-cbc', $key, true, $iv)
// 当key 24位的时候 相当于php openssl_decrypt(base64_decode($strEncode), 'aes-192-cbc', $key, true, $iv)
// 当key 32位的时候 相当于php openssl_decrypt(base64_decode($strEncode), ”, $key, true, $iv)
func NewAesBCrypto(key string, iv string) BCrypto {
	return &aesImpl{key: key, iv: iv}
}

// Encrypt returns encrypt string
func (c *aesImpl) Encrypt(s string) (string, error) {
	payload, err := crypto.AesEncrypt(s, c.key, c.iv)
	if err != nil {
		return "", err
	}

	return payload, err
}

// Decrypt returns decrypt payload
func (c *aesImpl) Decrypt(payload string) (string, error) {
	return crypto.AesDecrypt(payload, c.key, c.iv)
}
