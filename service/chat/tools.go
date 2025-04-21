package chat

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
	"log"
	"meetingagent/dao"
	"time"
)

type saveMemoryReq struct {
	KeyInfo string `json:"key_info" jsonschema:"description=key information for conversation"`
}

type saveMemoryResp struct {
	status string
}

func GetSaveMemoryTool() tool.InvokableTool {
	saveMemoryTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "save_memory",
			Desc: "提炼出与用户的重要回忆存入向量数据库",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"key_info": {
					Type:     schema.String,
					Desc:     "关于对话的关键信息",
					Required: true,
				},
			}),
		}, func(ctx context.Context, input *saveMemoryReq) (output *saveMemoryResp, err error) {
			// 获取当前时间
			currentTime := time.Now()

			// 获取时间戳（秒级别）
			timestamp := currentTime.Unix()

			// 将时间戳转换为字符串
			timestampStr := fmt.Sprintf("%d", timestamp)
			doc := chromem.Document{
				ID:      uuid.New().String(),
				Content: input.KeyInfo,
				Metadata: map[string]string{
					"timestamp": timestampStr,
					"cate":      "长期记忆",
				},
			}
			collection := dao.ChromemDB.GetCollection("memory", nil)
			err = collection.AddDocument(ctx, doc)
			if err != nil {
				log.Print("添加文档失败")
				return &saveMemoryResp{status: "fail to save"}, err
			}
			log.Print("ai调用了save工具,储存了：" + input.KeyInfo)
			return &saveMemoryResp{status: "ok"}, nil
		},
	)
	return saveMemoryTool
}
