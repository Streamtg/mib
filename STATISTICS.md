# Statistics Feature

This document describes the statistics feature that has been added to the Telegram File Stream Bot.

## Overview

The statistics feature tracks file processing metrics in real-time, providing insights into bot usage patterns.

## Features

### 📊 Statistics Tracking
- **Daily Statistics**: Files processed and total size for today
- **Yesterday Statistics**: Files processed and total size for yesterday
- **Weekly Statistics**: Files processed and total size for the last 7 days
- **All-time Statistics**: Total files processed and size since bot creation

### 🎯 Commands
- `/stats` - Display current statistics in the chat

### 🌐 API Endpoints
- `GET /api/stats` - JSON API endpoint for statistics

## Database

The statistics are stored in a SQLite database located at `data/fsb_stats.db`. The database automatically creates the necessary tables on first run.

### Database Schema

```sql
CREATE TABLE file_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL,
    file_count BIGINT NOT NULL DEFAULT 0,
    total_size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Usage

### Telegram Bot Commands

1. **View Statistics**: Send `/stats` to the bot
   ```
   📊 Bot Statistics

   Today: 145 files - 96.93 GB
   Yesterday: 1652 files - 1.62 TB
   Last 7 days: 10724 files - 11.03 TB
   All time: 50000 files - 50.5 TB

   🔄 Stats are updated in real-time
   ⏰ Last updated: 2024-01-15 14:30:25
   ```

### API Access

Get statistics via HTTP API:
```bash
curl https://your-bot-domain.com/api/stats
```

Response:
```json
{
  "success": true,
  "data": {
    "today": {
      "date": "2024-01-15T00:00:00Z",
      "file_count": 145,
      "total_size": 104073748480
    },
    "yesterday": {
      "date": "2024-01-14T00:00:00Z",
      "file_count": 1652,
      "total_size": 1781204459520
    },
    "last_week": {
      "start_date": "2024-01-08T00:00:00Z",
      "end_date": "2024-01-15T00:00:00Z",
      "file_count": 10724,
      "total_size": 11848234598400
    },
    "total": {
      "date": "2024-01-15T14:30:25Z",
      "file_count": 50000,
      "total_size": 54250000000000
    }
  }
}
```

## Implementation Details

### Automatic Tracking
Statistics are automatically recorded whenever a file is processed by the bot. No manual intervention is required.

### Real-time Updates
Statistics are updated in real-time as files are processed. The `/stats` command shows the most current data.

### File Size Formatting
File sizes are automatically formatted into human-readable units (B, KB, MB, GB, TB) for easy reading.

### Error Handling
- If the database is unavailable, the bot will continue to function normally
- Statistics recording errors are logged but don't affect file processing
- Graceful degradation ensures the bot remains operational

## Configuration

The statistics feature requires no additional configuration. It automatically:
- Creates the database directory (`data/`)
- Initializes the SQLite database
- Creates necessary tables
- Starts tracking statistics

## Deployment Notes

### Koyeb Deployment
When deploying to Koyeb:
1. The database file will be created in the `data/` directory
2. Statistics will persist across deployments
3. The API endpoint will be available at `https://your-app.koyeb.app/api/stats`

### File Permissions
Ensure the application has write permissions to create the `data/` directory and database file.

## Monitoring

### Logs
Statistics-related operations are logged with the following prefixes:
- `database` - Database initialization and operations
- `stats_cache` - Statistics cache operations
- `stats` - Statistics command handling

### Database Maintenance
The SQLite database is lightweight and requires minimal maintenance. Consider:
- Regular backups of the `data/fsb_stats.db` file
- Monitoring database size for very long-running instances
- Archiving old statistics data if needed

## Future Enhancements

Potential future improvements:
- Export statistics to external analytics platforms
- Custom date range queries
- User-specific statistics
- Performance metrics and trends
- Automated reporting 
