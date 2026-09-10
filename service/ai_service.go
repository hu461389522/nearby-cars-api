package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatWithAI 调用大模型API
func ChatWithAI(question string) (string, error) {
	// ⚠️ 换成你新创建的 DeepSeek Key
	apiKey := "sk-3f054cf477004de6a923c3f6e7342534"

	reqBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个专业的客服助手，回答简洁友好。"},
			{"role": "user", "content": question},
		},
		"temperature": 0.7,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("API返回异常: %s", string(body))
	}
	msg := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	return msg["content"].(string), nil
}