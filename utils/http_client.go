package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get 发送GET请求
func (c *HTTPClient) Get(url string, params url.Values) (string, error) {
	// 拼接URL参数
	if len(params) > 0 {
		if strings.Contains(url, "?") {
			url += "&" + params.Encode()
		} else {
			url += "?" + params.Encode()
		}
	}

	// 发送请求
	resp, err := c.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}

// Post 发送POST请求
func (c *HTTPClient) Post(url string, params url.Values) (string, error) {
	fmt.Println("发送请求:", url, params.Encode())

	// 发送请求
	resp, err := c.client.PostForm(url, params)
	if err != nil {
		return "", fmt.Errorf("POST请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}

// PostJSON 发送JSON格式的POST请求
func (c *HTTPClient) PostJSON(url string, jsonData []byte) (string, error) {
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}

// Post 简化版Post请求
func Post(headers map[string]string, targetURL, contentType string) (string, error) {
	// 创建请求对象
	data := url.Values{}
	for key, value := range headers {
		data.Add(key, value)
	}

	// 编码请求数据
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求对象失败: %v", err)
	}

	// 设置请求头
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}
