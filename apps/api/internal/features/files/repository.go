package files

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"gorm.io/gorm"
)

// Repository handles file database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new files repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new file record
func (r *Repository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// FindByID finds a file by ID
func (r *Repository) FindByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// FindByIDAndUser finds a file by ID and user ID
func (r *Repository) FindByIDAndUser(id, userID uint) (*models.File, error) {
	var file models.File
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// FindByStorageKey finds a file by storage key
func (r *Repository) FindByStorageKey(key string) (*models.File, error) {
	var file models.File
	err := r.db.Where("storage_key = ?", key).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// List lists files for a user with optional filtering
func (r *Repository) List(userID uint, parentID *uint, search string, page, perPage int) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.Model(&models.File{}).Where("user_id = ?", userID)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR original_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Update updates a file record
func (r *Repository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

// UpdateStatus updates file status
func (r *Repository) UpdateStatus(id uint, status models.FileStatus) error {
	return r.db.Model(&models.File{}).Where("id = ?", id).
		Update("status", status).Error
}

// Delete soft-deletes a file
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

// HardDelete permanently deletes a file
func (r *Repository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&models.File{}, id).Error
}

// GetUserFiles gets all files for a user (for backup)
func (r *Repository) GetUserFiles(userID uint) ([]models.File, error) {
	var files []models.File
	err := r.db.Where("user_id = ? AND status = ?", userID, models.FileStatusReady).
		Find(&files).Error
	return files, err
}

// GetTotalSizeByUser calculates total storage used by user
func (r *Repository) GetTotalSizeByUser(userID uint) (int64, error) {
	var total int64
	err := r.db.Model(&models.File{}).
		Where("user_id = ? AND status = ?", userID, models.FileStatusReady).
		Select("COALESCE(SUM(size), 0)").
		Scan(&total).Error
	return total, err
}
