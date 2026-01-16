package backup

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles backup job database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new backup repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new backup job
func (r *Repository) Create(job *models.BackupJob) error {
	return r.db.Create(job).Error
}

// FindByID finds a backup job by ID
func (r *Repository) FindByID(id uuid.UUID) (*models.BackupJob, error) {
	var job models.BackupJob
	err := r.db.First(&job, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// FindByIDAndUser finds a backup job by ID and user
func (r *Repository) FindByIDAndUser(id, userID uuid.UUID) (*models.BackupJob, error) {
	var job models.BackupJob
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// List lists backup jobs for a user
func (r *Repository) List(userID uuid.UUID, backupType, status string, page, perPage int) ([]models.BackupJob, int64, error) {
	var jobs []models.BackupJob
	var total int64

	query := r.db.Model(&models.BackupJob{}).Where("user_id = ?", userID)

	if backupType != "" {
		query = query.Where("type = ?", backupType)
	}
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

// Update updates a backup job
func (r *Repository) Update(job *models.BackupJob) error {
	return r.db.Save(job).Error
}

// GetActiveJob gets any active backup job for user
func (r *Repository) GetActiveJob(userID uuid.UUID) (*models.BackupJob, error) {
	var job models.BackupJob
	err := r.db.Where("user_id = ? AND status IN (?)", userID,
		[]models.BackupStatus{models.BackupStatusPending, models.BackupStatusRunning}).
		First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}
