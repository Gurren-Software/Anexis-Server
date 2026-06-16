package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/Gurren-Software/Anexis-Server/apps/api/internal/infrastructure/storage"
)

type Client struct {
	basePath string
}

type Config struct {
	BasePath string
}

func NewClient(cfg *Config) (*Client, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the CWD: %w", err)
	}

	basePath := cfg.BasePath
	if basePath == "" {
		basePath = "./data/storage"
	}

	if err = os.MkdirAll(path.Join(cwd, basePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Client{
		basePath: basePath,
	}, nil
}

func (c *Client) fullPath(key string) string {
	return filepath.Join(c.basePath, key)
}

func (c *Client) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	fullPath := c.fullPath(key)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Client) UploadLarge(ctx context.Context, key string, reader io.Reader, size int64, contentType string, concurrency int) error {
	return c.Upload(ctx, key, reader, size, contentType)
}

func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := c.fullPath(key)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

func (c *Client) DownloadRange(ctx context.Context, key string, start, end int64) (io.ReadCloser, error) {
	fullPath := c.fullPath(key)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	if _, err := file.Seek(start, 0); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	return &rangeReader{
		reader: file,
		end:    end,
	}, nil
}

type rangeReader struct {
	reader *os.File
	end    int64
	read   int64
}

func (r *rangeReader) Read(p []byte) (int, error) {
	if r.read >= r.end {
		return 0, io.EOF
	}

	remaining := r.end - r.read
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err := r.reader.Read(p)
	r.read += int64(n)
	return n, err
}

func (r *rangeReader) Close() error {
	return r.reader.Close()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	fullPath := c.fullPath(key)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (c *Client) GetURL(ctx context.Context, key string, expiresIn int) (string, error) {
	return "/api/v1/storage/" + key, nil
}

func (c *Client) GetStreamURL(ctx context.Context, key string, expiresIn int) (string, error) {
	return c.GetURL(ctx, key, expiresIn)
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := c.fullPath(key)
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *Client) GetMetadata(ctx context.Context, key string) (*storage.FileMetadata, error) {
	fullPath := c.fullPath(key)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	modTime := info.ModTime()
	return &storage.FileMetadata{
		Key:          key,
		Size:         info.Size(),
		LastModified: modTime.Unix(),
	}, nil
}

func (c *Client) List(ctx context.Context, prefix string, maxKeys int) ([]*storage.FileMetadata, error) {
	fullPath := c.fullPath(prefix)

	var files []*storage.FileMetadata

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	count := 0
	for _, entry := range entries {
		if maxKeys > 0 && count >= maxKeys {
			break
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		relPath := path.Join(prefix, entry.Name())
		files = append(files, &storage.FileMetadata{
			Key:          relPath,
			Size:         info.Size(),
			LastModified: info.ModTime().Unix(),
		})
		count++
	}

	return files, nil
}

func (c *Client) Copy(ctx context.Context, srcKey, dstKey string) error {
	srcPath := c.fullPath(srcKey)
	dstPath := c.fullPath(dstKey)

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	sourceFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

var _ storage.Provider = (*Client)(nil)

func (c *Client) ServeFile(key string) (io.ReadCloser, int64, error) {
	fullPath := c.fullPath(key)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, fmt.Errorf("failed to stat file: %w", err)
	}

	return file, info.Size(), nil
}

func (c *Client) GetBasePath() string {
	return c.basePath
}
