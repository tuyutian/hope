package bcrypt

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"backend/pkg/utils"
)

var _ BCrypto = (*customNumberImpl)(nil)

var numberRegexp = regexp.MustCompile("^([0-9]+)")

// customNumberImpl .
type customNumberImpl struct {
	secretKey string
	salt      string
}

// NewCustomNumber MD5 and Base64 are used, and custom character replacement
func NewCustomNumber(secretKey string, salt string) BCrypto {
	return &customNumberImpl{
		secretKey: secretKey,
		salt:      salt,
	}
}

// Encrypt passportKey string
func (c *customNumberImpl) Encrypt(s string) (string, error) {
	// 获取锁的长度
	lockKeyRunes := []rune(c.salt)
	lockKeyRunesLen := len(lockKeyRunes)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 从指定字符串中随机获取一个字符
	lockKeyRandom := r.Intn(lockKeyRunesLen - 1)
	lockKeyRandomString := string(lockKeyRunes[lockKeyRandom])

	// 生成随机密钥, 并进行 MD5 加密
	encryptedSecretKey := utils.Md5(c.secretKey + lockKeyRandomString)

	// 将明文加密
	encryptedText := base64.StdEncoding.EncodeToString([]byte(s))

	dataLen := len([]rune(encryptedText))
	secretKeyLen := len([]rune(encryptedSecretKey))

	var j, k int
	var b strings.Builder
	for i := 0; i < dataLen; i++ {
		if k == secretKeyLen {
			k = 0
		}
		n1 := strings.Index(c.salt, string(encryptedText[i]))
		n2 := int(rune(encryptedSecretKey[k]))
		j = (n1 + lockKeyRandom + n2) % lockKeyRunesLen
		b.WriteString(string(c.salt[j]))
		k++
	}

	b.WriteString(lockKeyRandomString)
	return b.String(), nil
}

// Decrypt returns decrypt payload
func (c *customNumberImpl) Decrypt(payload string) (originStr string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("recover: %v", e)
		}
	}()

	// 获取锁的长度
	lockKeyRunes := []rune(c.salt)
	payloadRunes := []rune(payload)
	lockKeyRuneLen := len(lockKeyRunes)
	payloadRuneLen := len(payloadRunes)

	// 从密文中提取一个最后一个字符
	payloadLastChar := string(payloadRunes[payloadRuneLen-1])
	lockKeyIndex := strings.Index(c.salt, payloadLastChar)

	// 生成密钥，并进行 MD5 加密
	encryptedSecretKey := utils.Md5(c.secretKey + payloadLastChar)

	encryptedSecretKeyLen := len([]rune(encryptedSecretKey))

	str := payload[0 : payloadRuneLen-1]

	var j, k int
	var resultBuilder strings.Builder
	for i := 0; i < payloadRuneLen-1; i++ {
		if k == encryptedSecretKeyLen {
			k = 0
		}
		n1 := strings.Index(c.salt, string(str[i]))
		n2 := int(rune(encryptedSecretKey[k]))
		j = n1 - lockKeyIndex - n2
		for {
			if j >= 0 {
				break
			}
			j += lockKeyRuneLen
		}
		resultBuilder.WriteString(string(c.salt[j]))
		k++
	}

	// 解密
	data, err := base64.StdEncoding.DecodeString(resultBuilder.String())
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %v", err)
	}

	if numberStr := numberRegexp.FindSubmatch(data); len(numberStr) != 0 {
		originStr = string(numberStr[0])
		return originStr, nil
	}

	return "0", nil
}
