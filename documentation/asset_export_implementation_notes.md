# Asset Export Implementation Notes

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                       Handler Layer                          │
│  - ExportAssetList                                          │
│  - ExportAssetStatistics                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│  - ExportAssetList                                          │
│  - exportAssetListToPDF                                     │
│  - exportAssetListToExcel                                   │
│  - ExportAssetStatistics                                    │
│  - exportAssetStatisticsToPDF                               │
│  - generateStatusChart                                      │
│  - generateConditionChart                                   │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                           │
│  - GetAssetsForExport                                       │
│  - GetAssetStatistics                                       │
└─────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. No Pagination for Export
**Decision**: `GetAssetsForExport` fetches all matching records without pagination

**Rationale**:
- Export should contain all filtered data
- User expects complete dataset in export file

**Considerations**:
- Memory usage for large datasets
- May need to add limit/streaming for very large exports (10k+ records)

**Alternative**: Implement streaming export for large datasets

### 2. Chart Generation Strategy
**Decision**: Generate charts as temporary HTML files, then embed in PDF

**Rationale**:
- go-echarts generates HTML output natively
- Maroto supports image embedding in PDF
- Clean separation of concerns

**Considerations**:
- Temporary files need cleanup (handled with defer)
- Limited to static charts (no interactivity)

**Alternative**: Use a chart library that directly generates images (e.g., go-chart)

### 3. PDF Layout
**Decision**:
- Asset List: Landscape orientation
- Statistics: Portrait orientation

**Rationale**:
- Asset list has many columns → landscape fits better
- Statistics with charts → portrait is more standard

### 4. Error Handling
**Decision**: Return errors immediately, no partial exports

**Rationale**:
- Ensures data integrity
- Prevents incomplete exports

**Alternative**: Could implement partial export with warnings

## Performance Optimizations

### Current Implementation
1. Single query for all export data
2. In-memory processing
3. Synchronous generation

### Potential Improvements

#### 1. Streaming for Large Datasets
```go
// Instead of loading all in memory
func (s *Service) ExportAssetListStreaming(ctx context.Context, payload *domain.ExportAssetListPayload, w io.Writer) error {
    // Stream data directly to writer
    // Process in batches
}
```

#### 2. Background Job Processing
```go
// For very large exports
type ExportJob struct {
    ID        string
    Status    string
    FileURL   string
    CreatedAt time.Time
}

func (s *Service) CreateExportJob(ctx context.Context, payload *domain.ExportAssetListPayload) (*ExportJob, error) {
    // Create job and process in background
    // Notify user when complete
}
```

#### 3. Caching Statistics
```go
// Cache statistics for improved performance
func (s *Service) GetAssetStatistics(ctx context.Context) (domain.AssetStatistics, error) {
    // Check cache first
    // Return cached if recent (e.g., < 5 minutes old)
    // Otherwise, fetch and cache
}
```

#### 4. Parallel Chart Generation
```go
// Generate multiple charts concurrently
func (s *Service) generateAllCharts(stats domain.AssetStatistics) (map[string]string, error) {
    var wg sync.WaitGroup
    charts := make(map[string]string)

    wg.Add(2)
    go func() {
        defer wg.Done()
        charts["status"], _ = s.generateStatusChart(stats.ByStatus)
    }()

    go func() {
        defer wg.Done()
        charts["condition"], _ = s.generateConditionChart(stats.ByCondition)
    }()

    wg.Wait()
    return charts, nil
}
```

## Security Considerations

### 1. Authentication & Authorization
**Current**: Uses middleware for auth checks
**Improvement**: Add role-based access control for export permissions

```go
// In handler
middleware.AuthorizeRole(domain.RoleAdmin, domain.RoleUser)
```

### 2. Data Sanitization
**Current**: Direct database output
**Improvement**: Sanitize sensitive data before export

```go
func (s *Service) sanitizeAssetForExport(asset *domain.Asset) *domain.Asset {
    // Remove or mask sensitive fields
    if shouldMaskPriceData(user) {
        asset.PurchasePrice = nil
    }
    return asset
}
```

### 3. Rate Limiting
**Current**: None
**Improvement**: Add rate limiting for export endpoints

```go
// Middleware to limit exports per user
middleware.RateLimit(5, time.Minute) // 5 exports per minute
```

### 4. File Size Limits
**Current**: No limits
**Improvement**: Add max file size checks

```go
const MaxExportRows = 50000

func (s *Service) ExportAssetList(...) ([]byte, string, error) {
    count, _ := s.Repo.CountAssets(ctx, params)
    if count > MaxExportRows {
        return nil, "", domain.ErrBadRequest("Too many records to export")
    }
    // ...
}
```

## Testing Strategy

### Unit Tests

#### 1. Service Layer Tests
```go
func TestExportAssetListToPDF(t *testing.T) {
    // Test PDF generation with sample data
    // Verify PDF structure and content
}

func TestExportAssetListToExcel(t *testing.T) {
    // Test Excel generation with sample data
    // Verify Excel structure and data
}

func TestGenerateStatusChart(t *testing.T) {
    // Test chart generation
    // Verify chart file is created
}
```

#### 2. Handler Tests
```go
func TestExportAssetListHandler(t *testing.T) {
    // Test with valid payload
    // Test with invalid format
    // Test with filters
    // Test auth requirements
}
```

### Integration Tests
```go
func TestExportEndToEnd(t *testing.T) {
    // Create test assets in DB
    // Call export endpoint
    // Verify file content
    // Cleanup test data
}
```

## Monitoring & Logging

### Metrics to Track
1. Export request count by format
2. Export generation time
3. Export file sizes
4. Failed exports
5. Memory usage during exports

### Logging Implementation
```go
func (s *Service) ExportAssetList(ctx context.Context, payload *domain.ExportAssetListPayload, langCode string) ([]byte, string, error) {
    start := time.Now()
    log.Info().
        Str("format", string(payload.Format)).
        Str("language", langCode).
        Msg("Starting asset export")

    // ... export logic ...

    log.Info().
        Str("format", string(payload.Format)).
        Int("asset_count", len(assets)).
        Dur("duration", time.Since(start)).
        Msg("Asset export completed")

    return data, filename, nil
}
```

## Future Enhancements

### 1. Additional Export Formats
- CSV
- JSON
- XML

### 2. Custom Templates
Allow users to define custom export templates:
```go
type ExportTemplate struct {
    ID      string
    Name    string
    Format  ExportFormat
    Columns []string
    Filters map[string]interface{}
}

func (s *Service) ExportWithTemplate(ctx context.Context, templateID string) ([]byte, string, error) {
    // Load template
    // Apply template configuration
    // Generate export
}
```

### 3. Scheduled Exports
```go
type ScheduledExport struct {
    ID       string
    Name     string
    Schedule string // cron expression
    Template *ExportTemplate
    Email    string
}

// Background job to run scheduled exports
func (s *Service) ProcessScheduledExports(ctx context.Context) error {
    // Find due exports
    // Generate files
    // Send via email or save to storage
}
```

### 4. Export History
```go
type ExportHistory struct {
    ID          string
    UserID      string
    Format      ExportFormat
    RowCount    int
    FileSize    int64
    GeneratedAt time.Time
    DownloadURL string
}

// Track all exports for auditing
func (s *Service) LogExport(ctx context.Context, export *ExportHistory) error {
    // Save to database
}
```

### 5. Cloud Storage Integration
```go
// Save exports to S3/GCS instead of direct download
func (s *Service) ExportToCloud(ctx context.Context, payload *domain.ExportAssetListPayload) (*ExportResult, error) {
    data, filename, err := s.ExportAssetList(ctx, payload, "en")
    if err != nil {
        return nil, err
    }

    // Upload to cloud storage
    url, err := s.CloudStorage.Upload(ctx, filename, data)
    if err != nil {
        return nil, err
    }

    return &ExportResult{
        URL:       url,
        Filename:  filename,
        Size:      len(data),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }, nil
}
```

### 6. Advanced Charts
- Bar charts for category distribution
- Line charts for asset trends over time
- Heat maps for location distribution
- Custom color schemes

### 7. Email Delivery
```go
func (s *Service) EmailExport(ctx context.Context, payload *domain.ExportAssetListPayload, email string) error {
    // Generate export
    data, filename, err := s.ExportAssetList(ctx, payload, "en")
    if err != nil {
        return err
    }

    // Send via email service
    return s.EmailService.SendWithAttachment(ctx, email, "Asset Export", filename, data)
}
```

## Dependencies

### Current
- `github.com/xuri/excelize/v2` - Excel generation
- `github.com/johnfercher/maroto/v2` - PDF generation
- `github.com/go-echarts/go-echarts/v2` - Chart generation

### Potential Additions
- `github.com/jung-kurt/gofpdf` - Alternative PDF library
- `github.com/wcharczuk/go-chart` - Native image chart generation
- `github.com/aws/aws-sdk-go-v2` - S3 integration
- `gopkg.in/gomail.v2` - Email delivery

## Migration Guide

If upgrading from a previous version without export functionality:

1. Install new dependencies:
```bash
go get github.com/xuri/excelize/v2
go get github.com/johnfercher/maroto/v2@v2.3.1
go get github.com/go-echarts/go-echarts/v2
```

2. Run database migrations (if any new tables needed for export history)

3. Update API documentation

4. Test all export endpoints thoroughly

5. Monitor performance metrics after deployment

## Troubleshooting

### Common Issues

#### 1. Memory Usage Too High
**Symptom**: Server runs out of memory during large exports
**Solution**: Implement streaming or pagination

#### 2. Slow Export Generation
**Symptom**: Exports take too long to generate
**Solution**:
- Add caching for statistics
- Optimize database queries
- Use concurrent processing for charts

#### 3. Chart Images Not Appearing in PDF
**Symptom**: PDF generated but charts are missing
**Solution**:
- Check temp file permissions
- Verify chart generation doesn't error
- Check Maroto image component usage

#### 4. Incorrect Data in Excel
**Symptom**: Excel file contains wrong or corrupted data
**Solution**:
- Verify data mapping from domain to Excel
- Check for null pointer dereferences
- Test with various data scenarios
