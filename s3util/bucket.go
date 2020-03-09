package s3util

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// bucket represents an S3 bucket and handles read, write and delete operations.
type Bucket struct {
	name   string
	client *s3.S3

	// The prefix should end with "/", so that the resulting bucket operates in a subfolder.
	prefix        string
	useLegacyList bool
}

func NewBucket(ctx context.Context, c *Config) (*Bucket, error) {
	sess, err := SetupS3(c)
	if err != nil {
		return nil, err
	}

	return openBucket(ctx, sess, c.BucketName, c.Prefix, false)
}

// openBucket returns an S3 Bucket.
func openBucket(ctx context.Context, sess *session.Session, bucketName string, prefix string, useLegacyList bool) (*Bucket, error) {
	if sess == nil {
		return nil, errors.New("s3blob.OpenBucket: sess is required")
	}
	if bucketName == "" {
		return nil, errors.New("s3blob.OpenBucket: bucketName is required")
	}
	if prefix == "" {
		prefix = "tmp/"
	}

	return &Bucket{
		name:          bucketName,
		client:        s3.New(sess),
		prefix:        prefix,
		useLegacyList: useLegacyList,
	}, nil
}
