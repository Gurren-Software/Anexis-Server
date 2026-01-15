package auth

import (
	"errors"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
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
func (r *Repository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
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
func (r *Repository) UpdatePassword(userID uint, passwordHash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

// UpdateStorageUsed updates user's storage usage
func (r *Repository) UpdateStorageUsed(userID uint, delta int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("storage_used", gorm.Expr("storage_used + ?", delta)).Error
}

// EmailExists checks if email is already registered
func (r *Repository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
