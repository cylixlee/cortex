package main

import (
	"context"
	"log"

	"github.com/open-portfolios/cortex/internal/config"
	"github.com/open-portfolios/cortex/pkg/storage"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	minioClient, err := storage.NewMinIOClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect MinIO: %v", err)
	}

	if err := minioClient.CreateBucketIfNotExists(ctx); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	log.Printf("Bucket '%s' is ready", cfg.MinIOBucket)
}
