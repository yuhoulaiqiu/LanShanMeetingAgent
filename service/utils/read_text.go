package utils

import (
	"encoding/json"
	"fmt"
	"meetingagent/models"
	"os"
	"path/filepath"
)

// ReadMeetingSummaryText 读取指定会议摘要的文本内容
func ReadMeetingSummaryText(meetingID string) (string, error) {
	// 获取该会议ID的锁并加读锁
	lock := GetFileLockManager().GetLock(meetingID)
	lock.RLock()
	defer lock.RUnlock()

	// 构建文件路径
	filePath := filepath.Join("data", "meetings", meetingID+".json")

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取会议摘要文件失败: %v", err)
	}

	// 解析JSON
	var meeting models.SummarizedMeeting
	if err := json.Unmarshal(data, &meeting); err != nil {
		return "", fmt.Errorf("解析会议摘要数据失败: %v", err)
	}

	return meeting.Summary.Text, nil
}
