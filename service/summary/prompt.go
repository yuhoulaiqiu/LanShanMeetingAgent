package summary

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewSummaryTemplate() *prompt.DefaultChatTemplate {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个智能助手，帮助用户总结会议内容"),
		schema.SystemMessage("现在的时间是: {current_time}"),
		schema.SystemMessage("已经总结过的文本内容: {summarized_text}"),
		schema.UserMessage("需要你总结的文本内容: {unsummarized_text}"),
	)
	return template
}
