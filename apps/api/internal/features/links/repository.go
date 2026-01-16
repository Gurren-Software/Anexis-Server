package links

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles link database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new links repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new link
func (r *Repository) Create(link *models.Link) error {
	return r.db.Create(link).Error
}

// FindByID finds a link by ID
func (r *Repository) FindByID(id uuid.UUID) (*models.Link, error) {
	var link models.Link
	err := r.db.Preload("File").First(&link, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

// FindByToken finds a link by token
func (r *Repository) FindByToken(token string) (*models.Link, error) {
	var link models.Link
	err := r.db.Preload("File").Where("token = ?", token).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

// FindByIDAndUser finds a link by ID and user
func (r *Repository) FindByIDAndUser(id, userID uuid.UUID) (*models.Link, error) {
	var link models.Link
	err := r.db.Preload("File").Where("id = ? AND user_id = ?", id, userID).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

// List lists links for a user
func (r *Repository) List(userID uuid.UUID, fileID *uuid.UUID, linkType string, page, perPage int) ([]models.Link, int64, error) {
	var links []models.Link
	var total int64

	query := r.db.Model(&models.Link{}).Where("user_id = ?", userID)

	if fileID != nil {
		query = query.Where("file_id = ?", *fileID)
	}
	if linkType != "" {
		query = query.Where("type = ?", linkType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Preload("File").Offset(offset).Limit(perPage).
		Order("created_at DESC").Find(&links).Error; err != nil {
		return nil, 0, err
	}

	return links, total, nil
}

// Update updates a link
func (r *Repository) Update(link *models.Link) error {
	return r.db.Save(link).Error
}

// IncrementDownloadCount increments download count
func (r *Repository) IncrementDownloadCount(id uuid.UUID) error {
	return r.db.Model(&models.Link{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

// UpdateLastAccessed updates last accessed timestamp
func (r *Repository) UpdateLastAccessed(id uuid.UUID) error {
	return r.db.Model(&models.Link{}).Where("id = ?", id).
		Update("last_accessed_at", gorm.Expr("NOW()")).Error
}

// Delete deletes a link
func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Link{}, "id = ?", id).Error
}
