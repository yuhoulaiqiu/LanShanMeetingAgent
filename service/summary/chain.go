package summary

import "C"
import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
	"meetingagent/config"
)

type texts struct {
	unsummarized_text string
	summarized_text   string
}

var App compose.Runnable[texts, *schema.Message]
var SummaryChain *compose.Chain[texts, *schema.Message]

func InitChain() {
	ctx := context.Background()

	//chatModel, _ := openai.NewChatModel(ctx, &openai.ChatModelConfig{
	//	APIKey:  config.Cfg.ModelInfo.ApiKey,
	//	BaseURL: config.Cfg.ModelInfo.BaseUrl,
	//	Model:   config.Cfg.ModelInfo.ModelName,
	//
	//	Temperature: &config.Cfg.ModelInfo.Temperature,
	//})
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		// 服务配置
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3", // 服务地址
		Region:  "cn-beijing",                               // 区域
		// 认证配置（二选一）
		APIKey:    config.Cfg.ModelInfo.ApiKey, // API Key 认证
		AccessKey: config.Cfg.VKDB.Ak,          // AK/SK 认证
		SecretKey: config.Cfg.VKDB.Sk,

		// 模型配置
		Model: config.Cfg.ModelInfo.ModelName, // 模型端点 ID

	})
	if err != nil {
		log.Println("创建模型失败:", err)
		return
	}

	//agent, err := react.NewAgent(ctx, &react.AgentConfig{
	//	Model: chatModel,
	//})
	//if err != nil {
	//	log.Println("创建Agent失败:", err)
	//}

	//agentLambda, _ := compose.AnyLambda(agent.Generate, agent.Stream, nil, nil)

	template := NewSummaryTemplate()

	//store := Store()
	load := LoadText()
	SummaryChain = compose.NewChain[texts, *schema.Message]()
	SummaryChain.
		AppendLambda(load).
		AppendChatTemplate(template).
		AppendChatModel(chatModel)
	//	AppendLambda(store)
	//AppendLambda(agentLambda)

	App, err = SummaryChain.Compile(ctx)
	if err != nil {
		log.Println("创建Chain失败:", err)
	} else {
		log.Println("创建Chain成功")
	}
}
