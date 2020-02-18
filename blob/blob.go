package blob

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gosimple/slug"

	"github.com/gleez/pkg/contextkeys"
)

var (
	// ErrNotFound is returned when given blob is not found
	ErrNotFound = errors.New("Blob not found")

	// ErrInvalidKeyFormat is returned when blob key is in invalid format
	ErrInvalidKeyFormat = errors.New("Blob key is in invalid format")
)

var (
	storageMu sync.RWMutex
	storages  = make(map[string]Storage)
)

type Config struct {
	Type string
	S3   struct {
		Endpoint        string
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		BucketName      string
		Prefix          string
	}
	FS struct {
		Path string
	}
	Debug bool
}

type Storage interface {
	Name() string
	Enabled() bool
	Init(*Config)
}

type Blob struct {
	Size        int64
	Content     []byte
	ContentType string
}

type ListBlobs struct {
	Result []string
}

type GetBlobByKey struct {
	Key string

	Result *Blob
}

type StoreBlob struct {
	Key         string
	Content     []byte
	ContentType string
}

type DeleteBlob struct {
	Key string
}

func Init(cfg *Config) {

	for _, st := range storages {
		if st.Enabled() {
			st.Init(cfg)
		}
	}

}

func Reset() {
	storageMu.Lock()
	defer storageMu.Unlock()

	storages = make(map[string]Storage)
}

// Register register sharding storage with name
func Register(name string, storageFactory Storage) {
	storageMu.Lock()
	defer storageMu.Unlock()

	if storageFactory == nil {
		panic("register storage factory is nil")
	}

	if _, dup := storages[name]; dup {
		panic("register called twice for storage " + name)
	}

	storages[name] = storageFactory
}

// LoadStorage load storage by name
func LoadStorage(name string) (Storage, error) {
	storageFactory := storages[name]
	if storageFactory == nil {
		return nil, fmt.Errorf("cannot load storage from %s", name)
	}

	return storageFactory, nil
}

// SanitizeFileName replaces invalid characters from given filename
func SanitizeFileName(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	ext := filepath.Ext(fileName)
	if ext != "" {
		return slug.Make(fileName[0:len(fileName)-len(ext)]) + ext
	}

	return slug.Make(fileName)
}

// ValidateKey checks if key is is valid format
func ValidateKey(key string) error {
	if len(key) == 0 || len(key) > 512 || strings.Contains(key, " ") {
		return ErrInvalidKeyFormat
	}

	if strings.HasPrefix(key, "/") || strings.HasSuffix(key, "/") {
		return ErrInvalidKeyFormat
	}
	return nil
}

func TenantFromContext(ctx context.Context) int64 {

	if tid, ok := ctx.Value(contextkeys.Tenant).(int64); ok {
		//fmt.Printf("Tenant Id from Blob CONTEXT %d \n", tid)

		return tid
	}

	return 0
}
