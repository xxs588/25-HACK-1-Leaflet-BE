# AI 内容审核功能说明

## 功能概述

本项目集成了硅基流动（SiliconFlow）的 AI 内容审核功能，用于自动审核用户上传的内容，确保平台内容的合规性和安全性。

## 审核范围

AI 审核功能已集成到以下接口：

-   `UploadProblem` - 上传问题时审核问题内容
-   `UpdateProblem` - 修改问题时审核新内容
-   `ChangeProblem` - 修改问题内容时审核新内容
-   `UploadSolve` - 上传解决方案时审核解决方案内容

## 环境配置

在 `.env` 文件中需要配置以下参数：

```env
# 硅基流动API配置
SILICONFLOW_API_KEY=your_api_key_here
SILICONFLOW_BASE_URL=https://api.siliconflow.cn/v1
SILICONFLOW_MODEL=Qwen/Qwen2.5-72B-Instruct

# 审核配置
CONTENT_MODERATION_ENABLED=true
MAX_CONTENT_LENGTH=10000
IMAGE_MAX_SIZE_MB=10
```

## 审核规则

AI 会检查以下方面的内容合规性：

1. 违法、暴力、色情、政治敏感等不当内容
2. 人身攻击、歧视性言论
3. 垃圾信息、广告等无关内容
4. 社区规范和道德标准

## 审核流程

1. 用户提交内容
2. 系统调用 AI 审核服务
3. AI 分析内容并返回审核结果
4. 根据审核结果决定是否允许内容发布

## 审核结果格式

```json
{
  "is_approved": true/false,
  "reason": "审核通过/拒绝的具体原因",
  "confidence": 0.0-1.0
}
```

## 错误处理

-   如果审核服务不可用，系统会拒绝内容上传并返回错误信息
-   如果 API 调用失败，系统会记录错误日志并拒绝内容上传
-   所有审核操作都会记录详细的日志信息

## 日志记录

系统会记录以下审核相关日志：

-   内容审核成功/失败
-   审核拒绝的原因和置信度
-   API 调用错误信息
-   用户操作记录

## 使用示例

### 正常内容（通过审核）

```json
{
    "context": "今天天气真好，我想去公园散步"
}
```

### 违规内容（拒绝审核）

```json
{
    "context": "垃圾信息广告推广"
}
```

## 响应示例

### 审核通过

```json
{
    "message": "上传成功"
}
```

### 审核拒绝

```json
{
    "error": "内容未通过审核",
    "reason": "内容包含垃圾信息",
    "confidence": 0.95
}
```

## 注意事项

1. 确保硅基流动 API 密钥有效且有足够的配额
2. 审核功能默认启用，可通过环境变量关闭
3. 内容长度超过限制会被自动拒绝
4. 建议定期检查审核日志，优化审核规则

## 技术实现

-   服务文件：`service/moderation.go`
-   集成位置：`controller/communication.go`
-   依赖：HTTP 客户端、JSON 解析、环境变量管理

## 性能考虑

-   API 调用超时设置为 30 秒
-   使用连接池优化 HTTP 请求
-   异步日志记录不影响主流程
-   失败时采用保守策略（拒绝内容）
