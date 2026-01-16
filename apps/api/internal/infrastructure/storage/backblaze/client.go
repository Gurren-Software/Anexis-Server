package backblaze

import (
	"context"
	"fmt"
	"io"

	"github.com/Backblaze/blazer/b2"
	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
)

// Client implements the storage.Provider interface for Backblaze B2
type Client struct {
	client     *b2.Client
	bucket     *b2.Bucket
	bucketName string
}

// Config holds Backblaze B2 configuration
type Config struct {
	KeyID          string
	ApplicationKey string
	BucketName     string
}

// NewClient creates a new Backblaze B2 client
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	client, err := b2.NewClient(ctx, cfg.KeyID, cfg.ApplicationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create B2 client: %w", err)
	}

	// Get the bucket
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	var bucket *b2.Bucket
	for _, b := range buckets {
		if b.Name() == cfg.BucketName {
			bucket = b
			break
		}
	}

	if bucket == nil {
		return nil, fmt.Errorf("bucket %q not found", cfg.BucketName)
	}

	return &Client{
		client:     client,
		bucket:     bucket,
		bucketName: cfg.BucketName,
	}, nil
}

// Upload implements storage.Provider.Upload
func (c *Client) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	obj := c.bucket.Object(key)
	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return w.Close()
}

// UploadLarge implements storage.Provider.UploadLarge
func (c *Client) UploadLarge(ctx context.Context, key string, reader io.Reader, size int64, contentType string, concurrency int) error {
	obj := c.bucket.Object(key)
	w := obj.NewWriter(ctx)
	w.ConcurrentUploads = concurrency

	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return fmt.Errorf("failed to upload large file: %w", err)
	}

	return w.Close()
}

// Download implements storage.Provider.Download
func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	obj := c.bucket.Object(key)
	r := obj.NewReader(ctx)
	return r, nil
}

// DownloadRange implements storage.Provider.DownloadRange
func (c *Client) DownloadRange(ctx context.Context, key string, start, end int64) (io.ReadCloser, error) {
	obj := c.bucket.Object(key)
	r := obj.NewRangeReader(ctx, start, end-start+1)
	return r, nil
}

// Delete implements storage.Provider.Delete
func (c *Client) Delete(ctx context.Context, key string) error {
	obj := c.bucket.Object(key)
	return obj.Delete(ctx)
}

// GetURL implements storage.Provider.GetURL
func (c *Client) GetURL(ctx context.Context, key string, expiresIn int) (string, error) {
	obj := c.bucket.Object(key)
	return obj.URL(), nil
}

// GetStreamURL implements storage.Provider.GetStreamURL
func (c *Client) GetStreamURL(ctx context.Context, key string, expiresIn int) (string, error) {
	// For B2, streaming URL is same as regular URL
	return c.GetURL(ctx, key, expiresIn)
}

// Exists implements storage.Provider.Exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	obj := c.bucket.Object(key)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		// Check if it's a "not found" error
		return false, nil
	}
	return attrs.Name != "", nil
}

// GetMetadata implements storage.Provider.GetMetadata
func (c *Client) GetMetadata(ctx context.Context, key string) (*storage.FileMetadata, error) {
	obj := c.bucket.Object(key)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	return &storage.FileMetadata{
		Key:          key,
		Size:         attrs.Size,
		ContentType:  attrs.ContentType,
		LastModified: attrs.UploadTimestamp.Unix(),
	}, nil
}

// List implements storage.Provider.List
func (c *Client) List(ctx context.Context, prefix string, maxKeys int) ([]*storage.FileMetadata, error) {
	var files []*storage.FileMetadata
	iter := c.bucket.List(ctx, b2.ListPrefix(prefix))

	count := 0
	for iter.Next() {
		if maxKeys > 0 && count >= maxKeys {
			break
		}

		obj := iter.Object()
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			continue
		}

		files = append(files, &storage.FileMetadata{
			Key:          attrs.Name,
			Size:         attrs.Size,
			ContentType:  attrs.ContentType,
			LastModified: attrs.UploadTimestamp.Unix(),
		})
		count++
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// Copy implements storage.Provider.Copy
func (c *Client) Copy(ctx context.Context, srcKey, dstKey string) error {
	// B2 doesn't have native copy, so we download and re-upload
	reader, err := c.Download(ctx, srcKey)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	defer reader.Close()

	// Get source metadata
	meta, err := c.GetMetadata(ctx, srcKey)
	if err != nil {
		return fmt.Errorf("failed to get source metadata: %w", err)
	}

	return c.Upload(ctx, dstKey, reader, meta.Size, meta.ContentType)
}

// Ensure Client implements Provider interface
var _ storage.Provider = (*Client)(nil)
