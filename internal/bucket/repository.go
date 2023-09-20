package bucket

import "io"

type FileStorer interface {
	Store(key, mime string, file io.Reader) error
	Fetch(key string) (io.ReadCloser, error)
	Destroy(key string) error
	DestroyMany(keys []string) error
}
