package gcp

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/libnat/infra/value"

	"cloud.google.com/go/storage"
)

var DefaultACL = []storage.ACLRule{{
	Entity: storage.AllUsers,
	Role:   storage.RoleReader,
}}

type Bucket struct {
	name         string
	handle       *storage.BucketHandle
	acl          []storage.ACLRule
	cacheControl string
	baseURL      string
}

func NewBucket(name, baseURL, cacheControl string, handle *storage.BucketHandle, acl []storage.ACLRule) *Bucket {
	return &Bucket{
		name:         name,
		handle:       handle,
		acl:          acl,
		cacheControl: cacheControl,
		baseURL:      baseURL,
	}
}

func (b *Bucket) Name() string {
	return b.name
}

func (b *Bucket) Save(ctx context.Context, obj *value.Object) (string, error) {
	wc := b.handle.Object(obj.Name).NewWriter(ctx)
	wc.ACL = b.acl
	wc.ContentType = obj.MIMEType
	wc.CacheControl = b.cacheControl
	if _, err := wc.Write(obj.Content); err != nil {
		return "", fmt.Errorf("write: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("close: %w", err)
	}
	return filepath.Join(b.baseURL, obj.Name), nil
}
