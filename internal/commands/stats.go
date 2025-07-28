package commands

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/cache"
	"EverythingSuckz/fsb/internal/types"
	"EverythingSuckz/fsb/internal/utils"
	"fmt"
	"time"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
)

func (m *command) LoadStats(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("stats")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(handlers.NewCommand("stats", stats))
}

func stats(ctx *ext.Context, u *ext.Update) error {
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}
	
	// Check if user is allowed (if restrictions are enabled)
	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, "You are not allowed to use this bot.", nil)
		return dispatcher.EndGroups
	}

	// Get statistics
	statsCache := cache.GetStatsCache()
	if statsCache == nil {
		ctx.Reply(u, "❌ Statistics service is not available at the moment.", nil)
		return dispatcher.EndGroups
	}

	stats, err := statsCache.GetCompleteStats()
	if err != nil {
		// Log error but don't expose it to user
		ctx.Reply(u, "❌ Failed to retrieve statistics. Please try again later.", nil)
		return dispatcher.EndGroups
	}

	// Format the statistics message
	message := formatStatisticsMessage(stats)
	
	ctx.Reply(u, message, nil)
	return dispatcher.EndGroups
}

func formatStatisticsMessage(stats types.StatisticsResponse) string {
	message := "📊 Bot Statistics\n\n"
	
	// Today's stats
	message += fmt.Sprintf("Today: %d files - %s\n", 
		stats.Today.FileCount, 
		utils.FormatFileSizeShort(stats.Today.TotalSize))
	
	// Yesterday's stats
	message += fmt.Sprintf("Yesterday: %d files - %s\n", 
		stats.Yesterday.FileCount, 
		utils.FormatFileSizeShort(stats.Yesterday.TotalSize))
	
	// Last 7 days stats
	message += fmt.Sprintf("Last 7 days: %d files - %s\n", 
		stats.LastWeek.FileCount, 
		utils.FormatFileSizeShort(stats.LastWeek.TotalSize))
	
	// Total stats
	message += fmt.Sprintf("All time: %d files - %s\n\n", 
		stats.Total.FileCount, 
		utils.FormatFileSizeShort(stats.Total.TotalSize))
	
	message += "🔄 Stats are updated in real-time\n"
	message += "⏰ Last updated: " + time.Now().Format("2006-01-02 15:04:05") + "."
	
	return message
} 
