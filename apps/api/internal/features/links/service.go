package links

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLinkNotFound    = errors.New("link not found")
	ErrLinkExpired     = errors.New("link has expired")
	ErrDownloadLimit   = errors.New("download limit reached")
	ErrAccessDenied    = errors.New("access denied")
	ErrInvalidPassword = errors.New("invalid password")
)

// FileRepository interface for file operations
type FileRepository interface {
	FindByIDAndUser(id, userID uint) (*models.File, error)
	FindByID(id uint) (*models.File, error)
}

// Service handles link business logic
type Service struct {
	repo     *Repository
	fileRepo FileRepository
	storage  storage.Provider
	baseURL  string
}

// NewService creates a new links service
func NewService(repo *Repository, fileRepo FileRepository, storage storage.Provider, baseURL string) *Service {
	return &Service{
		repo:     repo,
		fileRepo: fileRepo,
		storage:  storage,
		baseURL:  baseURL,
	}
}

// Create creates a new access link
func (s *Service) Create(userID uint, req *CreateLinkRequest) (*models.Link, error) {
	// Verify file ownership
	file, err := s.fileRepo.FindByIDAndUser(req.FileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file not found")
	}

	// Generate unique token
	token, err := generateToken(32)
	if err != nil {
		return nil, err
	}

	link := &models.Link{
		UserID:       userID,
		FileID:       req.FileID,
		Token:        token,
		Type:         models.LinkType(req.Type),
		AccessType:   models.LinkAccessPublic,
		MaxDownloads: req.MaxDownloads,
		Name:         req.Name,
		Description:  req.Description,
	}

	if req.AccessType != "" {
		link.AccessType = models.LinkAccessType(req.AccessType)
	}

	// Hash password if provided
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashStr := string(hash)
		link.Password = &hashStr
	}

	// Set expiration for temporal links
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Second)
		link.ExpiresAt = &expiresAt
	}

	if err := s.repo.Create(link); err != nil {
		return nil, err
	}

	return link, nil
}

// GetByToken retrieves a link by token and validates access
func (s *Service) GetByToken(token string, password string) (*models.Link, error) {
	link, err := s.repo.FindByToken(token)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}

	// Check expiration
	if link.IsExpired() {
		return nil, ErrLinkExpired
	}

	// Check download limit
	if !link.HasDownloadsRemaining() {
		return nil, ErrDownloadLimit
	}

	// Check password
	if link.Password != nil && *link.Password != "" {
		if password == "" {
			return nil, ErrAccessDenied
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*link.Password), []byte(password)); err != nil {
			return nil, ErrInvalidPassword
		}
	}

	// Update access tracking
	_ = s.repo.UpdateLastAccessed(link.ID)

	return link, nil
}

// AccessFile provides file access via link
func (s *Service) AccessFile(ctx context.Context, token string, password string) (io.ReadCloser, *models.Link, error) {
	link, err := s.GetByToken(token, password)
	if err != nil {
		return nil, nil, err
	}

	// Get file
	file, err := s.fileRepo.FindByID(link.FileID)
	if err != nil || file == nil {
		return nil, nil, errors.New("file not found")
	}

	// Download from storage
	reader, err := s.storage.Download(ctx, file.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	// Increment download count
	_ = s.repo.IncrementDownloadCount(link.ID)

	return reader, link, nil
}

// GetStreamURL returns a streaming URL for the file
func (s *Service) GetStreamURL(ctx context.Context, token string, password string) (string, *models.Link, error) {
	link, err := s.GetByToken(token, password)
	if err != nil {
		return "", nil, err
	}

	file, err := s.fileRepo.FindByID(link.FileID)
	if err != nil || file == nil {
		return "", nil, errors.New("file not found")
	}

	url, err := s.storage.GetStreamURL(ctx, file.StoragePath, 3600) // 1 hour
	if err != nil {
		return "", nil, err
	}

	return url, link, nil
}

// List lists user's links
func (s *Service) List(userID uint, req *ListLinksRequest) ([]models.Link, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.List(userID, req.FileID, req.Type, page, perPage)
}

// Update updates a link
func (s *Service) Update(userID, linkID uint, req *UpdateLinkRequest) (*models.Link, error) {
	link, err := s.repo.FindByIDAndUser(linkID, userID)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}

	if req.Name != "" {
		link.Name = req.Name
	}
	if req.Description != "" {
		link.Description = req.Description
	}
	if req.MaxDownloads != nil {
		link.MaxDownloads = req.MaxDownloads
	}
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Second)
		link.ExpiresAt = &expiresAt
	}
	if req.Password != nil {
		if *req.Password == "" {
			link.Password = nil
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			hashStr := string(hash)
			link.Password = &hashStr
		}
	}

	if err := s.repo.Update(link); err != nil {
		return nil, err
	}

	return link, nil
}

// Delete deletes a link
func (s *Service) Delete(userID, linkID uint) error {
	link, err := s.repo.FindByIDAndUser(linkID, userID)
	if err != nil {
		return err
	}
	if link == nil {
		return ErrLinkNotFound
	}

	return s.repo.Delete(linkID)
}

// ToLinkResponse converts link model to response DTO
func (s *Service) ToLinkResponse(link *models.Link) *LinkResponse {
	resp := &LinkResponse{
		ID:             link.ID,
		Token:          link.Token,
		Type:           string(link.Type),
		AccessType:     string(link.AccessType),
		URL:            fmt.Sprintf("%s/api/v1/links/%s", s.baseURL, link.Token),
		MaxDownloads:   link.MaxDownloads,
		DownloadCount:  link.DownloadCount,
		ExpiresAt:      link.ExpiresAt,
		LastAccessedAt: link.LastAccessedAt,
		Name:           link.Name,
		Description:    link.Description,
		FileID:         link.FileID,
		CreatedAt:      link.CreatedAt.Format(time.RFC3339),
	}

	if link.File.ID != 0 {
		resp.FileName = link.File.OriginalName
	}

	return resp
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
