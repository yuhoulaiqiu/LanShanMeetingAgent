package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"meetingagent/models"
	utils1 "meetingagent/service/utils"
)

type createTodoReq struct {
	MeetingID string `json:"meeting_id" jsonschema:"description=会议ID"`
	Assignee  string `json:"assignee" jsonschema:"description=负责人名称"`
	Task      string `json:"task" jsonschema:"description=任务内容"`
	Level     string `json:"level" jsonschema:"description=任务优先级"`
	State     string `json:"state" jsonschema:"description=任务状态"`
	Deadline  string `json:"deadline" jsonschema:"description=任务截止日期,格式为YYYY-MM-DD-HH:MM:SS"`
}
type createTodoResp struct {
	Status string `json:"status" jsonschema:"description=操作状态"`
}

type UpdateTodoReq struct {
	TodoID    string `json:"todo_id" jsonschema:"description=任务ID"`
	MeetingID string `json:"meeting_id" jsonschema:"description=会议ID"`
	Assignee  string `json:"assignee" jsonschema:"description=负责人名称"`
	Task      string `json:"task" jsonschema:"description=任务内容"`
	Level     string `json:"level" jsonschema:"description=任务优先级"`
	State     string `json:"state" jsonschema:"description=任务状态"`
	Deadline  string `json:"deadline" jsonschema:"description=任务截止日期,格式为YYYY-MM-DD-HH:MM:SS"`
}
type UpdateTodoResp struct {
	Status string `json:"status" jsonschema:"description=操作状态"`
}

type GetAllTodoReq struct {
	MeetingID string `json:"meeting_id" jsonschema:"description=会议ID"`
}
type GetAllTodoResp struct {
	TodoList string `json:"todo_list" jsonschema:"description=todo列表"`
}
type DeleteTodoReq struct {
	TodoID    string `json:"todo_id" jsonschema:"description=任务ID"`
	MeetingID string `json:"meeting_id" jsonschema:"description=会议ID"`
}
type DeleteTodoResp struct {
	Status string `json:"status" jsonschema:"description=操作状态"`
}

func DeleteTodo() tool.InvokableTool {
	updateTodoTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "delete_todo",
			Desc: "删除一条todo",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"todo_id": {
					Type:     schema.String,
					Desc:     "任务ID",
					Required: true,
				},
				"meeting_id": {
					Type:     schema.String,
					Desc:     "会议ID",
					Required: true,
				},
			}),
		}, func(ctx context.Context, input *DeleteTodoReq) (output DeleteTodoResp, err error) {
			err = utils1.DeleteMeetingActionItem(input.MeetingID, input.TodoID)
			if err != nil {
				println("删除todo失败：", err)
				output.Status = "删除todo失败"
				return output, err
			}
			output.Status = "删除todo成功"
			fmt.Println("删除的todo：", input.TodoID)
			return output, nil
		},
	)
	return updateTodoTool
}
func GetAllTodo() tool.InvokableTool {
	updateTodoTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "get_all_todo",
			Desc: "获取所有的todo",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"meeting_id": {
					Type:     schema.String,
					Desc:     "会议ID",
					Required: true,
				},
			}),
		}, func(ctx context.Context, input *GetAllTodoReq) (output GetAllTodoResp, err error) {
			meetingID := input.MeetingID
			todos, err := utils1.ReadMeetingActionItems(meetingID)
			if err != nil {
				fmt.Println("读取会议待办事项失败:", err)
			}
			jsonBytes, err := json.Marshal(todos)
			if err != nil {
				fmt.Errorf("序列化结构体失败: %v", err)
			}
			todoString := string(jsonBytes)
			output.TodoList = todoString
			return output, nil
		},
	)
	return updateTodoTool
}

func UpdateTodo() tool.InvokableTool {
	updateTodoTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "update_todo",
			Desc: "更新一条todo",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"todo_id": {
					Type:     schema.String,
					Desc:     "任务ID",
					Required: true,
				},
				"meeting_id": {
					Type:     schema.String,
					Desc:     "会议ID",
					Required: true,
				},
				"assignee": {
					Type:     schema.String,
					Desc:     "负责人的名称",
					Required: true,
				},
				"task": {
					Type:     schema.String,
					Desc:     "任务的内容",
					Required: true,
				},
				"level": {
					Type:     schema.String,
					Desc:     "任务的优先级",
					Required: true,
				},
				"state": {
					Type:     schema.String,
					Desc:     "任务的状态",
					Required: true,
				},
				"deadline": {
					Type:     schema.String,
					Desc:     "任务的截止日期",
					Required: true,
				},
			}),
		}, func(ctx context.Context, input *UpdateTodoReq) (output UpdateTodoResp, err error) {
			mid := models.ActionItem{
				Assignee: input.Assignee,
				Task:     input.Task,
				Level:    input.Level,
				State:    input.State,
				Deadline: input.Deadline,
			}
			fmt.Println("更新的todo：", mid)
			err = utils1.UpdateMeetingActionItemByID(input.MeetingID, input.TodoID, mid)
			if err != nil {
				println("更新todo失败：", err)
			}
			return output, nil
		},
	)
	return updateTodoTool
}

func CreateTodo() tool.InvokableTool {
	saveMemoryTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "create_todo",
			Desc: "创建一个新的todo",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"meeting_id": {
					Type:     schema.String,
					Desc:     "会议ID",
					Required: true,
				},
				"assignee": {
					Type:     schema.String,
					Desc:     "负责人的名称",
					Required: true,
				},
				"task": {
					Type:     schema.String,
					Desc:     "任务的内容",
					Required: true,
				},
				"level": {
					Type:     schema.String,
					Desc:     "任务的优先级",
					Required: true,
				},
				"state": {
					Type:     schema.String,
					Desc:     "任务的状态",
					Required: true,
				},
				"deadline": {
					Type:     schema.String,
					Desc:     "任务的截止日期",
					Required: true,
				},
			}),
		}, func(ctx context.Context, input *createTodoReq) (output createTodoResp, err error) {
			mid := models.ActionItem{
				Assignee: input.Assignee,
				Task:     input.Task,
				Level:    input.Level,
				State:    input.State,
				Deadline: input.Deadline,
			}
			fmt.Println("创建的todo：", mid)
			err = utils1.AddMeetingActionItem(input.MeetingID, mid)
			if err != nil {
				println("添加todo失败：", err)
			}
			return output, nil
		},
	)
	return saveMemoryTool
}
