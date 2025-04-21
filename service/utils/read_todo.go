package utils

import (
	"encoding/json"
	"fmt"
	"meetingagent/models"
	"os"
	"path/filepath"
)

// ReadMeetingActionItems 读取指定会议摘要的ActionItems内容
func ReadMeetingActionItems(meetingID string) ([]models.ActionItem, error) {
	// 获取该会议ID的锁并加读锁
	lock := GetFileLockManager().GetLock(meetingID)
	lock.RLock()
	defer lock.RUnlock()

	// 构建文件路径
	filePath := filepath.Join("data", "meetings", meetingID+".json")

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取会议摘要文件失败: %v", err)
	}

	// 解析JSON
	var meeting models.SummarizedMeeting
	if err := json.Unmarshal(data, &meeting); err != nil {
		return nil, fmt.Errorf("解析会议摘要数据失败: %v", err)
	}

	return meeting.Summary.ActionItems, nil
}

// UpdateMeetingActionItems 更新指定会议摘要的ActionItems内容
func UpdateMeetingActionItems(meetingID string, newActionItems []models.ActionItem) error {
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

	// 更新ActionItems内容
	meeting.Summary.ActionItems = newActionItems

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
