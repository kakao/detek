package detek

import "time"

type ReportExportingFormat struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Labels      []string  `json:"labels"`
	CreatedAt   time.Time `json:"created_at"`

	Level            SeverityLevel            `json:"level"`
	LevelDescription SeverityLevelDescription `json:"level_description"`

	CurrentState string       `json:"current_state"`
	Solution     string       `json:"solution"`
	Problem      JSONableData `json:"problem"`

	Attachments []JSONableData `json:"attachments"`
}
