package s3util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/sandeepone/go-escape"
)

// SignedURLOptions sets options for SignedURL.
type SignedURLOptions struct {
	// Expiry sets how long the returned URL is valid for. It is guaranteed to be > 0.
	Expiry time.Duration

	// Method is the HTTP method that can be used on the URL; one of "GET", "PUT",
	// or "DELETE". Drivers must implement all 3.
	Method string

	// ContentType specifies the Content-Type HTTP header the user agent is
	// permitted to use in the PUT request. It must match exactly. See
	// EnforceAbsentContentType for behavior when ContentType is the empty string.
	// If this field is not empty and the bucket cannot enforce the Content-Type
	// header, it must return an Unimplemented error.
	//
	// This field will not be set for any non-PUT requests.
	ContentType string

	// If EnforceAbsentContentType is true and ContentType is the empty string,
	// then PUTing to the signed URL must fail if the Content-Type header is
	// present or the implementation must return an error if it cannot enforce
	// this. If EnforceAbsentContentType is false and ContentType is the empty
	// string, implementations should validate the Content-Type header if possible.
	// If EnforceAbsentContentType is true and the bucket cannot enforce the
	// Content-Type header, it must return an Unimplemented error.
	//
	// This field will always be false for non-PUT requests.
	EnforceAbsentContentType bool

	// The canned ACL to apply to the object. For more information, see Canned ACL
	// (https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#CannedACL).
	ACL string

	// Can be used to specify caching behavior along the request/reply chain. For
	// more information, see http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9
	// (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9).
	CacheControl string
}

func (b *Bucket) SignedURL(ctx context.Context, key string, opts *SignedURLOptions) (string, error) {
	if b.prefix != "" {
		key = b.prefix + key
	}
	key = escapeKey(key)

	if opts.Expiry <= 0 {
		opts.Expiry = 60
	}

	switch opts.Method {
	case http.MethodGet:
		in := &s3.GetObjectInput{
			Bucket: aws.String(b.name),
			Key:    aws.String(key),
		}
		req, _ := b.client.GetObjectRequest(in)
		return req.Presign(opts.Expiry * time.Second)
	case http.MethodPut:
		in := &s3.PutObjectInput{
			Bucket:      aws.String(b.name),
			Key:         aws.String(key),
			ContentType: aws.String(opts.ContentType),
			ACL:         aws.String(opts.ACL),
		}
		req, _ := b.client.PutObjectRequest(in)
		return req.Presign(opts.Expiry * time.Second)
	case http.MethodDelete:
		in := &s3.DeleteObjectInput{
			Bucket: aws.String(b.name),
			Key:    aws.String(key),
		}
		req, _ := b.client.DeleteObjectRequest(in)
		return req.Presign(opts.Expiry * time.Second)
	default:
		return "", fmt.Errorf("unsupported Method %q", opts.Method)
	}
}

// escapeKey does all required escaping for UTF-8 strings to work with S3.
func escapeKey(key string) string {
	return escape.HexEscape(key, func(r []rune, i int) bool {
		c := r[i]
		switch {
		// S3 doesn't handle these characters (determined via experimentation).
		case c < 32:
			return true
		// For "../", escape the trailing slash.
		case i > 1 && c == '/' && r[i-1] == '.' && r[i-2] == '.':
			return true
		// For "//", escape the trailing slash. Otherwise, S3 drops it.
		case i > 0 && c == '/' && r[i-1] == '/':
			return true
		}
		return false
	})
}

// unescapeKey reverses escapeKey.
func unescapeKey(key string) string {
	return escape.HexUnescape(key)
}
