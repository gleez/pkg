package fs

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/gleez/pkg/blob"
	"github.com/gleez/pkg/errors"
)

var (
	perm os.FileMode = 0744
	cfg  *blob.Config
)

type Service struct{}

func (s Service) Name() string {
	return "FileSystem"
}

func (s Service) Enabled() bool {
	return cfg.Type == "fs"
}

func (s Service) Init(c *blob.Config) {
	cfg = c
}

func ListBlobs(ctx context.Context, q *blob.ListBlobs) error {
	basePath := basePath(ctx)
	files := make([]string, 0)

	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path[len(basePath)+1:])
			}
			return nil
		})
	if err != nil {
		return errors.Wrapf(err, "failed to read dir '%s'", basePath)
	}

	sort.Strings(files)
	q.Result = files
	return nil
}

func GetBlobByKey(ctx context.Context, q *blob.GetBlobByKey) error {
	fullPath := keyFullPath(ctx, q.Key)
	stats, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return blob.ErrNotFound
		}
		return errors.Wrapf(err, "failed to get stats '%s' from FileSystem", q.Key)
	}

	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return errors.Wrapf(err, "failed to get read '%s' from FileSystem", q.Key)
	}

	q.Result = &blob.Blob{
		Content:     file,
		ContentType: http.DetectContentType(file),
		Size:        stats.Size(),
	}
	return nil
}

func StoreBlob(ctx context.Context, c *blob.StoreBlob) error {
	if err := blob.ValidateKey(c.Key); err != nil {
		return errors.Wrapf(err, "failed to validate blob key '%s'", c.Key)
	}

	fullPath := keyFullPath(ctx, c.Key)
	err := os.MkdirAll(filepath.Dir(fullPath), perm)

	if err != nil {
		return errors.Wrapf(err, "failed to create folder '%s' on FileSystem", fullPath)
	}

	err = ioutil.WriteFile(fullPath, c.Content, perm)
	if err != nil {
		return errors.Wrapf(err, "failed to create file '%s' on FileSystem", fullPath)
	}

	return nil
}

func DeleteBlob(ctx context.Context, c *blob.DeleteBlob) error {
	fullPath := keyFullPath(ctx, c.Key)
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed to delete file '%s' from FileSystem", c.Key)
	}
	return nil
}

func keyFullPath(ctx context.Context, key string) string {
	return path.Join(basePath(ctx), key)
}

func basePath(ctx context.Context) string {
	startPath := cfg.FS.Path
	tenant := blob.TenantFromContext(ctx)
	if tenant > 0 {
		return path.Join(startPath, "tenants", strconv.FormatInt(tenant, 10))
	}

	return path.Join(startPath)
}
