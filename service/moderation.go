package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// SiliconFlowAPI 硅基流动API结构
type SiliconFlowAPI struct {
	BaseURL string
	APIKey  string
	Model   string
}

// ModerationRequest 审核请求结构
type ModerationRequest struct {
	Model    string `json:"model"`
	Messages []Message `json:"messages"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ModerationResponse 审核响应结构
type ModerationResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice 选择结构
type Choice struct {
	Message Message `json:"message"`
}

// ModerationResult 审核结果
type ModerationResult struct {
	IsApproved bool   `json:"is_approved"`
	Reason     string `json:"reason"`
	Confidence float64 `json:"confidence"`
}

// NewModerationService 创建新的审核服务
func NewModerationService() *SiliconFlowAPI {
	return &SiliconFlowAPI{
		BaseURL: os.Getenv("SILICONFLOW_BASE_URL"),
		APIKey:  os.Getenv("SILICONFLOW_API_KEY"),
		Model:   os.Getenv("SILICONFLOW_MODEL"),
	}
}

// ModerateContent 审核内容
func (api *SiliconFlowAPI) ModerateContent(content string) (*ModerationResult, error) {
	// 检查是否启用审核
	enabledStr := os.Getenv("CONTENT_MODERATION_ENABLED")
	if enabledStr == "" {
		enabledStr = "true" // 默认启用
	}
	
	enabled, err := strconv.ParseBool(enabledStr)
	if err != nil || !enabled {
		// 如果未启用审核，直接通过
		return &ModerationResult{
			IsApproved: true,
			Reason:     "审核功能未启用",
			Confidence: 1.0,
		}, nil
	}

	// 检查内容长度
	maxLengthStr := os.Getenv("MAX_CONTENT_LENGTH")
	maxLength := 10000 // 默认最大长度
	if maxLengthStr != "" {
		if length, err := strconv.Atoi(maxLengthStr); err == nil {
			maxLength = length
		}
	}
	
	if len(content) > maxLength {
		return &ModerationResult{
			IsApproved: false,
			Reason:     fmt.Sprintf("内容长度超过限制，最大允许%d字符", maxLength),
			Confidence: 1.0,
		}, nil
	}

	// 构建审核提示词
	prompt := fmt.Sprintf(`请审核以下内容是否合规。请检查以下方面：
1. 是否包含违法、暴力、色情、政治敏感等不当内容
2. 是否包含人身攻击、歧视性言论
3. 是否包含垃圾信息、广告等无关内容
4. 是否符合社区规范和道德标准

请以JSON格式返回审核结果，格式如下：
{
  "is_approved": true/false,
  "reason": "审核通过/拒绝的具体原因",
  "confidence": 0.0-1.0
}

待审核内容：%s`, content)

	// 构建请求
	req := ModerationRequest{
		Model: api.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 发送请求
	result, err := api.callAPI(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"content": content[:min(100, len(content))], // 只记录前100字符
		}).Error("AI审核API调用失败")
		
		// API调用失败时，为了不影响用户体验，可以选择放行或拒绝
		// 这里选择保守策略：拒绝
		return &ModerationResult{
			IsApproved: false,
			Reason:     "审核服务暂时不可用，请稍后再试",
			Confidence: 0.0,
		}, nil
	}

	// 记录审核结果日志
	if !result.IsApproved {
		logrus.WithFields(logrus.Fields{
			"content_preview": content[:min(200, len(content))], // 记录前200字符用于分析
			"reason":          result.Reason,
			"confidence":      result.Confidence,
			"content_length":  len(content),
		}).Warn("内容审核未通过")
	} else {
		logrus.WithFields(logrus.Fields{
			"content_preview": content[:min(100, len(content))], // 记录前100字符
			"confidence":      result.Confidence,
			"content_length":  len(content),
		}).Info("内容审核通过")
	}

	return result, nil
}

// callAPI 调用硅基流动API
func (api *SiliconFlowAPI) callAPI(req ModerationRequest) (*ModerationResult, error) {
	// 序列化请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/chat/completions", api.BaseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+api.APIKey)

	// 设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var apiResp ModerationResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("API响应中没有选择项")
	}

	// 解析AI返回的JSON结果
	var result ModerationResult
	if err := json.Unmarshal([]byte(apiResp.Choices[0].Message.Content), &result); err != nil {
		// 如果解析失败，可能是AI没有返回标准格式，进行简单判断
		content := apiResp.Choices[0].Message.Content
		if containsKeywords(content, []string{"不合规", "违规", "拒绝", "不通过", "禁止"}) {
			return &ModerationResult{
				IsApproved: false,
				Reason:     "内容可能包含不当信息",
				Confidence: 0.5,
			}, nil
		}
		
		// 默认通过
		return &ModerationResult{
			IsApproved: true,
			Reason:     "内容审核通过",
			Confidence: 0.8,
		}, nil
	}

	return &result, nil
}

// containsKeywords 检查文本是否包含关键词
func containsKeywords(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if contains(text, keyword) {
			return true
		}
	}
	return false
}

// contains 简单的字符串包含检查
func contains(text, keyword string) bool {
	return len(text) >= len(keyword) && 
		   (text == keyword || 
		    len(text) > len(keyword) && 
		    (text[:len(keyword)] == keyword || 
		     text[len(text)-len(keyword):] == keyword ||
		     containsSubstring(text, keyword)))
}

// containsSubstring 检查子字符串
func containsSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}