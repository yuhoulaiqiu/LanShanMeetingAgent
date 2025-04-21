package models

type SummarizedMeeting struct {
	MeetingID string `json:"meeting_id"`
	Summary   struct {
		Text        string   `json:"text"`
		KeyPoints   []string `json:"key_points"`
		ActionItems []struct {
			Assignee string `json:"assignee"`
			Task     string `json:"task"`
			Level    string `json:"level"`
			State    string `json:"state"`
			Deadline string `json:"deadline"`
		} `json:"action_items"`
	} `json:"summary"`
}
