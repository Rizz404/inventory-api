# Cron Job Configuration for Asset Management

## 📅 Scheduled Tasks

### Warranty Expiring Check
**Schedule**: Every day at 9:00 AM
**Cron Expression**: `0 0 9 * * *`
**Function**: `checkWarrantyExpiring()`

**What it does**:
- Queries assets dengan warranty expiring dalam 30 hari
- Hanya untuk assets yang assigned ke user
- Sends notification ke setiap assigned user
- Logs jumlah assets yang ditemukan

**Query Logic**:
```sql
SELECT * FROM assets
WHERE warranty_end IS NOT NULL
  AND warranty_end > NOW()
  AND warranty_end <= NOW() + INTERVAL '30 days'
  AND assigned_to IS NOT NULL
```

---

### Warranty Expired Check
**Schedule**: Every day at 9:30 AM
**Cron Expression**: `0 30 9 * * *`
**Function**: `checkExpiredWarranties()`

**What it does**:
- Queries assets dengan warranty yang expired dalam 24 jam terakhir
- Hanya untuk assets yang assigned ke user
- Sends notification ke setiap assigned user
- Logs jumlah assets yang ditemukan

**Query Logic**:
```sql
SELECT * FROM assets
WHERE warranty_end IS NOT NULL
  AND warranty_end < NOW()
  AND warranty_end >= NOW() - INTERVAL '1 day'
  AND assigned_to IS NOT NULL
```

## 🎯 Cron Expression Format

```
┌─────────── second (0 - 59)
│ ┌───────── minute (0 - 59)
│ │ ┌─────── hour (0 - 23)
│ │ │ ┌───── day of month (1 - 31)
│ │ │ │ ┌─── month (1 - 12)
│ │ │ │ │ ┌─ day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │ │
* * * * * *
```

### Examples:
```
0 0 9 * * *     # Daily at 9:00 AM
0 30 9 * * *    # Daily at 9:30 AM
0 0 0 * * *     # Daily at midnight
0 0 12 * * 1    # Every Monday at noon
0 0 9 1 * *     # First day of every month at 9:00 AM
0 */15 * * * *  # Every 15 minutes
```

## ⚙️ Configuration

### Change Schedule
Edit `services/asset/cron_service.go`:

```go
func (cs *CronService) Start() error {
    // Change the schedule here
    _, err := cs.cron.AddFunc("0 0 9 * * *", cs.checkWarrantyExpiring)
    if err != nil {
        return err
    }

    _, err = cs.cron.AddFunc("0 30 9 * * *", cs.checkExpiredWarranties)
    if err != nil {
        return err
    }

    cs.cron.Start()
    return nil
}
```

### Change Warranty Window
Edit `services/asset/cron_service.go`:

```go
func (cs *CronService) checkWarrantyExpiring() {
    // Change 30 to desired number of days
    assets, err := cs.assetRepo.GetAssetsWithWarrantyExpiring(ctx, 30)
    // ...
}
```

## 🔍 Monitoring

### Check Logs
```bash
# View cron execution logs
grep "cron service" application.log

# View warranty check logs
grep "warranty check" application.log

# View notification logs
grep "warranty.*notification" application.log
```

### Expected Log Output:
```
[INFO] Asset cron service started successfully
[INFO] Running warranty expiring check...
[INFO] Warranty expiring check completed. Found 5 assets with warranties expiring within 30 days
[INFO] Successfully created warranty expiring notification for asset ID: xxx, user ID: xxx
[INFO] Running expired warranty check...
[INFO] Expired warranty check completed. Found 2 assets with expired warranties
[INFO] Successfully created warranty expired notification for asset ID: xxx, user ID: xxx
```

## 🧪 Testing

### Manual Trigger (Development)
Add temporary function in your test file:

```go
func TestCronManually() {
    // Initialize services...
    cronService := asset.NewCronService(assetRepo, notificationService)

    // Manually trigger functions
    cronService.checkWarrantyExpiring()
    cronService.checkExpiredWarranties()
}
```

### Test with Mock Data
1. Create assets dengan warranty_end dates:
   - Some expiring in 15 days (should trigger)
   - Some expiring in 45 days (should NOT trigger)
   - Some expired yesterday (should trigger)
   - Some expired 2 days ago (should NOT trigger)

2. Assign assets to test users

3. Run cron functions manually or wait for scheduled time

4. Check:
   - Logs untuk execution details
   - Database untuk created notifications
   - FCM untuk delivered notifications

## 🚨 Troubleshooting

### Cron Not Running
**Symptoms**: No log entries untuk cron execution

**Solutions**:
1. Check cron service started di `app/main.go`:
   ```go
   assetCronService := asset.NewCronService(...)
   if err := assetCronService.Start(); err != nil {
       log.Fatalf("Failed to start: %v", err)
   }
   ```

2. Verify cron expression syntax
3. Check server timezone settings
4. Look for startup errors in logs

### No Notifications Sent
**Symptoms**: Cron runs but no notifications created

**Solutions**:
1. Check if assets exist dengan matching criteria
2. Verify assets have `assigned_to` users
3. Check NotificationService is not nil
4. Review notification creation logs
5. Verify FCM configuration

### Duplicate Notifications
**Symptoms**: Users receive same notification multiple times

**Solutions**:
1. Verify cron schedules don't overlap
2. Check only one cron instance running
3. Review warranty date filters
4. Check notification deduplication logic

### Performance Issues
**Symptoms**: Cron execution takes too long

**Solutions**:
1. Add database indexes:
   ```sql
   CREATE INDEX idx_assets_warranty_assigned
   ON assets(warranty_end, assigned_to);
   ```

2. Limit query results if needed
3. Optimize preloads
4. Consider pagination untuk large datasets

## 📊 Performance Metrics

### Expected Performance:
- Query time: < 100ms for 10,000 assets
- Notification creation: < 50ms per notification
- Total execution: < 5 seconds for 100 notifications

### Database Load:
- 2 queries per day (one at 9:00 AM, one at 9:30 AM)
- Filtered queries dengan proper indexes
- Minimal impact on database performance

## 🔐 Security Considerations

1. **User Privacy**: Only send notifications to assigned users
2. **Data Access**: Repository methods respect access controls
3. **Error Handling**: Don't expose sensitive data in logs
4. **Rate Limiting**: Consider if needed for notification service

## 🌍 Timezone Handling

### Important Notes:
- Cron schedules use **server timezone**
- Warranty dates stored in UTC
- Notification times displayed in user's locale

### Recommendations:
1. Document server timezone in deployment docs
2. Consider timezone when setting schedules
3. Test with users in different timezones

### Example Configuration:
```go
// For UTC server
"0 0 9 * * *"   // 9:00 AM UTC

// For Jakarta (UTC+7)
"0 0 2 * * *"   // 9:00 AM Jakarta = 2:00 AM UTC
```

## 📈 Scaling Considerations

### For Large Deployments:
1. **Pagination**: Process assets in batches
2. **Parallel Processing**: Use goroutines dengan rate limiting
3. **Distributed Cron**: Consider distributed lock mechanism
4. **Queue System**: Use message queue untuk notifications

### Example Batch Processing:
```go
func (cs *CronService) checkWarrantyExpiringBatched() {
    batchSize := 1000
    offset := 0

    for {
        assets := cs.getAssetBatch(offset, batchSize)
        if len(assets) == 0 {
            break
        }

        cs.processAssetBatch(assets)
        offset += batchSize
    }
}
```

## 🔧 Advanced Configuration

### Multiple Schedules
Add more cron jobs as needed:

```go
// Check high-value assets daily
_, err = cs.cron.AddFunc("0 0 10 * * *", cs.checkHighValueAssets)

// Weekly summary every Monday
_, err = cs.cron.AddFunc("0 0 9 * * 1", cs.sendWeeklySummary)

// Monthly report first day of month
_, err = cs.cron.AddFunc("0 0 9 1 * *", cs.sendMonthlyReport)
```

### Custom Notification Logic
Override notification methods untuk custom behavior:

```go
func (cs *CronService) sendWarrantyExpiringNotification(ctx context.Context, asset *domain.Asset) {
    // Add custom logic here
    if asset.PurchasePrice != nil && *asset.PurchasePrice > 10000 {
        // Send to admin also
        cs.sendToAdmin(ctx, asset)
    }

    // Original logic
    cs.NotificationService.CreateNotification(...)
}
```

## 📚 References

- [robfig/cron Documentation](https://pkg.go.dev/github.com/robfig/cron/v3)
- [Cron Expression Guide](https://crontab.guru/)
- Asset Notification Guide: `documentation/asset_notification_complete_guide.md`
- Implementation Summary: `documentation/ASSET_NOTIFICATION_IMPLEMENTATION_SUMMARY.md`

## ✅ Checklist for Deployment

- [ ] Configure cron schedules untuk production timezone
- [ ] Test cron execution dengan production-like data
- [ ] Set up monitoring dan alerting
- [ ] Document maintenance procedures
- [ ] Plan for graceful shutdown handling
- [ ] Configure database indexes
- [ ] Test notification delivery
- [ ] Set up log rotation
- [ ] Document troubleshooting procedures
- [ ] Test error recovery scenarios

---

**Last Updated**: 2025-01-16
**Version**: 1.0.0
**Maintainer**: Development Team
