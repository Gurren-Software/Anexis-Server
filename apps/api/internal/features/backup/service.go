package backup

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/google/uuid"
)

var (
	ErrJobNotFound     = errors.New("backup job not found")
	ErrActiveJobExists = errors.New("active backup job already exists")
)

// FileRepository interface for file operations
type FileRepository interface {
	GetUserFiles(userID uint) ([]models.File, error)
}

// Service handles backup business logic
type Service struct {
	repo     *Repository
	fileRepo FileRepository
	storage  storage.Provider
	baseURL  string
}

// NewService creates a new backup service
func NewService(repo *Repository, fileRepo FileRepository, storage storage.Provider, baseURL string) *Service {
	return &Service{
		repo:     repo,
		fileRepo: fileRepo,
		storage:  storage,
		baseURL:  baseURL,
	}
}

// StartExport starts a new export backup job
func (s *Service) StartExport(ctx context.Context, userID uint) (*models.BackupJob, error) {
	// Check for active job
	activeJob, err := s.repo.GetActiveJob(userID)
	if err != nil {
		return nil, err
	}
	if activeJob != nil {
		return nil, ErrActiveJobExists
	}

	job := &models.BackupJob{
		UserID: userID,
		Type:   models.BackupTypeExport,
		Status: models.BackupStatusPending,
	}

	if err := s.repo.Create(job); err != nil {
		return nil, err
	}

	// Start background export
	go s.runExport(context.Background(), job.ID, userID)

	return job, nil
}

func (s *Service) runExport(ctx context.Context, jobID, userID uint) {
	job, err := s.repo.FindByID(jobID)
	if err != nil || job == nil {
		return
	}

	now := time.Now()
	job.Status = models.BackupStatusRunning
	job.StartedAt = &now
	s.repo.Update(job)

	// Get all user files
	files, err := s.fileRepo.GetUserFiles(userID)
	if err != nil {
		job.Status = models.BackupStatusFailed
		job.LastError = err.Error()
		s.repo.Update(job)
		return
	}

	job.TotalFiles = len(files)
	s.repo.Update(job)

	// Create ZIP archive
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	processedFiles := 0
	for _, file := range files {
		if file.IsFolder() {
			continue
		}

		reader, err := s.storage.Download(ctx, file.StoragePath)
		if err != nil {
			job.LastError = fmt.Sprintf("Failed to download %s: %v", file.Name, err)
			continue
		}

		writer, err := zipWriter.Create(file.OriginalName)
		if err != nil {
			reader.Close()
			continue
		}

		_, err = io.Copy(writer, reader)
		reader.Close()
		if err != nil {
			continue
		}

		processedFiles++
		job.ProcessedFiles = processedFiles
		s.repo.Update(job)
	}

	if err := zipWriter.Close(); err != nil {
		job.Status = models.BackupStatusFailed
		job.LastError = err.Error()
		s.repo.Update(job)
		return
	}

	// Upload archive to storage
	archiveKey := fmt.Sprintf("backups/%d/%s.zip", userID, uuid.New().String())
	archiveData := buf.Bytes()

	if err := s.storage.Upload(ctx, archiveKey, bytes.NewReader(archiveData), int64(len(archiveData)), "application/zip"); err != nil {
		job.Status = models.BackupStatusFailed
		job.LastError = err.Error()
		s.repo.Update(job)
		return
	}

	// Complete job
	completedAt := time.Now()
	expiresAt := completedAt.Add(7 * 24 * time.Hour) // 7 days
	job.Status = models.BackupStatusCompleted
	job.ArchiveKey = archiveKey
	job.ArchiveSize = int64(len(archiveData))
	job.CompletedAt = &completedAt
	job.ExpiresAt = &expiresAt
	s.repo.Update(job)
}

// GetJob gets a backup job
func (s *Service) GetJob(userID, jobID uint) (*models.BackupJob, error) {
	job, err := s.repo.FindByIDAndUser(jobID, userID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}
	return job, nil
}

// ListJobs lists user's backup jobs
func (s *Service) ListJobs(userID uint, req *ListBackupsRequest) ([]models.BackupJob, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.List(userID, req.Type, req.Status, page, perPage)
}

// GetDownloadURL returns download URL for a backup
func (s *Service) GetDownloadURL(ctx context.Context, userID, jobID uint) (string, error) {
	job, err := s.repo.FindByIDAndUser(jobID, userID)
	if err != nil {
		return "", err
	}
	if job == nil {
		return "", ErrJobNotFound
	}

	if job.Status != models.BackupStatusCompleted || job.ArchiveKey == "" {
		return "", errors.New("backup not ready for download")
	}

	// Check expiration
	if job.ExpiresAt != nil && time.Now().After(*job.ExpiresAt) {
		return "", errors.New("backup has expired")
	}

	return s.storage.GetURL(ctx, job.ArchiveKey, 3600) // 1 hour URL
}

// ToBackupResponse converts job model to response DTO
func (s *Service) ToBackupResponse(job *models.BackupJob) *BackupResponse {
	resp := &BackupResponse{
		ID:             job.ID,
		Type:           string(job.Type),
		Status:         string(job.Status),
		ArchiveKey:     job.ArchiveKey,
		ArchiveSize:    job.ArchiveSize,
		TotalFiles:     job.TotalFiles,
		ProcessedFiles: job.ProcessedFiles,
		Progress:       job.Progress(),
		LastError:      job.LastError,
		CreatedAt:      job.CreatedAt.Format(time.RFC3339),
	}

	if job.StartedAt != nil {
		resp.StartedAt = job.StartedAt.Format(time.RFC3339)
	}
	if job.CompletedAt != nil {
		resp.CompletedAt = job.CompletedAt.Format(time.RFC3339)
	}
	if job.ExpiresAt != nil {
		resp.ExpiresAt = job.ExpiresAt.Format(time.RFC3339)
	}

	return resp
}
