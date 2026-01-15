# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go workspace files
COPY go.work go.work.sum* ./
COPY packages/database/go.mod packages/database/go.sum ./packages/database/
COPY apps/api/go.mod apps/api/go.sum ./apps/api/

# Download dependencies
RUN cd packages/database && go mod download
RUN cd apps/api && go mod download

# Copy source code
COPY packages/ ./packages/
COPY apps/api/ ./apps/api/

# Build the application
WORKDIR /app/apps/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api-server ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api-server .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./api-server"]
