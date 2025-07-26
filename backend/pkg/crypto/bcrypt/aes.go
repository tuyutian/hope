package bcrypt

import (
	"strings"

	"backend/pkg/crypto"
)

// EncryptedPrefix encrypt prefix
const EncryptedPrefix = "encrypt_"

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

	return strings.Join([]string{EncryptedPrefix, payload}, ""), err
}

// Decrypt returns decrypt payload
func (c *aesImpl) Decrypt(payload string) (string, error) {
	payload = strings.TrimPrefix(payload, EncryptedPrefix)
	return crypto.AesDecrypt(payload, c.key, c.iv)
}
