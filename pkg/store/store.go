package store

import "context"

type Store interface {
	Upload(ctx context.Context, filename string, data []byte) error
	Download(ctx context.Context, filename string) ([]byte, error)
}
