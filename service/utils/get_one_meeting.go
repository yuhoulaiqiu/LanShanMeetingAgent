package utils

import (
	"errors"
	"meetingagent/models"
	"sync"
)

// MeetingStore 封装会议存储和操作
type MeetingStore struct {
	meetings map[string]models.Meeting
	mu       sync.RWMutex
}

var storeInstance *MeetingStore
var once sync.Once

// GetMeetingStore 获取 MeetingStore 单例
func GetMeetingStore() *MeetingStore {
	once.Do(func() {
		storeInstance = &MeetingStore{
			meetings: make(map[string]models.Meeting),
		}
	})
	return storeInstance
}

// GetMeeting retrieves a single meeting by ID
func (ms *MeetingStore) GetMeeting(meetingID string) (models.Meeting, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	meeting, exists := ms.meetings[meetingID]
	if !exists {
		return models.Meeting{}, errors.New("meeting not found")
	}

	return meeting, nil
}

// SetMeeting stores a meeting by ID
func (ms *MeetingStore) SetMeeting(meetingID string, meeting models.Meeting) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.meetings[meetingID] = meeting
}

// GetAllMeetings retrieves all meetings
func (ms *MeetingStore) GetAllMeetings() []models.Meeting {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	meetingsList := make([]models.Meeting, 0, len(ms.meetings))
	for _, meeting := range ms.meetings {
		meetingsList = append(meetingsList, meeting)
	}
	return meetingsList
}
