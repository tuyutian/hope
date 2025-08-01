// Package utils 用于基础函数的封装
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"backend/pkg/response/message"
)

// Uuid 生成 version4 的uuid
// 返回格式:eba1e8cd0460491049c644bdf3cf024d
func Uuid() string {
	u, err := uuid.NewRandom()
	if err != nil {
		return strings.Replace(RndUuid(), "-", "", -1)
	}

	return strings.Replace(u.String(), "-", "", -1)
}

// RndUuid 基于时间ns和随机数实现唯一的uuid
// 在单台机器上是不会出现重复的uuid
// 如果要在分布式的架构上生成不重复的uuid
// 只需要在rndStr的前面加一些自定义的字符串
// 返回格式:eba1e8cd0460491049c644bdf3cf024d
func RndUuid() string {
	ns := time.Now().UnixNano()
	rndStr := strings.Join([]string{
		strconv.FormatInt(ns, 10), strconv.FormatInt(RandInt64(1000, 9999), 10),
	}, "")
	s := Md5(rndStr)
	s = fmt.Sprintf("%s%s%s%s%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
	return s
}

// RandInt64 生成m-n之间的 int64 随机数
func RandInt64(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if min >= max || min == 0 || max == 0 {
		return max
	}

	return r.Int63n(max-min) + min
}

// Md5 md5 string
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetShopName(shop string) (string, error) {
	if !strings.HasPrefix(shop, "https://") {
		shop = "https://" + shop
	}
	parsedURL, err := url.Parse(shop)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()

	if strings.HasSuffix(host, ".myshopify.com") {
		shopName := strings.TrimSuffix(host, ".myshopify.com")
		return shopName, nil
	}
	return "", message.ErrInvalidAccount
}

func PaseTimeToStamp(timeStr string) int64 {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// 如果解析失败，返回 0
		return 0
	}
	// 返回 Unix 时间戳（秒级）
	return t.Unix()
}

func GetIdFromShopifyGraphqlId(gid string) int64 {
	if gid == "" {
		return 0
	}

	var idStr string
	if strings.HasPrefix(gid, "gid://shopify/") {
		parts := strings.Split(gid, "/")
		if len(parts) > 0 {
			idStr = parts[len(parts)-1]
		}
	} else {
		idStr = gid
	}

	// 将字符串转换为 int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0 // 转换失败时返回 0
	}

	return id
}

// ParseShopifyTime 解析 Shopify 时间格式
func ParseShopifyTime(timeStr string) int64 {
	if timeStr == "" {
		return 0
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0
	}

	return t.Unix()
}
