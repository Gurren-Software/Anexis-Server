package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/Treefle-labs/anexis-server/apps/api/internal/infrastructure/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client     *s3.Client
	bucket     string
	basePath   string
	presignURL func(ctx context.Context, key string, expires time.Duration) (string, error)
}

type Config struct {
	Endpoint       string
	Region         string
	Bucket         string
	AccessKey      string
	SecretKey      string
	ForcePathStyle bool
	BasePath       string
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var endpoint *url.URL
	if cfg.Endpoint != "" {
		endpoint, err = url.Parse(cfg.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid endpoint URL: %w", err)
		}
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = cfg.ForcePathStyle
		if endpoint != nil && endpoint.Scheme == "http" {
			o.UseAccelerate = false
		}
	})

	presignClient := s3.NewPresignClient(s3Client)

	return &Client{
		client:   s3Client,
		bucket:   cfg.Bucket,
		basePath: cfg.BasePath,
		presignURL: func(ctx context.Context, key string, expires time.Duration) (string, error) {
			req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(cfg.Bucket),
				Key:    aws.String(key),
			}, func(opts *s3.PresignOptions) {
				opts.Expires = expires
			})
			if err != nil {
				return "", err
			}
			return req.URL, nil
		},
	}, nil
}

func (c *Client) fullKey(key string) string {
	if c.basePath == "" {
		return key
	}
	return path.Join(c.basePath, key)
}

func (c *Client) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(c.fullKey(key)),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (c *Client) UploadLarge(ctx context.Context, key string, reader io.Reader, size int64, contentType string, concurrency int) error {
	return c.Upload(ctx, key, reader, size, contentType)
}

func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.fullKey(key)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return resp.Body, nil
}

func (c *Client) DownloadRange(ctx context.Context, key string, start, end int64) (io.ReadCloser, error) {
	resp, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.fullKey(key)),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", start, end)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download range: %w", err)
	}
	return resp.Body, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.fullKey(key)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (c *Client) GetURL(ctx context.Context, key string, expiresIn int) (string, error) {
	return c.presignURL(ctx, c.fullKey(key), time.Duration(expiresIn)*time.Second)
}

func (c *Client) GetStreamURL(ctx context.Context, key string, expiresIn int) (string, error) {
	return c.GetURL(ctx, key, expiresIn)
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.fullKey(key)),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (c *Client) GetMetadata(ctx context.Context, key string) (*storage.FileMetadata, error) {
	resp, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.fullKey(key)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	return &storage.FileMetadata{
		Key:         key,
		Size:        *resp.ContentLength,
		ContentType: *resp.ContentType,
	}, nil
}

func (c *Client) List(ctx context.Context, prefix string, maxKeys int) ([]*storage.FileMetadata, error) {
	var files []*storage.FileMetadata

	result, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(c.bucket),
		Prefix:  aws.String(c.fullKey(prefix)),
		MaxKeys: aws.Int32(int32(maxKeys)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	for _, obj := range result.Contents {
		files = append(files, &storage.FileMetadata{
			Key:          *obj.Key,
			Size:         *obj.Size,
			LastModified: obj.LastModified.Unix(),
		})
	}

	return files, nil
}

func (c *Client) Copy(ctx context.Context, srcKey, dstKey string) error {
	_, err := c.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(c.bucket),
		CopySource: aws.String(c.bucket + "/" + c.fullKey(srcKey)),
		Key:        aws.String(c.fullKey(dstKey)),
	})
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

var _ storage.Provider = (*Client)(nil)
