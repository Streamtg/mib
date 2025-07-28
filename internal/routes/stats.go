package routes

import (
	"EverythingSuckz/fsb/internal/cache"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (r *allRoutes) LoadStatsAPI(route *Route) {
	route.Engine.GET("/api/stats", r.getStats)
}

func (r *allRoutes) getStats(c *gin.Context) {
	statsCache := cache.GetStatsCache()
	if statsCache == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Statistics service is not available",
		})
		return
	}

	stats, err := statsCache.GetCompleteStats()
	if err != nil {
		r.log.Error("Failed to get statistics", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
} 
