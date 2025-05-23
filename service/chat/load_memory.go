package chat

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/philippgille/chromem-go"
	"log"
	"meetingagent/config"
	"meetingagent/dao"
)

func InitLoadMemory() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		question := input["question"].(string)
		meetingID := input["meeting_id"].(string)
		collection, err := dao.ChromemDB.GetOrCreateCollection("rag", nil, chromem.NewEmbeddingFuncOpenAICompat(
			"https://ark.cn-beijing.volces.com/api/v3",
			config.Cfg.ModelInfo.ApiKey,
			"doubao-embedding-large-text-240915",
			nil,
		))
		if err != nil {
			log.Fatal("fail to get memory", err)
		}
		count := collection.Count()
		var nResult int
		if count >= 10 {
			nResult = 10
		} else {
			if count == 0 {
				nResult = 0
			} else {
				nResult = count
			}
		}
		where := map[string]string{"source_meeting_id": meetingID}
		result, err := collection.Query(ctx, question, nResult, where, nil)
		if err != nil {
			log.Println("fail to get memory", err)
		}
		memories := ""
		for i, item := range result {
			if item.Similarity < 0.4 {
				continue
			}
			memories += item.Content
			memories += fmt.Sprintf("%f", item.Similarity)
			if i != len(result)-1 {
				memories += ","
			}
		}
		log.Print("get memory:" + memories)
		rows, err := dao.SqliteDB.Query("SELECT role,content,timestamp FROM messages ORDER BY timestamp DESC LIMIT 5")
		if err != nil {
			log.Print("fail to get context")
		}
		defer rows.Close()

		var chatHistory []*schema.Message
		for rows.Next() {
			var role, content, timestamp string
			err := rows.Scan(&role, &content, &timestamp)
			if err != nil {
				log.Println("fail to scan", err)
			}
			var msg *schema.Message
			if role == "user" {
				msg = schema.UserMessage(content)
			} else if role == "ai" {
				msg = schema.AssistantMessage(content, nil)
			} else {
				log.Fatal("未知的角色")
			}
			chatHistory = append(chatHistory, msg)

		}
		if err := rows.Err(); err != nil {
			log.Print("转换的过程有误", err)
		}
		output := map[string]any{
			"context":      memories,
			"chat_history": chatHistory,
			"question":     question,
			"meeting_id":   meetingID,
		}
		go InsertMessage(dao.SqliteDB, "user", question)
		return output, nil
	})
}
