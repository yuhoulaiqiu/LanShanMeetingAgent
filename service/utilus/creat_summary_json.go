package utilus

import (
	"encoding/json"
	"fmt"
	"meetingagent/models"
	"os"
	"path/filepath"
)

func CreateMeetingJSON(text string, meetingID string) error {
	// 获取该会议ID的锁并加写锁
	lock := GetFileLockManager().GetLock(meetingID)
	lock.Lock()
	defer lock.Unlock()

	// 创建会议结构
	meeting := models.SummarizedMeeting{
		MeetingID: meetingID,
	}

	// 设置摘要信息
	meeting.Summary.Text = text
	meeting.Summary.KeyPoints = []string{}

	// 转换为JSON
	jsonData, err := json.MarshalIndent(meeting, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化会议数据失败: %v", err)
	}

	// 确保目录存在
	meetingsDir := filepath.Join("data", "meetings")
	if err := os.MkdirAll(meetingsDir, 0755); err != nil {
		return fmt.Errorf("创建会议数据目录失败: %v", err)
	}

	// 写入文件
	filePath := filepath.Join(meetingsDir, meetingID+".json")
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入会议文件失败: %v", err)
	}

	return nil
}
