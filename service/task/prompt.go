package task

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewSummaryTemplate() *prompt.DefaultChatTemplate {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个会议Todo记录助手,你需要从自然语言中提取那些明确提出的todo任务 使用工具创建新的todo,如果没有明确提出的todo,直接返回没有"),
		schema.SystemMessage("一定要是明确提出的需要做的事情才能算作todo,如果是模糊的或者不明确的事情,请不要记录"),
		schema.SystemMessage("你还可以通过更新操作整合或者更新已有的todo,避免todo数量过多,优先考虑更新已有的todo"),
		schema.SystemMessage("现在的时间是: {current_time}"),
		schema.SystemMessage("这个会议的会议ID是: {meeting_id}"),
		schema.SystemMessage("这是已经记录的json: {todo_list}"),
		schema.UserMessage("需要你提取的文本内容: {unsummarized_text}"),
	)
	return template
}
