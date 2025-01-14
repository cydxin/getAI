package getAi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// AIClient is the interface for the AI client.
type AIClient interface {
	SetResponse(body interface{})
	SetAPIKey(key string)
	SetURL(url string)
	SendRequest(request interface{}) error
}

// GPTRequest 是已经预设好的chatGpt返回值结构体
type GPTRequest struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
	Stream   bool                `json:"stream"`
}

// GPTResponse 是已经预设好的chatGpt返回值结构体
type GPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AiClientImpl 是 AIClient 接口的实现
type AiClientImpl struct {
	apiKey   string
	url      string
	response interface{}
}

// NewAiClient 创建一个新的 AiClientImpl 实例
func NewAiClient() AIClient {
	return &AiClientImpl{}
}

// SetResponse 设置返回值
func (c *AiClientImpl) SetResponse(body interface{}) {
	if body == nil {
		body = GPTResponse{}
	}
	c.response = body
}

// SetAPIKey 设置 API Key
func (c *AiClientImpl) SetAPIKey(key string) {
	c.apiKey = key
}

// SetURL 设置请求的 URL
func (c *AiClientImpl) SetURL(url string) {
	c.url = url
}

// SendRequest 发送请求并获取响应
func (c *AiClientImpl) SendRequest(requestBody interface{}) error {
	// 将请求体转换为 JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	return json.NewDecoder(resp.Body).Decode(c.response)
}
