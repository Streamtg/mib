package utils

import (
	"fmt"
	"math"
)

// FormatFileSize formats bytes into human readable format
func FormatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	
	const unit = 1024
	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	pre := "KMGTPE"
	if exp == 0 {
		return fmt.Sprintf("%d B", bytes)
	}
	
	exp--
	if exp >= len(pre) {
		exp = len(pre) - 1
	}
	
	val := float64(bytes) / math.Pow(unit, float64(exp+1))
	return fmt.Sprintf("%.2f %cB", val, pre[exp])
}

// FormatFileSizeShort formats bytes into short human readable format (e.g., 1.5 GB)
func FormatFileSizeShort(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	
	const unit = 1024
	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	pre := "KMGTPE"
	if exp == 0 {
		return fmt.Sprintf("%d B", bytes)
	}
	
	exp--
	if exp >= len(pre) {
		exp = len(pre) - 1
	}
	
	val := float64(bytes) / math.Pow(unit, float64(exp+1))
	return fmt.Sprintf("%.1f %cB", val, pre[exp])
} 
