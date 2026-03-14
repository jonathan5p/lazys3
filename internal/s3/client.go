package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Client is the S3 access interface used throughout the application.
type Client interface {
	ListBuckets(ctx context.Context) ([]string, error)
	ListObjects(ctx context.Context, bucket, prefix string) ([]Object, error)
	HeadObject(ctx context.Context, bucket, key string) (map[string]string, error)
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	DownloadObject(ctx context.Context, bucket, key, dest string) error
}

// s3API is the subset of the AWS S3 client used by AWSClient.
type s3API interface {
	ListBuckets(ctx context.Context, params *awss3.ListBucketsInput, optFns ...func(*awss3.Options)) (*awss3.ListBucketsOutput, error)
	ListObjectsV2(ctx context.Context, params *awss3.ListObjectsV2Input, optFns ...func(*awss3.Options)) (*awss3.ListObjectsV2Output, error)
	HeadObject(ctx context.Context, params *awss3.HeadObjectInput, optFns ...func(*awss3.Options)) (*awss3.HeadObjectOutput, error)
	GetObject(ctx context.Context, params *awss3.GetObjectInput, optFns ...func(*awss3.Options)) (*awss3.GetObjectOutput, error)
}

// AWSClient wraps the AWS SDK S3 client.
type AWSClient struct {
	api s3API
}

// NewAWSClient creates a new AWSClient configured with the given parameters.
func NewAWSClient(ctx context.Context, region, endpoint, profile string) (*AWSClient, error) {
	opts := []func(*awsconfig.LoadOptions) error{}
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}
	if profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(profile))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	s3Opts := []func(*awss3.Options){}
	if endpoint != "" {
		s3Opts = append(s3Opts, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}
	return &AWSClient{api: awss3.NewFromConfig(cfg, s3Opts...)}, nil
}

// ListBuckets returns all bucket names accessible with the configured credentials.
func (c *AWSClient) ListBuckets(ctx context.Context) ([]string, error) {
	out, err := c.api.ListBuckets(ctx, &awss3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("list buckets: %w", err)
	}
	names := make([]string, 0, len(out.Buckets))
	for _, b := range out.Buckets {
		if b.Name != nil {
			names = append(names, *b.Name)
		}
	}
	return names, nil
}

// ListObjects lists objects under a given prefix, returning virtual directories
// as Object entries with IsDir=true.
func (c *AWSClient) ListObjects(ctx context.Context, bucket, prefix string) ([]Object, error) {
	input := &awss3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}
	out, err := c.api.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("list objects: %w", err)
	}
	var objects []Object
	for _, cp := range out.CommonPrefixes {
		if cp.Prefix != nil {
			objects = append(objects, Object{Key: *cp.Prefix, IsDir: true})
		}
	}
	for _, item := range out.Contents {
		obj := objectFromContent(item)
		objects = append(objects, obj)
	}
	return objects, nil
}

func objectFromContent(item types.Object) Object {
	obj := Object{}
	if item.Key != nil {
		obj.Key = *item.Key
	}
	if item.Size != nil {
		obj.Size = *item.Size
	}
	if item.LastModified != nil {
		obj.LastModified = *item.LastModified
	}
	return obj
}

// HeadObject returns the metadata for a specific S3 object.
func (c *AWSClient) HeadObject(ctx context.Context, bucket, key string) (map[string]string, error) {
	out, err := c.api.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("head object: %w", err)
	}
	meta := make(map[string]string, len(out.Metadata)+2)
	for k, v := range out.Metadata {
		meta[k] = v
	}
	if out.ContentType != nil {
		meta["ContentType"] = *out.ContentType
	}
	if out.ContentLength != nil {
		meta["ContentLength"] = fmt.Sprintf("%d", *out.ContentLength)
	}
	return meta, nil
}

// GetObject streams the content of an S3 object.
func (c *AWSClient) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	out, err := c.api.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return out.Body, nil
}

// DownloadObject downloads an S3 object to a local file path.
func (c *AWSClient) DownloadObject(ctx context.Context, bucket, key, dest string) error {
	rc, err := c.GetObject(ctx, bucket, key)
	if err != nil {
		return err
	}
	defer rc.Close()
	return writeToFile(rc, dest)
}
