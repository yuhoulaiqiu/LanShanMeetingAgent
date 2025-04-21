package chat

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"log"
	"meetingagent/config"
)

var App compose.Runnable[string, *schema.Message]
var ChatChain *compose.Chain[string, *schema.Message]
var Ctx context.Context

func InitChain() {
	Ctx = context.Background()
	chatModel, err := ark.NewChatModel(Ctx, &ark.ChatModelConfig{
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
	memorytool := GetSaveMemoryTool()
	tools := []tool.BaseTool{}
	tools = append(tools, memorytool)
	toolsConfig := compose.ToolsNodeConfig{Tools: tools}
	var agent *react.Agent
	agent, err = react.NewAgent(Ctx, &react.AgentConfig{
		Model:       chatModel,
		ToolsConfig: toolsConfig,
	})
	agentLamda, _ := compose.AnyLambda(agent.Generate, agent.Stream, nil, nil)
	template := NewChatTemplate()
	loadMemoryLamda := InitLoadMemory()
	load := LoadText()
	ChatChain = compose.NewChain[string, *schema.Message]()
	ChatChain.
		AppendLambda(loadMemoryLamda).
		AppendLambda(load).
		AppendChatTemplate(template).
		AppendLambda(agentLamda)

	App, err = ChatChain.Compile(Ctx)
	if err != nil {
		log.Println("创建Chain失败:", err)
	} else {
		log.Println("创建Chain成功")
	}
}

func Stream(input string) (*schema.StreamReader[*schema.Message], error) {
	return App.Stream(Ctx, input)
}
