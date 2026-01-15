package migration

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"gorm.io/gorm"
)

// Repository handles migration job database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new migration repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new migration job
func (r *Repository) Create(job *models.MigrationJob) error {
	return r.db.Create(job).Error
}

// FindByID finds a migration job by ID
func (r *Repository) FindByID(id uint) (*models.MigrationJob, error) {
	var job models.MigrationJob
	err := r.db.First(&job, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// FindByIDAndUser finds a migration job by ID and user
func (r *Repository) FindByIDAndUser(id, userID uint) (*models.MigrationJob, error) {
	var job models.MigrationJob
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// List lists migration jobs for a user
func (r *Repository) List(userID uint, status string, page, perPage int) ([]models.MigrationJob, int64, error) {
	var jobs []models.MigrationJob
	var total int64

	query := r.db.Model(&models.MigrationJob{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// Update updates a migration job
func (r *Repository) Update(job *models.MigrationJob) error {
	return r.db.Save(job).Error
}

// UpdateStatus updates migration job status
func (r *Repository) UpdateStatus(id uint, status models.MigrationStatus) error {
	return r.db.Model(&models.MigrationJob{}).Where("id = ?", id).
		Update("status", status).Error
}

// UpdateProgress updates migration job progress
func (r *Repository) UpdateProgress(id uint, processedFiles, failedFiles int, processedBytes int64) error {
	return r.db.Model(&models.MigrationJob{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed_files": processedFiles,
			"failed_files":    failedFiles,
			"processed_bytes": processedBytes,
		}).Error
}

// GetActiveJob gets any active migration job for user
func (r *Repository) GetActiveJob(userID uint) (*models.MigrationJob, error) {
	var job models.MigrationJob
	err := r.db.Where("user_id = ? AND status IN (?)", userID,
		[]models.MigrationStatus{models.MigrationStatusPending, models.MigrationStatusRunning}).
		First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}
