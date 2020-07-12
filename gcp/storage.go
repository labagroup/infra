package gcp

import (
	"context"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/gopub/log"
)

const (
	cacheControl = "public, max-age=14400"
)

type Storage struct {
	client  *storage.Client
	baseURL string
}

func NewStorage(baseURL string) *Storage {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Save client failed: %v", err)
	}
	return &Storage{client: client, baseURL: baseURL}
}

func (s *Storage) Bucket(name string, acl []storage.ACLRule, cacheControl string) *Bucket {
	baseURL := filepath.Join(s.baseURL, name)
	handle := s.client.Bucket(name)
	return NewBucket(name, baseURL, cacheControl, handle, acl)
}
