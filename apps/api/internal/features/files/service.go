package files

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound         = errors.New("file not found")
	ErrAccessDenied         = errors.New("access denied")
	ErrStorageQuotaExceeded = errors.New("storage quota exceeded")
	ErrInvalidFile          = errors.New("invalid file")
)

// Service handles file business logic
type Service struct {
	repo     *Repository
	storage  storage.Provider
	authRepo AuthRepository
}

// AuthRepository interface for user operations
type AuthRepository interface {
	UpdateStorageUsed(userID uint, delta int64) error
}

// NewService creates a new files service
func NewService(repo *Repository, storage storage.Provider, authRepo AuthRepository) *Service {
	return &Service{
		repo:     repo,
		storage:  storage,
		authRepo: authRepo,
	}
}

// Upload handles file upload
func (s *Service) Upload(ctx context.Context, userID uint, storageQuota, storageUsed int64, file *multipart.FileHeader, opts *UploadRequest) (*models.File, error) {
	// Check storage quota
	if storageUsed+file.Size > storageQuota {
		return nil, ErrStorageQuotaExceeded
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file contents
	contents, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate checksum
	hash := sha256.Sum256(contents)
	checksum := hex.EncodeToString(hash[:])

	// Process file (compress if requested)
	var processedData []byte
	isCompressed := false
	if opts.Compress {
		var buf bytes.Buffer
		gzWriter := gzip.NewWriter(&buf)
		if _, err := gzWriter.Write(contents); err != nil {
			return nil, fmt.Errorf("failed to compress: %w", err)
		}
		gzWriter.Close()
		processedData = buf.Bytes()
		isCompressed = true
	} else {
		processedData = contents
	}

	// Generate storage key
	ext := filepath.Ext(file.Filename)
	storageKey := fmt.Sprintf("%d/%s%s", userID, uuid.New().String(), ext)
	storagePath := fmt.Sprintf("files/%s", storageKey)

	// Upload to storage
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	err = s.storage.Upload(ctx, storagePath, bytes.NewReader(processedData), int64(len(processedData)), contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}

	// Create database record
	fileRecord := &models.File{
		UserID:       userID,
		Name:         file.Filename,
		OriginalName: file.Filename,
		MimeType:     contentType,
		Size:         file.Size,
		StorageKey:   storageKey,
		StoragePath:  storagePath,
		Checksum:     checksum,
		Status:       models.FileStatusReady,
		IsEncrypted:  opts.Encrypt,
		IsCompressed: isCompressed,
		ParentID:     opts.ParentID,
		Description:  opts.Description,
		Tags:         opts.Tags,
		UploadedAt:   time.Now(),
	}

	if err := s.repo.Create(fileRecord); err != nil {
		// Try to clean up storage
		_ = s.storage.Delete(ctx, storagePath)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// Update user storage used
	if err := s.authRepo.UpdateStorageUsed(userID, file.Size); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to update storage usage: %v\n", err)
	}

	return fileRecord, nil
}

// Download retrieves a file for download
func (s *Service) Download(ctx context.Context, userID, fileID uint) (io.ReadCloser, *models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, nil, err
	}
	if file == nil {
		return nil, nil, ErrFileNotFound
	}

	reader, err := s.storage.Download(ctx, file.StoragePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download from storage: %w", err)
	}

	// If compressed, decompress
	if file.IsCompressed {
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			reader.Close()
			return nil, nil, fmt.Errorf("failed to decompress: %w", err)
		}
		return gzReader, file, nil
	}

	return reader, file, nil
}

// GetFile retrieves file metadata
func (s *Service) GetFile(userID, fileID uint) (*models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

// ListFiles lists user's files
func (s *Service) ListFiles(userID uint, req *ListFilesRequest) ([]models.File, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	return s.repo.List(userID, req.ParentID, req.Search, page, perPage)
}

// CreateFolder creates a new folder
func (s *Service) CreateFolder(userID uint, req *CreateFolderRequest) (*models.File, error) {
	folder := &models.File{
		UserID:       userID,
		Name:         req.Name,
		OriginalName: req.Name,
		MimeType:     "application/x-directory",
		Size:         0,
		StorageKey:   fmt.Sprintf("%d/folder-%s", userID, uuid.New().String()),
		StoragePath:  "",
		Checksum:     "",
		Status:       models.FileStatusReady,
		ParentID:     req.ParentID,
		UploadedAt:   time.Now(),
	}

	if err := s.repo.Create(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// Rename renames a file or folder
func (s *Service) Rename(userID, fileID uint, req *RenameRequest) (*models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}

	file.Name = req.Name
	if err := s.repo.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}

// Move moves a file to a different folder
func (s *Service) Move(userID, fileID uint, req *MoveRequest) (*models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}

	// Validate target folder exists (if specified)
	if req.TargetParentID != nil {
		targetFolder, err := s.repo.FindByIDAndUser(*req.TargetParentID, userID)
		if err != nil {
			return nil, err
		}
		if targetFolder == nil || !targetFolder.IsFolder() {
			return nil, errors.New("invalid target folder")
		}
	}

	file.ParentID = req.TargetParentID
	if err := s.repo.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}

// Delete deletes a file
func (s *Service) Delete(ctx context.Context, userID, fileID uint) error {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return err
	}
	if file == nil {
		return ErrFileNotFound
	}

	// Delete from storage (if not a folder)
	if !file.IsFolder() && file.StoragePath != "" {
		if err := s.storage.Delete(ctx, file.StoragePath); err != nil {
			// Log but continue with deletion
			fmt.Printf("Warning: failed to delete from storage: %v\n", err)
		}
	}

	// Update user storage
	if !file.IsFolder() {
		if err := s.authRepo.UpdateStorageUsed(userID, -file.Size); err != nil {
			fmt.Printf("Warning: failed to update storage usage: %v\n", err)
		}
	}

	return s.repo.Delete(fileID)
}

// ToFileResponse converts file model to response DTO
func ToFileResponse(file *models.File) *FileResponse {
	return &FileResponse{
		ID:           file.ID,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		MimeType:     file.MimeType,
		Size:         file.Size,
		Status:       string(file.Status),
		IsEncrypted:  file.IsEncrypted,
		IsCompressed: file.IsCompressed,
		Description:  file.Description,
		Tags:         file.Tags,
		UploadedAt:   file.UploadedAt.Format(time.RFC3339),
		ParentID:     file.ParentID,
	}
}
