---
sidebar_position: 1
---

# Introduction

Anexis Server is an open-source cloud file storage API server built with Go. It provides a complete solution for file management, sharing, and backup operations.

## Features

- **Multi-provider Storage**: Store files locally, on Backblaze B2, or any S3-compatible storage
- **File Management**: Upload, download, organize files in folders, rename, move, delete
- **Access Links**: Share files via permanent, temporal, streaming, or download links
- **Cloud Migration**: Import files from Google Drive, Amazon S3, OneDrive, Dropbox
- **Backup & Restore**: Export your data as ZIP archives, restore from backups
- **Dual Deployment Modes**: SaaS mode (multi-user with quotas) or Standalone mode (self-hosted)

## Architecture

Anexis Server follows a clean architecture pattern with vertical slices:

```
apps/api/
├── cmd/server/          # Entry point
└── internal/
    ├── config/          # Configuration
    ├── features/        # Business logic (auth, files, links, migration, backup)
    └── infrastructure/  # HTTP, storage providers
packages/database/       # Database models and migrations
```

## Quick Links

- [Getting Started](./getting-started)
- [Configuration](./configuration)
- [Deployment](./deployment/docker)
- [API Reference](./api/auth)