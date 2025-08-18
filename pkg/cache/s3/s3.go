package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Cache implements the BytesCache interface using Amazon S3 or compatible storage
type S3Cache struct {
	client     *s3.Client
	bucketName string
}

// NewS3Cache creates a new S3 cache with the specified bucket, endpoint, region, and optional access keys
func NewS3Cache(bucketName, endpoint, region, accessKey, secretKey string) (*S3Cache, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error when loading the default S3 configuration: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		if region != "" {
			o.Region = region
		}
		if accessKey != "" || secretKey != "" {
			o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
		}
	})

	return &S3Cache{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Push stores the data in the S3 bucket
func (sc *S3Cache) Push(key string, value []byte) error {
	_, err := sc.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(sc.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	})
	return err
}

// isNotFoundError checks if the error is a "not found" error from S3
func isNotFoundError(err error) bool {
	// Check for common S3 error messages indicating a missing object
	return err != nil && (
	// Look for common error message patterns in S3 errors
	strings.Contains(err.Error(), "NoSuchKey") ||
		strings.Contains(err.Error(), "NoSuchBucket") ||
		strings.Contains(err.Error(), "NoSuchUpload"))
}

// Pull retrieves the data from the S3 bucket
func (sc *S3Cache) Pull(key string) ([]byte, error) {
	result, err := sc.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(sc.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if the error is a "not found" error
		if isNotFoundError(err) {
			return nil, fmt.Errorf("not exist")
		}
		return nil, err
	}
	defer func() {
		_ = result.Body.Close()
	}()

	return io.ReadAll(result.Body)
}
