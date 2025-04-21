package mcp

import (
	"context"
	"fmt"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetStudioTool(ctx context.Context, command string, env []string, args []string) []tool.BaseTool {
	cli, err := client.NewStdioMCPClient(command, env, args...)
	if err != nil {
		fmt.Println("mcp加载出错")
		return nil
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		fmt.Println("mcp初始化出错")
		return nil
	}
	//res, err := cli.ListTools(ctx, mcp.ListToolsRequest{})

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
	if err != nil {
		fmt.Println("eino转换mcp出错")
		return nil
	}
	return tools
}
