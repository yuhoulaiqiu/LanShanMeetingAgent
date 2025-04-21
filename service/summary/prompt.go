package summary

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewSummaryTemplate() *prompt.DefaultChatTemplate {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个智能会议总结助手，帮助用户总结会议内容,用markdown的形式返回，你需要阅读已总结的内容，并把未总结的内容进行总结，返回的内容需要包含已总结的内容和未总结的内容的总结，不要省略掉已经总结的内容"),
		schema.SystemMessage("你只需要返回markdown格式的内容，不需要其他任何内容"),
		schema.SystemMessage("现在的时间是: {current_time}"),
		schema.SystemMessage("已经总结过的文本内容: {summarized_text}"),
		schema.UserMessage("需要你总结的文本内容: {unsummarized_text}"),
	)
	return template
}
