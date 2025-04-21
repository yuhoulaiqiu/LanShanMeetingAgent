package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"meetingagent/models"
	"os"
	"path/filepath"
)

// generateRandomID 生成指定长度的随机ID
func generateRandomID(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AddMeetingActionItem 向指定会议的ActionItems列表中添加一个新的待办事项
func AddMeetingActionItem(meetingID string, newItem models.ActionItem) error {
	// 获取该会议ID的锁并加写锁
	lock := GetFileLockManager().GetLock(meetingID)
	lock.Lock()
	defer lock.Unlock()

	// 生成随机TodoID
	todoID, err := generateRandomID(8)
	if err != nil {
		return fmt.Errorf("生成TodoID失败: %v", err)
	}
	newItem.TodoID = todoID
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

	// 添加���的ActionItem
	meeting.Summary.ActionItems = append(meeting.Summary.ActionItems, newItem)

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
