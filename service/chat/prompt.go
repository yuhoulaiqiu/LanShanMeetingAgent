package chat

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewChatTemplate() *prompt.DefaultChatTemplate {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个智能助手，根据会议摘要，回答用户问题"),
		schema.SystemMessage("当用户提到todo时，积极地调用工具去获取，修改，删除todo"),
		schema.SystemMessage("现在的时间是: {current_time}"),
		schema.SystemMessage("现在的会议ID是: {meeting_id}"),
		schema.SystemMessage("从向量数据库中获取的相关上下文: {context}"),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("提问: {question}"),
	)
	return template
}
