// Package utils for some common function.
package utils

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	secret := "O6yErxY6HoHV6Ym5BPZEK"
	timestamp := time.Now().Unix()
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 发送markdown消息
	urlP := "https://open.feishu.cn/open-apis/bot/v2/hook/e3421807-ddaf-4aa2-a3d3-697ef93516bb"
	requestBody := struct {
		Timestamp int64  `json:"timestamp"`
		Sign      string `json:"sign"`
		MsgType   string `json:"msg_type"`
		Content   struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		Timestamp: timestamp,
		Sign:      signature,
		MsgType:   "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: error,
		},
	}
	jsonStr, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", urlP, bytes.NewBuffer(jsonStr))
	if req == nil {
		fmt.Println(err)
		return
	}
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
