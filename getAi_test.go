package getAi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSendRequest 测试 SendRequest 方法
func TestSendRequest(t *testing.T) {
	// 使用httptest来模拟一个测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header 'Bearer test-api-key', got '%s'", r.Header.Get("Authorization"))
		}

		// 模拟响应
		response := GPTResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Content string `json:"content"`
					}{
						Content: "Hello, GPT!",
					},
				},
			},
		}

		// 将响应编码为 JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建 AiClientImpl 实例
	client := &AiClientImpl{
		apiKey: "test-api-key", // 不想用测试的 可以把这换成AI的key
		url:    server.URL,     // 不想用测试的 可以把这换成AI的地址
	}
	type GPTRequest struct {
		Model    string              `json:"model"`
		Messages []map[string]string `json:"messages"`
		Stream   bool                `json:"stream"`
	}
	// 设置请求体
	requestBody := GPTRequest{
		Model:    "gpt-4o",
		Messages: []map[string]string{},
		Stream:   false,
	}

	// 设置响应体
	var response GPTResponse
	client.SetResponse(&response)

	// 发送请求
	err := client.SendRequest(&requestBody)
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}

	// 验证响应
	if len(response.Choices) == 0 || response.Choices[0].Message.Content != "Hello, GPT!" {
		t.Errorf("Expected response content 'Hello, GPT!', got '%s'", response.Choices[0].Message.Content)
	}
}
