package cache

import (
	"EverythingSuckz/fsb/internal/database"
	"EverythingSuckz/fsb/internal/types"
	"fmt"
	"time"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

type StatsCache struct {
	db  *gorm.DB
	log *zap.Logger
}

var statsCache *StatsCache

func InitStatsCache(log *zap.Logger) {
	log = log.Named("stats_cache")
	defer log.Sugar().Info("Initialized stats cache")
	
	db := database.GetDB()
	if db == nil {
		log.Error("Database not initialized")
		return
	}
	
	statsCache = &StatsCache{
		db:  db,
		log: log,
	}
}

func GetStatsCache() *StatsCache {
	return statsCache
}

// RecordFileProcessed records a file processing event
func (sc *StatsCache) RecordFileProcessed(fileSize int64) error {
	today := time.Now().Truncate(24 * time.Hour)
	
	var stats types.Stats
	result := sc.db.Where("date = ?", today).First(&stats)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new record for today
			stats = types.Stats{
				Date:      today,
				FileCount: 1,
				TotalSize: fileSize,
			}
			return sc.db.Create(&stats).Error
		}
		return result.Error
	}
	
	// Update existing record
	stats.FileCount++
	stats.TotalSize += fileSize
	return sc.db.Save(&stats).Error
}

// GetTodayStats returns today's statistics
func (sc *StatsCache) GetTodayStats() (types.DailyStats, error) {
	today := time.Now().Truncate(24 * time.Hour)
	
	var stats types.Stats
	result := sc.db.Where("date = ?", today).First(&stats)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return types.DailyStats{
				Date:      today,
				FileCount: 0,
				TotalSize: 0,
			}, nil
		}
		return types.DailyStats{}, result.Error
	}
	
	return types.DailyStats{
		Date:      stats.Date,
		FileCount: stats.FileCount,
		TotalSize: stats.TotalSize,
	}, nil
}

// GetYesterdayStats returns yesterday's statistics
func (sc *StatsCache) GetYesterdayStats() (types.DailyStats, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	
	var stats types.Stats
	result := sc.db.Where("date = ?", yesterday).First(&stats)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return types.DailyStats{
				Date:      yesterday,
				FileCount: 0,
				TotalSize: 0,
			}, nil
		}
		return types.DailyStats{}, result.Error
	}
	
	return types.DailyStats{
		Date:      stats.Date,
		FileCount: stats.FileCount,
		TotalSize: stats.TotalSize,
	}, nil
}

// GetLastWeekStats returns the last 7 days statistics
func (sc *StatsCache) GetLastWeekStats() (types.WeeklyStats, error) {
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -7)
	
	var result struct {
		FileCount int64 `gorm:"column:file_count"`
		TotalSize int64 `gorm:"column:total_size"`
	}
	
	err := sc.db.Model(&types.Stats{}).
		Select("COALESCE(SUM(file_count), 0) as file_count, COALESCE(SUM(total_size), 0) as total_size").
		Where("date >= ? AND date < ?", startDate, endDate).
		Scan(&result).Error
	
	if err != nil {
		return types.WeeklyStats{}, err
	}
	
	return types.WeeklyStats{
		StartDate: startDate,
		EndDate:   endDate,
		FileCount: result.FileCount,
		TotalSize: result.TotalSize,
	}, nil
}

// GetTotalStats returns all-time statistics
func (sc *StatsCache) GetTotalStats() (types.DailyStats, error) {
	var result struct {
		FileCount int64 `gorm:"column:file_count"`
		TotalSize int64 `gorm:"column:total_size"`
	}
	
	err := sc.db.Model(&types.Stats{}).
		Select("COALESCE(SUM(file_count), 0) as file_count, COALESCE(SUM(total_size), 0) as total_size").
		Scan(&result).Error
	
	if err != nil {
		return types.DailyStats{}, err
	}
	
	return types.DailyStats{
		Date:      time.Now(),
		FileCount: result.FileCount,
		TotalSize: result.TotalSize,
	}, nil
}

// GetCompleteStats returns all statistics in one call
func (sc *StatsCache) GetCompleteStats() (types.StatisticsResponse, error) {
	today, err := sc.GetTodayStats()
	if err != nil {
		return types.StatisticsResponse{}, fmt.Errorf("failed to get today stats: %w", err)
	}
	
	yesterday, err := sc.GetYesterdayStats()
	if err != nil {
		return types.StatisticsResponse{}, fmt.Errorf("failed to get yesterday stats: %w", err)
	}
	
	lastWeek, err := sc.GetLastWeekStats()
	if err != nil {
		return types.StatisticsResponse{}, fmt.Errorf("failed to get last week stats: %w", err)
	}
	
	total, err := sc.GetTotalStats()
	if err != nil {
		return types.StatisticsResponse{}, fmt.Errorf("failed to get total stats: %w", err)
	}
	
	return types.StatisticsResponse{
		Today:     today,
		Yesterday: yesterday,
		LastWeek:  lastWeek,
		Total:     total,
	}, nil
} 
