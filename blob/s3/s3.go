package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/gleez/pkg/blob"
	"github.com/gleez/pkg/errors"
)

var (
	//DefaultClient is an S3 Client
	DefaultClient *s3.S3

	cfg *blob.Config
)

type Service struct{}

func (s Service) Name() string {
	return "S3"
}

func (s Service) Enabled() bool {
	return cfg.Type == "s3"
}

func (s Service) Init(c *blob.Config) {
	cfg = c
	s3EnvConfig := cfg.S3

	if s3EnvConfig.Endpoint != "" {
		s3Config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(s3EnvConfig.AccessKeyID, s3EnvConfig.SecretAccessKey, ""),
			Endpoint:         aws.String(s3EnvConfig.Endpoint),
			Region:           aws.String(s3EnvConfig.Region),
			DisableSSL:       aws.Bool(strings.HasSuffix(s3EnvConfig.Endpoint, "http://")),
			S3ForcePathStyle: aws.Bool(true),
		}
		awsSession := session.New(s3Config)
		DefaultClient = s3.New(awsSession)
	}

}

func ListBlobs(ctx context.Context, q *blob.ListBlobs) error {
	tenant := blob.TenantFromContext(ctx)
	basePath := fmt.Sprintf("tenants/%d/", tenant)

	response, err := DefaultClient.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket:  aws.String(cfg.S3.BucketName),
		MaxKeys: aws.Int64(3000),
		Prefix:  aws.String(basePath),
	})
	if err != nil {
		return wrap(err, "failed to list blobs from S3")
	}

	files := make([]string, 0)
	for _, item := range response.Contents {
		key := *item.Key
		files = append(files, key[len(basePath):])
	}

	sort.Strings(files)
	q.Result = files
	return nil
}

func GetBlobByKey(ctx context.Context, q *blob.GetBlobByKey) error {
	resp, err := DefaultClient.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3.BucketName),
		Key:    aws.String(keyFullPathURL(ctx, q.Key)),
	})
	if err != nil {
		if isNotFound(err) {
			return blob.ErrNotFound
		}
		return wrap(err, "failed to get blob '%s' from S3", q.Key)
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wrap(err, "failed to read blob body '%s' from S3", q.Key)
	}

	q.Result = &blob.Blob{
		Content:     bytes,
		ContentType: *resp.ContentType,
		Size:        *resp.ContentLength,
	}
	return nil
}

func StoreBlob(ctx context.Context, c *blob.StoreBlob) error {
	if err := blob.ValidateKey(c.Key); err != nil {
		return wrap(err, "failed to validate blob key '%s'", c.Key)
	}

	reader := bytes.NewReader(c.Content)
	_, err := DefaultClient.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3.BucketName),
		Key:         aws.String(keyFullPathURL(ctx, c.Key)),
		ContentType: aws.String(c.ContentType),
		ACL:         aws.String(s3.ObjectCannedACLPrivate),
		Body:        reader,
	})
	if err != nil {
		return wrap(err, "failed to upload blob '%s' to S3", c.Key)
	}
	return nil
}

func DeleteBlob(ctx context.Context, c *blob.DeleteBlob) error {
	_, err := DefaultClient.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3.BucketName),
		Key:    aws.String(keyFullPathURL(ctx, c.Key)),
	})
	if err != nil && !isNotFound(err) {
		return wrap(err, "failed to delete blob '%s' from S3", c.Key)
	}
	return nil
}

func keyFullPathURL(ctx context.Context, key string) string {
	tenant := blob.TenantFromContext(ctx)
	if tenant > 0 {
		return path.Join("tenants", strconv.FormatInt(tenant, 10), key)
	}
	return key
}

func isNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Code() == s3.ErrCodeNoSuchKey
	}
	return false
}

func wrap(err error, format string, a ...interface{}) error {
	if awsErr, ok := err.(awserr.Error); ok {
		return errors.Wrapf(awsErr.OrigErr(), format, a...)
	}
	return errors.Wrapf(err, format, a...)
}
