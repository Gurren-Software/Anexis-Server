package s3

import (
	"context"
	"testing"
	"time"
)

func TestFullKey(t *testing.T) {
	client := &Client{basePath: "tenant/root"}

	if got := client.fullKey("folder/file.txt"); got != "tenant/root/folder/file.txt" {
		t.Fatalf("unexpected full key %q", got)
	}

	client.basePath = ""
	if got := client.fullKey("folder/file.txt"); got != "folder/file.txt" {
		t.Fatalf("unexpected full key without base path %q", got)
	}
}

func TestGetURLUsesPresigner(t *testing.T) {
	client := &Client{
		basePath: "base",
		presignURL: func(ctx context.Context, key string, expires time.Duration) (string, error) {
			if key != "base/file.txt" {
				t.Fatalf("unexpected presign key %q", key)
			}
			if expires != time.Minute {
				t.Fatalf("unexpected expiration %s", expires)
			}
			return "https://signed.example.com/file.txt", nil
		},
	}

	got, err := client.GetURL(context.Background(), "file.txt", 60)
	if err != nil {
		t.Fatalf("GetURL failed: %v", err)
	}
	if got != "https://signed.example.com/file.txt" {
		t.Fatalf("unexpected url %q", got)
	}
}
