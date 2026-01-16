package auth

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles user database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new user
func (r *Repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail finds a user by email
func (r *Repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID
func (r *Repository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *Repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UpdatePassword updates user password
func (r *Repository) UpdatePassword(userID uuid.UUID, passwordHash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

// UpdateStorageUsed updates user's storage usage
func (r *Repository) UpdateStorageUsed(userID uuid.UUID, delta int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("storage_used", gorm.Expr("storage_used + ?", delta)).Error
}

// GetStorageUsed returns the current storage used
func (r *Repository) GetStorageUsed(userID uuid.UUID) (int64, error) {
	var user models.User
	err := r.db.Select("storage_used").First(&user, "id = ?", userID).Error
	if err != nil {
		return 0, err
	}
	return user.StorageUsed, nil
}
