package main

import (
	"context"
	"log"

	"github.com/cylixlee/cortex/internal/config"
	"github.com/cylixlee/cortex/pkg/storage"
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
