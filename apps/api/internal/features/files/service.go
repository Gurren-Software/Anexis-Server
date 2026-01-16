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
	"time"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrStorageExceeded = errors.New("storage quota exceeded")
)

// AuthRepository interface for auth operations
type AuthRepository interface {
	UpdateStorageUsed(userID uuid.UUID, delta int64) error
	GetStorageUsed(userID uuid.UUID) (int64, error)
}

// Service handles file business logic
type Service struct {
	repo     *Repository
	storage  storage.Provider
	authRepo AuthRepository
}

// NewService creates a new files service
func NewService(repo *Repository, storage storage.Provider, authRepo AuthRepository) *Service {
	return &Service{
		repo:     repo,
		storage:  storage,
		authRepo: authRepo,
	}
}

// Upload uploads a file
func (s *Service) Upload(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, compress bool, parentID *uuid.UUID, description string, storageQuota int64) (*models.File, error) {
	// Check storage quota
	currentUsed, err := s.authRepo.GetStorageUsed(userID)
	if err != nil {
		return nil, err
	}
	if currentUsed+file.Size > storageQuota {
		return nil, ErrStorageExceeded
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read file content
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// Calculate checksum
	hash := sha256.Sum256(content)
	checksum := hex.EncodeToString(hash[:])

	// Compress if requested
	var fileContent []byte
	isCompressed := false
	if compress {
		var buf bytes.Buffer
		gzWriter := gzip.NewWriter(&buf)
		if _, err := gzWriter.Write(content); err != nil {
			return nil, err
		}
		if err := gzWriter.Close(); err != nil {
			return nil, err
		}
		fileContent = buf.Bytes()
		isCompressed = true
	} else {
		fileContent = content
	}

	// Generate storage key
	storageKey := fmt.Sprintf("files/%s/%s/%s", userID.String(), time.Now().Format("2006/01/02"), uuid.New().String())

	// Upload to storage
	if err := s.storage.Upload(ctx, storageKey, bytes.NewReader(fileContent), int64(len(fileContent)), file.Header.Get("Content-Type")); err != nil {
		return nil, err
	}

	// Create file record
	now := time.Now()
	fileRecord := &models.File{
		UserID:       userID,
		Name:         file.Filename,
		OriginalName: file.Filename,
		MimeType:     file.Header.Get("Content-Type"),
		Size:         file.Size,
		StorageKey:   storageKey,
		StoragePath:  storageKey,
		Checksum:     checksum,
		Status:       models.FileStatusReady,
		IsCompressed: isCompressed,
		ParentID:     parentID,
		Description:  description,
		UploadedAt:   &now,
		ProcessedAt:  &now,
	}

	if err := s.repo.Create(fileRecord); err != nil {
		return nil, err
	}

	// Update storage used
	if err := s.authRepo.UpdateStorageUsed(userID, file.Size); err != nil {
		return nil, err
	}

	return fileRecord, nil
}

// Download downloads a file
func (s *Service) Download(ctx context.Context, userID, fileID uuid.UUID) (io.ReadCloser, *models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, nil, err
	}
	if file == nil {
		return nil, nil, ErrFileNotFound
	}

	reader, err := s.storage.Download(ctx, file.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	// If file was compressed, decompress it
	if file.IsCompressed {
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			reader.Close()
			return nil, nil, err
		}
		// Wrap to close both readers
		return &decompressReader{gzReader: gzReader, underlying: reader}, file, nil
	}

	return reader, file, nil
}

// decompressReader wraps gzip reader to close underlying reader too
type decompressReader struct {
	gzReader   *gzip.Reader
	underlying io.ReadCloser
}

func (d *decompressReader) Read(p []byte) (n int, err error) {
	return d.gzReader.Read(p)
}

func (d *decompressReader) Close() error {
	d.gzReader.Close()
	return d.underlying.Close()
}

// GetFile gets file metadata
func (s *Service) GetFile(userID, fileID uuid.UUID) (*models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

// ListFiles lists files
func (s *Service) ListFiles(userID uuid.UUID, req *ListFilesRequest) ([]models.File, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.List(userID, req.ParentID, req.Search, page, perPage)
}

// CreateFolder creates a folder
func (s *Service) CreateFolder(userID uuid.UUID, req *CreateFolderRequest) (*models.File, error) {
	folder := &models.File{
		UserID:       userID,
		Name:         req.Name,
		OriginalName: req.Name,
		MimeType:     "application/x-directory",
		Size:         0,
		StorageKey:   fmt.Sprintf("folders/%s/%s", userID.String(), uuid.New().String()),
		StoragePath:  "",
		Checksum:     "",
		Status:       models.FileStatusReady,
		ParentID:     req.ParentID,
	}

	if err := s.repo.Create(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// Rename renames a file or folder
func (s *Service) Rename(userID, fileID uuid.UUID, newName string) (*models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}

	file.Name = newName
	if err := s.repo.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}

// Move moves a file or folder
func (s *Service) Move(userID, fileID uuid.UUID, newParentID *uuid.UUID) (*models.File, error) {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}

	file.ParentID = newParentID
	if err := s.repo.Update(file); err != nil {
		return nil, err
	}

	return file, nil
}

// Delete deletes a file
func (s *Service) Delete(ctx context.Context, userID, fileID uuid.UUID) error {
	file, err := s.repo.FindByIDAndUser(fileID, userID)
	if err != nil {
		return err
	}
	if file == nil {
		return ErrFileNotFound
	}

	// Delete from storage if not a folder
	if !file.IsFolder() {
		if err := s.storage.Delete(ctx, file.StoragePath); err != nil {
			return err
		}

		// Update storage used
		if err := s.authRepo.UpdateStorageUsed(userID, -file.Size); err != nil {
			return err
		}
	}

	return s.repo.Delete(fileID)
}

// ToFileResponse converts file model to response
func (s *Service) ToFileResponse(file *models.File) *FileResponse {
	return &FileResponse{
		ID:           file.ID,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		MimeType:     file.MimeType,
		Size:         file.Size,
		IsFolder:     file.IsFolder(),
		ParentID:     file.ParentID,
		Description:  file.Description,
		IsEncrypted:  file.IsEncrypted,
		IsCompressed: file.IsCompressed,
		CreatedAt:    file.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    file.UpdatedAt.Format(time.RFC3339),
	}
}
