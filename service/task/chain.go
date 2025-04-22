package task

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

var TodoApp compose.Runnable[map[string]any, *schema.Message]
var TodoChain *compose.Chain[map[string]any, *schema.Message]

func InitTodoChain() {
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
	tools := []tool.BaseTool{}
	createTodoTool := CreateTodo()
	updateTodoTool := UpdateTodo()
	deleteTodoTool := DeleteTodo()
	tools = append(tools, createTodoTool)
	tools = append(tools, updateTodoTool)
	tools = append(tools, deleteTodoTool)
	toolsConfig := compose.ToolsNodeConfig{
		Tools: tools,
	}
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model:       chatModel,
		ToolsConfig: toolsConfig,
	})
	if err != nil {
		log.Println("创建Agent失败:", err)
	}

	agentLambda, _ := compose.AnyLambda(agent.Generate, agent.Stream, nil, nil)

	template := NewSummaryTemplate()

	//store := Store()
	load := LoadText()
	TodoChain = compose.NewChain[map[string]any, *schema.Message]()
	TodoChain.
		AppendLambda(load).
		AppendChatTemplate(template).
		AppendLambda(agentLambda)

	TodoApp, err = TodoChain.Compile(ctx)
	if err != nil {
		log.Println("创建Chain失败:", err)
	} else {
		log.Println("创建Chain成功")
	}
}
