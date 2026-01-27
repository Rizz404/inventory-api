package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rizz404/inventory-api/domain"
	"github.com/Rizz404/inventory-api/internal/postgresql/gorm/model"
	"github.com/Rizz404/inventory-api/internal/postgresql/mapper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) applyUserFilters(db *gorm.DB, filters *domain.UserFilterOptions) *gorm.DB {
	if filters == nil {
		return db
	}

	if filters.Role != nil {
		db = db.Where("u.role = ?", filters.Role)
	}
	if filters.IsActive != nil {
		db = db.Where("u.is_active = ?", *filters.IsActive)
	}
	if filters.EmployeeID != nil {
		db = db.Where("u.employee_id = ?", *filters.EmployeeID)
	}
	return db
}

func (r *UserRepository) applyUserSorts(db *gorm.DB, sort *domain.UserSortOptions) *gorm.DB {
	if sort == nil || sort.Field == "" {
		return db.Order("u.created_at DESC")
	}

	// Map camelCase sort field to snake_case database column
	columnName := mapper.MapUserSortFieldToColumn(sort.Field)
	orderClause := "u." + columnName

	order := "DESC"
	if sort.Order == domain.SortOrderAsc {
		order = "ASC"
	}
	return db.Order(fmt.Sprintf("%s %s", orderClause, order))
}

// *===========================MUTATION===========================*
func (r *UserRepository) CreateUser(ctx context.Context, payload *domain.User) (domain.User, error) {
	modelUser := mapper.ToModelUserForCreate(payload)

	// Create user in database
	err := r.db.WithContext(ctx).Create(&modelUser).Error
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&modelUser), nil
}

func (r *UserRepository) BulkCreateUsers(ctx context.Context, users []domain.User) ([]domain.User, error) {
	if len(users) == 0 {
		return []domain.User{}, nil
	}

	models := make([]*model.User, len(users))
	for i := range users {
		m := mapper.ToModelUserForCreate(&users[i])
		models[i] = &m
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Omit(clause.Associations).
		Session(&gorm.Session{CreateBatchSize: 500}).
		Create(&models).Error; err != nil {
		tx.Rollback()
		return nil, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	created := make([]domain.User, len(models))
	for i := range models {
		created[i] = mapper.ToDomainUser(models[i])
	}
	return created, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userId string, payload *domain.UpdateUserPayload) (domain.User, error) {
	var updatedUser model.User

	// Build update map from payload
	updates := mapper.ToModelUserUpdateMap(payload)

	// Perform update
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Updates(updates).Error
	if err != nil {
		return domain.User{}, domain.ErrInternal(err)
	}

	// Get updated user
	err = r.db.WithContext(ctx).First(&updatedUser, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&updatedUser), nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, email, passwordHash string) error {
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Update("password_hash", passwordHash).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userId string) error {
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userId).
		Update("last_login", time.Now().UTC()).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId string) error {
	err := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", userId).Error
	if err != nil {
		return domain.ErrInternal(err)
	}
	return nil
}

func (r *UserRepository) BulkDeleteUsers(ctx context.Context, userIds []string) (domain.BulkDeleteUsers, error) {
	result := domain.BulkDeleteUsers{
		RequestedIDS: userIds,
		DeletedIDS:   []string{},
	}

	if len(userIds) == 0 {
		return result, nil
	}

	// First, find which users actually exist
	var existingUsers []model.User
	if err := r.db.WithContext(ctx).Select("id").Where("id IN ?", userIds).Find(&existingUsers).Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	// Collect existing user IDs
	existingIds := make([]string, 0, len(existingUsers))
	for _, user := range existingUsers {
		existingIds = append(existingIds, user.ID.String())
	}

	// If no users exist, return early
	if len(existingIds) == 0 {
		return result, nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return result, domain.ErrInternal(tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete users
	if err := tx.Delete(&model.User{}, "id IN ?", existingIds).Error; err != nil {
		tx.Rollback()
		return result, domain.ErrInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return result, domain.ErrInternal(err)
	}

	result.DeletedIDS = existingIds
	return result, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userId string, hashedPassword string) error {
	// Update only the password_hash and updated_at
	updates := map[string]interface{}{
		"password_hash": hashedPassword,
		"updated_at":    time.Now().UTC(),
	}
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Updates(updates).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound("user")
		}
		return domain.ErrInternal(err)
	}
	return nil
}

// *===========================QUERY===========================*
func (r *UserRepository) GetUsersPaginated(ctx context.Context, params domain.UserParams) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx).
		Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.name ILIKE ? OR u.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyUserFilters(db, params.Filters)

	// Apply sorting
	db = r.applyUserSorts(db, params.Sort)

	// Apply pagination
	if params.Pagination != nil {
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
		if params.Pagination.Offset > 0 {
			db = db.Offset(params.Pagination.Offset)
		}
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	return mapper.ToDomainUsers(users), nil
}

func (r *UserRepository) GetUsersCursor(ctx context.Context, params domain.UserParams) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx).
		Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.name ILIKE ? OR u.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyUserFilters(db, params.Filters)

	// Apply sorting - for cursor pagination, we need consistent ordering by ID
	if params.Sort != nil && params.Sort.Field != "" {
		db = r.applyUserSorts(db, params.Sort)
	} else {
		// Default to ID ASC for cursor pagination
		db = db.Order("u.id ASC")
	}

	// Apply cursor-based pagination
	if params.Pagination != nil {
		if params.Pagination.Cursor != "" {
			db = db.Where("u.id > ?", params.Pagination.Cursor)
		}
		if params.Pagination.Limit > 0 {
			db = db.Limit(params.Pagination.Limit)
		}
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	return mapper.ToDomainUsers(users), nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId string) (domain.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).First(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (domain.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user with name '" + name + "'")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound("user with email '" + email + "'")
		}
		return domain.User{}, domain.ErrInternal(err)
	}

	return mapper.ToDomainUser(&user), nil
}

func (r *UserRepository) CheckUserExists(ctx context.Context, userId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckNameExists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckNameExistsExcluding(ctx context.Context, name string, excludeUserId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("name = ? AND id != ?", name, excludeUserId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CheckEmailExistsExcluding(ctx context.Context, email string, excludeUserId string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ? AND id != ?", email, excludeUserId).Count(&count).Error
	if err != nil {
		return false, domain.ErrInternal(err)
	}
	return count > 0, nil
}

func (r *UserRepository) CountUsers(ctx context.Context, params domain.UserParams) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.name ILIKE ? OR u.full_name ILIKE ?", searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyUserFilters(db, params.Filters)

	if err := db.Count(&count).Error; err != nil {
		return 0, domain.ErrInternal(err)
	}
	return count, nil
}

func (r *UserRepository) GetUserStatistics(ctx context.Context) (domain.UserStatistics, error) {
	var stats domain.UserStatistics

	// Get total user count
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&totalCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.Total.Count = int(totalCount)

	// Get user counts by status (active/inactive)
	var activeCount, inactiveCount int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", true).Count(&activeCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", false).Count(&inactiveCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByStatus.Active = int(activeCount)
	stats.ByStatus.Inactive = int(inactiveCount)

	// Get user counts by role
	var adminCount, staffCount, employeeCount int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", domain.RoleAdmin).Count(&adminCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", domain.RoleStaff).Count(&staffCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", domain.RoleEmployee).Count(&employeeCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.ByRole.Admin = int(adminCount)
	stats.ByRole.Staff = int(staffCount)
	stats.ByRole.Employee = int(employeeCount)

	// Get registration trends (last 30 days)
	var registrationTrends []struct {
		Date  time.Time `json:"date"`
		Count int64     `json:"count"`
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '30 days'").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&registrationTrends).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.RegistrationTrends = make([]domain.RegistrationTrend, len(registrationTrends))
	for i, rt := range registrationTrends {
		stats.RegistrationTrends[i] = domain.RegistrationTrend{
			Date:  rt.Date,
			Count: int(rt.Count),
		}
	}

	// Calculate summary statistics
	stats.Summary.TotalUsers = int(totalCount)

	if totalCount > 0 {
		stats.Summary.ActiveUsersPercentage = float64(activeCount) / float64(totalCount) * 100
		stats.Summary.InactiveUsersPercentage = float64(inactiveCount) / float64(totalCount) * 100
		stats.Summary.AdminPercentage = float64(adminCount) / float64(totalCount) * 100
		stats.Summary.StaffPercentage = float64(staffCount) / float64(totalCount) * 100
		stats.Summary.EmployeePercentage = float64(employeeCount) / float64(totalCount) * 100
	}

	// Get earliest and latest registration dates
	var earliestDate, latestDate time.Time
	if err := r.db.WithContext(ctx).Model(&model.User{}).Select("MIN(created_at)").Scan(&earliestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	if err := r.db.WithContext(ctx).Model(&model.User{}).Select("MAX(created_at)").Scan(&latestDate).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.Summary.EarliestRegistrationDate = earliestDate
	stats.Summary.LatestRegistrationDate = latestDate

	// Calculate average users per day
	if !earliestDate.IsZero() && !latestDate.IsZero() {
		daysDiff := latestDate.Sub(earliestDate).Hours() / 24
		if daysDiff > 0 {
			stats.Summary.AverageUsersPerDay = float64(totalCount) / daysDiff
		}
	}

	return stats, nil
}

func (r *UserRepository) GetUsersForExport(ctx context.Context, params domain.UserParams) ([]domain.User, error) {
	var users []model.User
	db := r.db.WithContext(ctx).
		Table("users u")

	if params.SearchQuery != nil && *params.SearchQuery != "" {
		searchPattern := "%" + *params.SearchQuery + "%"
		db = db.Where("u.full_name ILIKE ? OR u.email ILIKE ? OR u.employee_id ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Apply filters
	db = r.applyUserFilters(db, params.Filters)

	// Apply sorting
	db = r.applyUserSorts(db, params.Sort)

	// No pagination for export - get all matching users
	if err := db.Find(&users).Error; err != nil {
		return nil, domain.ErrInternal(err)
	}

	// Convert to domain users
	return mapper.ToDomainUsers(users), nil
}

func (r *UserRepository) GetUserPersonalStatistics(ctx context.Context, userId string) (domain.UserPersonalStatistics, error) {
	var stats domain.UserPersonalStatistics

	// Get user info
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stats, domain.ErrNotFound("user")
		}
		return stats, domain.ErrInternal(err)
	}

	// * === ASSET STATISTICS ===
	// Get assigned assets (Checkout without return)
	type AssetItem struct {
		AssetID      string    `gorm:"column:asset_id"`
		AssetTag     string    `gorm:"column:asset_tag"`
		Name         string    `gorm:"column:name"`
		Category     string    `gorm:"column:category_name"`
		Condition    string    `gorm:"column:condition"`
		Value        float64   `gorm:"column:purchase_price"`
		AssignedDate time.Time `gorm:"column:assigned_date"`
	}

	var assetItems []AssetItem
	err := r.db.WithContext(ctx).
		Table("asset_movements am").
		Select(`
			a.id as asset_id,
			a.asset_tag,
			a.name,
			c.name as category_name,
			a.condition,
			a.purchase_price,
			am.created_at as assigned_date
		`).
		Joins("JOIN assets a ON am.asset_id = a.id").
		Joins("JOIN categories c ON a.category_id = c.id").
		Where("am.assigned_to = ?", userId).
		Where("am.movement_type = ?", "Checkout").
		Where("am.returned_at IS NULL").
		Where("a.deleted_at IS NULL").
		Where("am.deleted_at IS NULL").
		Find(&assetItems).Error

	if err != nil {
		return stats, domain.ErrInternal(err)
	}

	// Calculate asset statistics
	stats.Assets.Total.Count = len(assetItems)
	var totalValue float64
	conditionCount := make(map[string]int)

	for _, item := range assetItems {
		totalValue += item.Value
		conditionCount[item.Condition]++
	}

	stats.Assets.Total.TotalValue = totalValue
	stats.Assets.ByCondition.Good = conditionCount["Good"]
	stats.Assets.ByCondition.Fair = conditionCount["Fair"]
	stats.Assets.ByCondition.Poor = conditionCount["Poor"]
	stats.Assets.ByCondition.Damaged = conditionCount["Damaged"]

	// Convert to domain items
	stats.Assets.Items = make([]domain.UserPersonalAssetItem, len(assetItems))
	for i, item := range assetItems {
		stats.Assets.Items[i] = domain.UserPersonalAssetItem{
			AssetID:      item.AssetID,
			AssetTag:     item.AssetTag,
			Name:         item.Name,
			Category:     item.Category,
			Condition:    item.Condition,
			Value:        item.Value,
			AssignedDate: item.AssignedDate,
		}
	}

	// * === ISSUE REPORT STATISTICS ===
	// Get total count
	var totalIssueCount int64
	if err := r.db.WithContext(ctx).
		Model(&model.IssueReport{}).
		Where("reported_by = ?", userId).
		Count(&totalIssueCount).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}
	stats.IssueReports.Total.Count = int(totalIssueCount)

	// Get counts by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.IssueReport{}).
		Select("status, COUNT(*) as count").
		Where("reported_by = ?", userId).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "Open":
			stats.IssueReports.ByStatus.Open = int(sc.Count)
		case "InProgress":
			stats.IssueReports.ByStatus.InProgress = int(sc.Count)
		case "Resolved":
			stats.IssueReports.ByStatus.Resolved = int(sc.Count)
		case "Closed":
			stats.IssueReports.ByStatus.Closed = int(sc.Count)
		}
	}

	// Get counts by priority
	var priorityCounts []struct {
		Priority string
		Count    int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.IssueReport{}).
		Select("priority, COUNT(*) as count").
		Where("reported_by = ?", userId).
		Group("priority").
		Scan(&priorityCounts).Error; err != nil {
		return stats, domain.ErrInternal(err)
	}

	for _, pc := range priorityCounts {
		switch pc.Priority {
		case "High":
			stats.IssueReports.ByPriority.High = int(pc.Count)
		case "Medium":
			stats.IssueReports.ByPriority.Medium = int(pc.Count)
		case "Low":
			stats.IssueReports.ByPriority.Low = int(pc.Count)
		}
	}

	// Get recent issues (last 10)
	type IssueItem struct {
		IssueID      string     `gorm:"column:issue_id"`
		AssetID      *string    `gorm:"column:asset_id"`
		AssetTag     *string    `gorm:"column:asset_tag"`
		Title        string     `gorm:"column:title"`
		Priority     string     `gorm:"column:priority"`
		Status       string     `gorm:"column:status"`
		ReportedDate time.Time  `gorm:"column:reported_date"`
	}

	var recentIssues []IssueItem
	err = r.db.WithContext(ctx).
		Table("issue_reports ir").
		Select(`
			ir.id as issue_id,
			ir.asset_id,
			a.asset_tag,
			ir.title,
			ir.priority,
			ir.status,
			ir.created_at as reported_date
		`).
		Joins("LEFT JOIN assets a ON ir.asset_id = a.id").
		Where("ir.reported_by = ?", userId).
		Where("ir.deleted_at IS NULL").
		Order("ir.created_at DESC").
		Limit(10).
		Find(&recentIssues).Error

	if err != nil {
		return stats, domain.ErrInternal(err)
	}

	stats.IssueReports.RecentIssues = make([]domain.UserPersonalIssueReportItem, len(recentIssues))
	for i, issue := range recentIssues {
		stats.IssueReports.RecentIssues[i] = domain.UserPersonalIssueReportItem{
			IssueID:      issue.IssueID,
			AssetID:      issue.AssetID,
			AssetTag:     issue.AssetTag,
			Title:        issue.Title,
			Priority:     issue.Priority,
			Status:       issue.Status,
			ReportedDate: issue.ReportedDate,
		}
	}

	// Calculate summary
	stats.IssueReports.Summary.OpenIssuesCount = stats.IssueReports.ByStatus.Open
	stats.IssueReports.Summary.ResolvedIssuesCount = stats.IssueReports.ByStatus.Resolved

	// Calculate average resolution time for resolved issues
	var avgResolutionDays float64
	if err := r.db.WithContext(ctx).
		Model(&model.IssueReport{}).
		Select("AVG(EXTRACT(EPOCH FROM (resolved_at - created_at))/86400) as avg_days").
		Where("reported_by = ?", userId).
		Where("status IN ?", []string{"Resolved", "Closed"}).
		Where("resolved_at IS NOT NULL").
		Scan(&avgResolutionDays).Error; err != nil {
		// Non-critical, continue
		avgResolutionDays = 0
	}
	stats.IssueReports.Summary.AverageResolutionDays = avgResolutionDays

	// * === SUMMARY STATISTICS ===
	stats.Summary.AccountCreatedDate = user.CreatedAt
	stats.Summary.LastLogin = user.LastLogin

	// Calculate account age
	accountAge := time.Since(user.CreatedAt)
	days := int(accountAge.Hours() / 24)
	stats.Summary.AccountAge = fmt.Sprintf("%d days", days)

	// Check if has active issues
	stats.Summary.HasActiveIssues = stats.IssueReports.ByStatus.Open > 0 || stats.IssueReports.ByStatus.InProgress > 0

	// Calculate health score (100 base)
	healthScore := 100
	// -5 points per Fair condition asset
	healthScore -= stats.Assets.ByCondition.Fair * 5
	// -10 points per Poor condition asset
	healthScore -= stats.Assets.ByCondition.Poor * 10
	// -15 points per Damaged condition asset
	healthScore -= stats.Assets.ByCondition.Damaged * 15
	// -10 points per open High priority issue
	healthScore -= stats.IssueReports.ByPriority.High * 10
	// -5 points per open Medium priority issue (only count open/in-progress)
	openMediumIssues := 0
	if stats.IssueReports.ByStatus.Open > 0 || stats.IssueReports.ByStatus.InProgress > 0 {
		// Rough estimate, actual calculation would need more detailed query
		openMediumIssues = stats.IssueReports.ByPriority.Medium / 2
	}
	healthScore -= openMediumIssues * 5

	// Ensure health score is between 0-100
	if healthScore < 0 {
		healthScore = 0
	}
	stats.Summary.HealthScore = healthScore

	return stats, nil
}
