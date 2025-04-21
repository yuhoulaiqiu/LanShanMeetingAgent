package summary

import (
	"context"
	"encoding/json"
	"fmt"
	"meetingagent/service/task"
	"meetingagent/service/utils"
)

func SummaryMeeting(meetingID string) {
	meeting, err := utils.GetMeetingStore().GetMeeting(meetingID)
	if err != nil {
		fmt.Println("获取会议记录失败:", err)
		return
	}

	// 将会议内容转换为JSON字符串
	contentJSON, err := json.Marshal(meeting.Content)
	if err != nil {
		fmt.Println("序列化会议内容失败:", err)
		return
	}

	// 将会议内容解析为结构体（假设内容是包含contents数组的结构）
	var meetingData struct {
		Contents []struct {
			TimeFrom string `json:"time_from"`
			TimeTo   string `json:"time_to"`
			User     string `json:"user"`
			Content  struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"contents"`
	}

	if err := json.Unmarshal(contentJSON, &meetingData); err != nil {
		fmt.Println("解析会议内容失败:", err)
		return
	}

	// 分片存储结果
	var chunks []string
	currentChunk := ""
	currentRuneCount := 0
	maxRuneCount := 600 // 每个分片的最大字符数

	// 遍历会议内容构建分片
	for _, item := range meetingData.Contents {
		// 构建当前条目的文本
		entryText := fmt.Sprintf("%s (%s-%s): %s\n",
			item.User, item.TimeFrom, item.TimeTo, item.Content.Text)

		// 计算添加这条记录后的总字符数
		entryRunes := []rune(entryText)

		// 如果添加这条记录会超过限制，先保存当前分片
		if currentRuneCount+len(entryRunes) > maxRuneCount && currentRuneCount > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
			currentRuneCount = 0
		}

		// 添加当前记录到分片
		currentChunk += entryText
		currentRuneCount += len(entryRunes)
	}

	// 添加最后一个分片（如果有）
	if currentRuneCount > 0 {
		chunks = append(chunks, currentChunk)
	}

	for i, chunk := range chunks {
		fmt.Printf("Chunk %d:\n", i+1)
		thisTexts := texts{
			unsummarized_text: chunk,
			summarized_text:   "",
		}
		if i == 0 {
			output, err := App.Invoke(context.Background(), thisTexts)
			if err != nil {
				fmt.Println("调用模型失败:", err)
			}
			err = utils.CreateMeetingJSON(output.Content, meetingID)
			if err != nil {
				fmt.Println("创建会议摘要失败:", err)
			}
			todos, err := utils.ReadMeetingActionItems(meetingID)
			if err != nil {
				fmt.Println("读取会议待办事项失败:", err)
			}
			m := make(map[string]any)
			m["unsummarized_text"] = chunk
			jsonBytes, err := json.Marshal(todos)
			if err != nil {
				fmt.Errorf("序列化结构体失败: %v", err)
			}
			todoString := string(jsonBytes)
			m["todo_list"] = todoString
			m["meeting_id"] = meetingID
			go func() {
				output, err := task.TodoApp.Invoke(context.Background(), m)
				if err != nil {
					fmt.Println("调用todo模型失败:", err)
				}
				println("todo模型返回:", output.Content)
			}()
			continue
		}
		go func() {
			todos, err := utils.ReadMeetingActionItems(meetingID)
			if err != nil {
				fmt.Println("读取会议待办事项失败:", err)
			}
			m := make(map[string]any)
			m["unsummarized_text"] = chunk
			jsonBytes, err := json.Marshal(todos)
			if err != nil {
				fmt.Errorf("序列化结构体失败: %v", err)
			}
			todoString := string(jsonBytes)
			m["todo_list"] = todoString
			m["meeting_id"] = meetingID
			output, err := task.TodoApp.Invoke(context.Background(), m)
			if err != nil {
				fmt.Println("调用todo模型失败:", err)
			}
			println("todo模型返回:", output.Content)
		}()
		thisTexts.summarized_text, err = utils.ReadMeetingSummaryText(meetingID)
		if err != nil {
			fmt.Println("读取会议摘要失败:", err)
		}
		output, err := App.Invoke(context.Background(), thisTexts)
		if err != nil {
			fmt.Println("调用模型失败:", err)
		}
		err = utils.UpdateMeetingSummaryText(meetingID, output.Content)
		if err != nil {
			fmt.Println("更新会议摘要失败:", err)
		}
		err = Split(meetingID)
		if err != nil {
			fmt.Println("分片失败:", err)
		}
	}
}

// 这里可以将分片存储到数据库或其他存储中
