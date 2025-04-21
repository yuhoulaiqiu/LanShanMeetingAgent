package utils

import (
	"encoding/json"
	"fmt"
	"meetingagent/models"
	"os"
	"path/filepath"
)

// DeleteMeetingActionItem 删除指定会议中的某个待办事项
func DeleteMeetingActionItem(meetingID string, todoID string) error {
	// 获取该会议ID的锁并加写锁
	lock := GetFileLockManager().GetLock(meetingID)
	lock.Lock()
	defer lock.Unlock()

	// 构建文件路径
	filePath := filepath.Join("data", "meetings", meetingID+".json")

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取会议摘要文件失败: %v", err)
	}

	// 解析JSON
	var meeting models.SummarizedMeeting
	if err := json.Unmarshal(data, &meeting); err != nil {
		return fmt.Errorf("解析会议摘要数据失败: %v", err)
	}

	// 查找并删除指定的ActionItem
	found := false
	for i, item := range meeting.Summary.ActionItems {
		if item.TodoID == todoID {
			// 删除该项
			meeting.Summary.ActionItems = append(meeting.Summary.ActionItems[:i], meeting.Summary.ActionItems[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到指定的TodoID: %s", todoID)
	}

	// 转换为JSON
	jsonData, err := json.MarshalIndent(meeting, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化会议数据失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入会议文件失败: %v", err)
	}

	return nil
}
