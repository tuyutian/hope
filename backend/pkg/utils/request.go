package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HTTPClient 封装的 HTTP 客户端
type HTTPClient struct {
	client  *http.Client
	retries int
	timeout time.Duration
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		retries: 3,
		timeout: time.Second * 30,
	}
}

// WithRetries 设置重试次数
func (c *HTTPClient) WithRetries(retries int) *HTTPClient {
	c.retries = retries
	return c
}

// WithTimeout 设置超时时间
func (c *HTTPClient) WithTimeout(timeout time.Duration) *HTTPClient {
	c.timeout = timeout
	c.client.Timeout = timeout
	return c
}

// PostJSON 发送 JSON POST 请求并支持重试
func (c *HTTPClient) PostJSON(ctx context.Context, url string, payload interface{}, result interface{}) error {
	return c.requestWithRetry(ctx, "POST", url, payload, result)
}

// requestWithRetry 带重试机制的请求方法
func (c *HTTPClient) requestWithRetry(ctx context.Context, method, url string, payload interface{}, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			// 指数退避策略
			backoff := time.Duration(attempt) * time.Second
			log.Println(ctx, "HTTP request retry", "attempt", attempt, "backoff", backoff, "url", url)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := c.doRequest(ctx, method, url, payload, result)
		if err == nil {
			return nil
		}

		lastErr = err

		// 判断是否应该重试
		if !c.shouldRetry(err) {
			break
		}

		log.Println(ctx, "HTTP request failed, retrying", "attempt", attempt+1, "error", err, "url", url)
	}

	return fmt.Errorf("request failed after %d attempts: %w", c.retries+1, lastErr)
}

// doRequest 执行单次请求
func (c *HTTPClient) doRequest(ctx context.Context, method, url string, payload interface{}, result interface{}) error {
	var body io.Reader

	// 处理请求体
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Shopify-App/1.0")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
			URL:        url,
		}
	}

	// 解析响应
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// shouldRetry 判断是否应该重试
func (c *HTTPClient) shouldRetry(err error) bool {
	if httpErr, ok := err.(*HTTPError); ok {
		// 5xx 错误和特定的 4xx 错误可以重试
		return httpErr.StatusCode >= 500 ||
			httpErr.StatusCode == 408 || // Request Timeout
			httpErr.StatusCode == 429 // Too Many Requests
	}

	// 网络错误可以重试
	return true
}

// HTTPError HTTP错误类型
type HTTPError struct {
	StatusCode int
	Body       string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s (URL: %s)", e.StatusCode, e.Body, e.URL)
}

// IsHTTPError 判断是否为HTTP错误
func IsHTTPError(err error) (*HTTPError, bool) {
	httpErr, ok := err.(*HTTPError)
	return httpErr, ok
}
