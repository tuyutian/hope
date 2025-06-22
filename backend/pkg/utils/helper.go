// Package utils for some common function.
package utils

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetValueFromHTTPCtx 从请求上下文获取指定的key值
func GetValueFromHTTPCtx(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

// SetValueToHTTPCtx 将指定的key/val设置到上下文中
func SetValueToHTTPCtx(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}

// GetStringByCtx 从上下文上获取string的key
func GetStringByCtx(ctx context.Context, key string) string {
	val := GetContextValue(ctx, key)
	if val == nil {
		return ""
	}

	str, ok := val.(string)
	if !ok {
		return ""
	}

	return str
}

// GetContextValue 从ctx上获得key/val
func GetContextValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

// SetContextValue 设置key/val到标准的上下文中
func SetContextValue(ctx context.Context,
	key interface{}, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}
func CallWilding(error string) {
	secret := "SEC3283300ab92509db664438f325b85331f316580e1c6691d442ebd827fea5a504"
	timestamp := time.Now().UnixMilli()
	// 将timestamp和secret拼接成签名字符串
	signStr := strconv.FormatInt(timestamp, 10) + "\n" + secret
	// 使用HmacSHA256算法计算签名
	hmacSha256 := hmac.New(sha256.New, []byte(secret))
	hmacSha256.Write([]byte(signStr))
	signBytes := hmacSha256.Sum(nil)
	// 进行Base64 encode
	signBase64 := base64.StdEncoding.EncodeToString(signBytes)
	// 进行urlEncode
	signUrlEncode := url.QueryEscape(signBase64)

	// 发送markdown消息
	urlP := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=a843ec588eba6834b25f0dd8125f4a28503d431393041cbfea653a7d687e6ccf&timestamp=%d&sign=%s", timestamp, signUrlEncode)
	requestBody := fmt.Sprintf(`{"msgtype": "text","text": {"content":"%s"}}`, error)
	var jsonStr = []byte(requestBody)
	req, err := http.NewRequest("POST", urlP, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
