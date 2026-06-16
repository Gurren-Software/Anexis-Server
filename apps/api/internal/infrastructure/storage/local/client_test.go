package local

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestClientFileLifecycle(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(&Config{BasePath: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	key := "users/123/file.txt"
	content := "hello anexis"

	if err := client.Upload(ctx, key, strings.NewReader(content), int64(len(content)), "text/plain"); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	exists, err := client.Exists(ctx, key)
	if err != nil {
		t.Fatalf("exists failed: %v", err)
	}
	if !exists {
		t.Fatalf("expected uploaded file to exist")
	}

	reader, err := client.Download(ctx, key)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(data) != content {
		t.Fatalf("expected %q, got %q", content, string(data))
	}

	metadata, err := client.GetMetadata(ctx, key)
	if err != nil {
		t.Fatalf("metadata failed: %v", err)
	}
	if metadata.Key != key || metadata.Size != int64(len(content)) {
		t.Fatalf("unexpected metadata: %#v", metadata)
	}

	served, size, err := client.ServeFile(key)
	if err != nil {
		t.Fatalf("serve file failed: %v", err)
	}
	_ = served.Close()
	if size != int64(len(content)) {
		t.Fatalf("expected served size %d, got %d", len(content), size)
	}

	if err := client.Delete(ctx, key); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	exists, err = client.Exists(ctx, key)
	if err != nil {
		t.Fatalf("exists after delete failed: %v", err)
	}
	if exists {
		t.Fatalf("expected deleted file to not exist")
	}
}

func TestClientRangeCopyAndList(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(&Config{BasePath: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := client.Upload(ctx, "docs/source.txt", strings.NewReader("abcdef"), 6, "text/plain"); err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	rangeReader, err := client.DownloadRange(ctx, "docs/source.txt", 2, 3)
	if err != nil {
		t.Fatalf("download range failed: %v", err)
	}
	rangeData, err := io.ReadAll(rangeReader)
	_ = rangeReader.Close()
	if err != nil {
		t.Fatalf("read range failed: %v", err)
	}
	if string(rangeData) != "cde" {
		t.Fatalf("expected range cde, got %q", string(rangeData))
	}

	if err := client.Copy(ctx, "docs/source.txt", "docs/copy.txt"); err != nil {
		t.Fatalf("copy failed: %v", err)
	}

	files, err := client.List(ctx, "docs", 10)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 listed files, got %d", len(files))
	}

	if got, err := client.GetURL(ctx, "docs/source.txt", 60); err != nil || got != "/api/v1/storage/docs/source.txt" {
		t.Fatalf("unexpected url %q, err %v", got, err)
	}
}
