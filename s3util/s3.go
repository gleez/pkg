package s3util

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string

	// The prefix should end with "/", so that the resulting bucket operates in a subfolder.
	Prefix string
	Debug  bool
}

func SetupS3(c *Config) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Credentials:                   credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
		Endpoint:                      aws.String(c.Endpoint),
		Region:                        aws.String(c.Region),
		DisableSSL:                    aws.Bool(strings.HasSuffix(c.Endpoint, "http://")),
		S3ForcePathStyle:              aws.Bool(true),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})
}

func IsNotFound(err error) bool {
	if e, ok := err.(awserr.Error); ok {
		return e.Code() == s3.ErrCodeNoSuchKey || e.Code() == "NoSuchKey" || e.Code() == "NotFound" || e.Code() == s3.ErrCodeObjectNotInActiveTierError
	}
	return false
}
