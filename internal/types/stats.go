package types

import (
	"time"
)

// Stats represents the statistics for file processing
type Stats struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Date      time.Time `gorm:"index;not null"`
	FileCount int64     `gorm:"not null;default:0"`
	TotalSize int64     `gorm:"not null;default:0"` // in bytes
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// DailyStats represents today's statistics
type DailyStats struct {
	Date      time.Time `json:"date"`
	FileCount int64     `json:"file_count"`
	TotalSize int64     `json:"total_size"` // in bytes
}

// WeeklyStats represents the last 7 days statistics
type WeeklyStats struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	FileCount int64     `json:"file_count"`
	TotalSize int64     `json:"total_size"` // in bytes
}

// StatisticsResponse represents the complete statistics response
type StatisticsResponse struct {
	Today     DailyStats  `json:"today"`
	Yesterday DailyStats  `json:"yesterday"`
	LastWeek  WeeklyStats `json:"last_week"`
	Total     DailyStats  `json:"total"`
}

// TableName specifies the table name for Stats
func (Stats) TableName() string {
	return "file_stats"
} 
