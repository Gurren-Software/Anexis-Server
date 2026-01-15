package migration

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/features/migration/providers"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/Treefle-labs/anexis-server/packages/database/models"
)

var (
	ErrJobNotFound     = errors.New("migration job not found")
	ErrActiveJobExists = errors.New("active migration job already exists")
	ErrInvalidProvider = errors.New("invalid provider")
)

// Service handles migration business logic
type Service struct {
	repo       *Repository
	storage    storage.Provider
	activeJobs map[uint]context.CancelFunc
	mu         sync.RWMutex
}

// NewService creates a new migration service
func NewService(repo *Repository, storage storage.Provider) *Service {
	return &Service{
		repo:       repo,
		storage:    storage,
		activeJobs: make(map[uint]context.CancelFunc),
	}
}

// Start starts a new migration job
func (s *Service) Start(ctx context.Context, userID uint, req *StartMigrationRequest) (*models.MigrationJob, error) {
	// Check for active job
	activeJob, err := s.repo.GetActiveJob(userID)
	if err != nil {
		return nil, err
	}
	if activeJob != nil {
		return nil, ErrActiveJobExists
	}

	job := &models.MigrationJob{
		UserID:       userID,
		Provider:     models.ProviderType(req.Provider),
		Status:       models.MigrationStatusPending,
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	}

	if err := s.repo.Create(job); err != nil {
		return nil, err
	}

	// Start background migration
	go s.runMigration(context.Background(), job.ID, userID)

	return job, nil
}

// runMigration runs the migration job in background
func (s *Service) runMigration(ctx context.Context, jobID, userID uint) {
	ctx, cancel := context.WithCancel(ctx)

	s.mu.Lock()
	s.activeJobs[jobID] = cancel
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.activeJobs, jobID)
		s.mu.Unlock()
	}()

	job, err := s.repo.FindByID(jobID)
	if err != nil || job == nil {
		return
	}

	// Update status to running
	now := time.Now()
	job.Status = models.MigrationStatusRunning
	job.StartedAt = &now
	s.repo.Update(job)

	// Get provider client
	provider, err := s.getProvider(job.Provider, job.AccessToken, job.RefreshToken)
	if err != nil {
		job.Status = models.MigrationStatusFailed
		job.LastError = err.Error()
		s.repo.Update(job)
		return
	}

	// List files from provider
	files, err := provider.ListFiles(ctx)
	if err != nil {
		job.Status = models.MigrationStatusFailed
		job.LastError = "Failed to list files: " + err.Error()
		s.repo.Update(job)
		return
	}

	job.TotalFiles = len(files)
	s.repo.Update(job)

	// Process each file
	processedFiles := 0
	failedFiles := 0
	var processedBytes int64

	for _, file := range files {
		select {
		case <-ctx.Done():
			job.Status = models.MigrationStatusCancelled
			s.repo.Update(job)
			return
		default:
		}

		err := s.processFile(ctx, provider, userID, job.ID, file)
		if err != nil {
			failedFiles++
			job.LastError = err.Error()
		} else {
			processedFiles++
			processedBytes += file.Size
		}

		s.repo.UpdateProgress(jobID, processedFiles, failedFiles, processedBytes)
	}

	// Complete job
	completedAt := time.Now()
	job.Status = models.MigrationStatusCompleted
	job.ProcessedFiles = processedFiles
	job.FailedFiles = failedFiles
	job.ProcessedBytes = processedBytes
	job.CompletedAt = &completedAt
	s.repo.Update(job)
}

func (s *Service) processFile(ctx context.Context, provider providers.Provider, userID uint, jobID uint, file *providers.FileInfo) error {
	// Download from provider
	reader, err := provider.DownloadFile(ctx, file.ID)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Upload to our storage
	storageKey := generateStorageKey(userID, file.Name)
	if err := s.storage.Upload(ctx, storageKey, reader, file.Size, file.MimeType); err != nil {
		return err
	}

	return nil
}

func (s *Service) getProvider(providerType models.ProviderType, accessToken, refreshToken string) (providers.Provider, error) {
	switch providerType {
	case models.ProviderTypeGoogle:
		return providers.NewGoogleProvider(accessToken, refreshToken), nil
	case models.ProviderTypeAmazon:
		return providers.NewAmazonProvider(accessToken, refreshToken), nil
	case models.ProviderTypeMicrosoft:
		return providers.NewMicrosoftProvider(accessToken, refreshToken), nil
	case models.ProviderTypeDropbox:
		return providers.NewDropboxProvider(accessToken, refreshToken), nil
	default:
		return nil, ErrInvalidProvider
	}
}

// Cancel cancels a running migration job
func (s *Service) Cancel(userID, jobID uint) error {
	job, err := s.repo.FindByIDAndUser(jobID, userID)
	if err != nil {
		return err
	}
	if job == nil {
		return ErrJobNotFound
	}

	s.mu.RLock()
	cancel, exists := s.activeJobs[jobID]
	s.mu.RUnlock()

	if exists {
		cancel()
	}

	job.Status = models.MigrationStatusCancelled
	return s.repo.Update(job)
}

// GetJob gets a migration job
func (s *Service) GetJob(userID, jobID uint) (*models.MigrationJob, error) {
	job, err := s.repo.FindByIDAndUser(jobID, userID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}
	return job, nil
}

// ListJobs lists user's migration jobs
func (s *Service) ListJobs(userID uint, req *ListMigrationsRequest) ([]models.MigrationJob, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	return s.repo.List(userID, req.Status, page, perPage)
}

// ToMigrationResponse converts job model to response DTO
func ToMigrationResponse(job *models.MigrationJob) *MigrationResponse {
	resp := &MigrationResponse{
		ID:             job.ID,
		Provider:       string(job.Provider),
		Status:         string(job.Status),
		TotalFiles:     job.TotalFiles,
		ProcessedFiles: job.ProcessedFiles,
		FailedFiles:    job.FailedFiles,
		TotalBytes:     job.TotalBytes,
		ProcessedBytes: job.ProcessedBytes,
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

	return resp
}

func generateStorageKey(userID uint, filename string) string {
	return "files/" + string(rune(userID)) + "/" + filename
}
