package bcrypt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"backend/pkg/utils"
)

// customImpl .
type customImpl struct {
	secretKey string
}

// NewCustom use a randomly generated encryption key to perform an XOR operation
func NewCustom(secretKey string) BCrypto {
	return &customImpl{secretKey: secretKey}
}

// Encrypt string
func (c *customImpl) Encrypt(s string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成随机的十六进制 encryptKey (区别于 secretKey)
	encryptKey := fmt.Sprintf("%x", r.Intn(32000))

	ctr := 0
	var bytesBuilder bytes.Buffer
	for i := 0; i < len(s); i++ {
		if ctr == len(encryptKey) {
			ctr = 0
		}
		bytesBuilder.WriteByte(encryptKey[ctr])
		bytesBuilder.WriteByte(s[i] ^ encryptKey[ctr])
		ctr++
	}

	// 使用 passportKey 加密
	return base64.StdEncoding.EncodeToString(c.passportKey(bytesBuilder.Bytes(), c.secretKey)), nil
}

// passportKey handle txt with encryptKey
func (c *customImpl) passportKey(txt []byte, encryptKey string) []byte {
	encryptedKey := utils.Md5(encryptKey)
	ctr := 0
	var bytesBuilder bytes.Buffer
	for i := 0; i < len(txt); i++ {
		ctr = ctr % len(encryptedKey)
		bytesBuilder.WriteByte(txt[i] ^ encryptedKey[ctr])
		ctr++
	}
	return bytesBuilder.Bytes()
}

// Decrypt payload
func (c *customImpl) Decrypt(payload string) (originStr string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("recover: %v", e)
		}
	}()

	var decodeBytes []byte
	decodeBytes, err = base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %v", err)
	}

	// 使用 passportKey 解密
	tmp := c.passportKey(decodeBytes, c.secretKey)
	var bytesBuilder bytes.Buffer
	for i := 0; i < len(tmp); i++ {
		txtI0 := i
		txtI := i + 1
		if txtI < len(tmp) && txtI0 < len(tmp) {
			bytesBuilder.WriteByte(tmp[txtI] ^ tmp[txtI0])
			i++
		} else {
			bytesBuilder.WriteByte(0)
		}
	}

	originStr = bytesBuilder.String()
	return originStr, nil
}
