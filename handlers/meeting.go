package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"meetingagent/service/summary"
	"os"
	"path/filepath"
	"time"

	"meetingagent/models"
	"meetingagent/service/utilus"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sse"
)

// CreateMeeting handles the creation of a new meeting
func CreateMeeting(ctx context.Context, c *app.RequestContext) {
	var reqBody map[string]interface{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	fmt.Printf("create meeting: %s\n", string(jsonBody))

	// 生成唯一的会议 ID
	meetingID := "meeting_" + time.Now().Format("20060102150405")

	// 创建会议对象
	meeting := models.Meeting{
		ID:      meetingID,
		Content: reqBody,
	}

	// 存储会议数据
	utilus.GetMeetingStore().SetMeeting(meetingID, meeting)

	// 返回会议 ID
	response := models.PostMeetingResponse{
		ID: meetingID,
	}
	go summary.SummaryMeeting(meetingID)
	c.JSON(consts.StatusOK, response)
}

// ListMeetings handles listing all meetings
func ListMeetings(ctx context.Context, c *app.RequestContext) {
	// 从内存中获取所有会议
	meetingsList := utilus.GetMeetingStore().GetAllMeetings()

	// 如果没有会议记录，添加一个默认样例
	if len(meetingsList) == 0 {
		meetingsList = append(meetingsList, models.Meeting{
			ID: "meeting_123",
			Content: map[string]interface{}{
				"title":        "Sample Meeting",
				"description":  "因为没有会议记录，所以这里是一个默认的样例",
				"participants": []string{"John Doe", "Jane Smith"},
				"start_time":   "2025-04-20 08:00:00",
				"end_time":     "2025-04-20 09:00:00",
				"content":      "This is the content of the meeting",
			},
		})
	}

	response := models.GetMeetingsResponse{
		Meetings: meetingsList,
	}

	c.JSON(consts.StatusOK, response)
}

// GetMeetingSummary handles retrieving a meeting summary
func GetMeetingSummary(ctx context.Context, c *app.RequestContext) {
	meetingID := c.Query("meeting_id")

	if meetingID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "meeting_id参数必须提供"})
		return
	}
	fmt.Printf("获取会议摘要，meetingID: %s\n", meetingID)

	// 构建文件路径
	filePath := filepath.Join("data", "meetings", meetingID+".json")

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取会议摘要文件失败: %v\n", err)
		c.JSON(consts.StatusNotFound, utils.H{"error": fmt.Sprintf("会议摘要文件不存在: %v", err)})
		return
	}

	// 解析JSON
	var meeting models.SummarizedMeeting
	if err := json.Unmarshal(data, &meeting); err != nil {
		fmt.Printf("解析会议摘要数据失败: %v\n", err)
		c.JSON(consts.StatusInternalServerError, utils.H{"error": fmt.Sprintf("解析会议摘要数据失败: %v", err)})
		return
	}

	// 返回完整的摘要数据
	c.JSON(consts.StatusOK, meeting)
}

// GetOneMeeting handles retrieving a single meeting by ID
func GetOneMeeting(ctx context.Context, c *app.RequestContext) {
	meetingID := c.Query("meeting_id")
	if meetingID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "meeting_id is required"})
		return
	}

	meeting, err := utilus.GetMeetingStore().GetMeeting(meetingID)
	if err != nil {
		c.JSON(consts.StatusNotFound, utils.H{"error": err.Error()})
		return
	}

	c.JSON(consts.StatusOK, meeting)
}

// HandleChat handles the SSE chat session
func HandleChat(ctx context.Context, c *app.RequestContext) {
	meetingID := c.Query("meeting_id")
	sessionID := c.Query("session_id")
	message := c.Query("message")

	if meetingID == "" || sessionID == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "meeting_id and session_id are required"})
		return
	}

	if message == "" {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "message is required"})
		return
	}

	fmt.Printf("meetingID: %s, sessionID: %s, message: %s\n", meetingID, sessionID, message)

	// Set SSE headers
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")

	// Create SSE stream
	stream := sse.NewStream(c)

	// TODO: Implement actual chat logic
	// This is a simple example that sends a message every second
	ticker := time.NewTicker(time.Millisecond * 100)
	stopChan := make(chan struct{})
	go func() {
		time.AfterFunc(time.Second, func() {
			ticker.Stop()
			close(stopChan)
		})
	}()

	msg := fmt.Sprintf("Fake sample chat message: %s\n", time.Now().Format(time.RFC3339))

	for {
		select {
		case <-ticker.C:
			res := models.ChatMessage{
				Data: msg,
			}

			data, err := json.Marshal(res)
			if err != nil {
				return
			}

			event := &sse.Event{
				Data: data,
			}

			if err := stream.Publish(event); err != nil {
				return
			}
		case <-stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
