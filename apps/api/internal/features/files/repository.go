package files

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
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

// Create creates a new file
func (r *Repository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// FindByID finds a file by ID
func (r *Repository) FindByID(id uuid.UUID) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// FindByIDAndUser finds a file by ID and user
func (r *Repository) FindByIDAndUser(id, userID uuid.UUID) (*models.File, error) {
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

// List lists files for a user
func (r *Repository) List(userID uuid.UUID, parentID *uuid.UUID, search string, page, perPage int) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.Model(&models.File{}).Where("user_id = ?", userID)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Update updates a file
func (r *Repository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

// Delete deletes a file
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.File{}, "id = ?", id).Error
}

// GetUserStorageUsed calculates total storage used by user
func (r *Repository) GetUserStorageUsed(userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&models.File{}).
		Where("user_id = ? AND mime_type != ?", userID, "application/x-directory").
		Select("COALESCE(SUM(size), 0)").Scan(&total).Error
	return total, err
}

// GetUserFiles gets all files for a user
func (r *Repository) GetUserFiles(userID uuid.UUID) ([]models.File, error) {
	var files []models.File
	err := r.db.Where("user_id = ?", userID).Find(&files).Error
	return files, err
}
