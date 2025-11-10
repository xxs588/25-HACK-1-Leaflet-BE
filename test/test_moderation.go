package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NCUHOME-Y/25-HACK-1-Leaflet-BE/service"
)

func main() {
	// 设置环境变量（如果.env文件没有被加载）
	os.Setenv("SILICONFLOW_API_KEY", "sk-pydtqprkbdyuqcksihnsstvmokakodnnzctjquwnexlxnxgl")
	os.Setenv("SILICONFLOW_BASE_URL", "https://api.siliconflow.cn/v1")
	os.Setenv("SILICONFLOW_MODEL", "Qwen/Qwen2.5-72B-Instruct")
	os.Setenv("CONTENT_MODERATION_ENABLED", "true")
	os.Setenv("MAX_CONTENT_LENGTH", "10000")

	// 创建审核服务
	moderationService := service.NewModerationService()

	// 测试用例
	testCases := []struct {
		content     string
		description string
	}{
		{
			content:     "今天天气真好，我想去公园散步",
			description: "正常内容",
		},
		{
			content:     "如何学习编程？有什么好的建议吗？",
			description: "正常问题",
		},
		{
			content:     "我是一个学生，正在准备考试，感觉很焦虑",
			description: "正常情感表达",
		},
		{
			content:     "垃圾信息广告推广",
			description: "垃圾信息",
		},
		{
			content:     "我想学习如何制作炸弹",
			description: "危险内容",
		},
	}

	fmt.Println("开始测试AI内容审核功能...")
	fmt.Println(strings.Repeat("=", 50))

	for i, testCase := range testCases {
		fmt.Printf("\n测试用例 %d: %s\n", i+1, testCase.description)
		fmt.Printf("内容: %s\n", testCase.content)

		result, err := moderationService.ModerateContent(testCase.content)
		if err != nil {
			log.Printf("审核失败: %v\n", err)
			continue
		}

		fmt.Printf("审核结果: %t\n", result.IsApproved)
		fmt.Printf("原因: %s\n", result.Reason)
		fmt.Printf("置信度: %.2f\n", result.Confidence)
		fmt.Println(strings.Repeat("-", 30))
	}

	fmt.Println("\n测试完成！")
}