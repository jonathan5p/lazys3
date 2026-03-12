package s3

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jonathan5p/lazys3/cmd/lazys3-tui/logger"
)

type Object struct {
	Key          string
	Size         int64
	LastModified *time.Time
	IsDir        bool
}

type Client struct {
	client *s3.Client
	log    *logger.Logger
}

func NewClient(ctx context.Context, log *logger.Logger) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &Client{client: client, log: log}, nil
}

func (c *Client) ListBuckets(ctx context.Context) ([]string, error) {
	c.log.Info("S3 ListBuckets called")
	result, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		c.log.Error("S3 ListBuckets failed: %v", err)
		return nil, err
	}

	buckets := make([]string, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		if b.Name != nil {
			buckets = append(buckets, *b.Name)
		}
	}
	c.log.Info("S3 ListBuckets returned %d buckets", len(buckets))
	return buckets, nil
}

func (o Object) Format() string {
	if o.IsDir {
		dir := strings.TrimSuffix(o.Key, "/")
		parts := strings.Split(dir, "/")
		return parts[len(parts)-1] + "/"
	}
	parts := strings.Split(o.Key, "/")
	return parts[len(parts)-1]
}

func (c *Client) ListObjects(ctx context.Context, bucket, prefix string) ([]Object, error) {
	c.log.Info("S3 ListObjects called bucket=%q prefix=%q", bucket, prefix)

	result, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		c.log.Error("S3 ListObjects failed: %v", err)
		return nil, err
	}

	var objects []Object

	for _, commonPrefix := range result.CommonPrefixes {
		if commonPrefix.Prefix != nil {
			objects = append(objects, Object{
				Key:   *commonPrefix.Prefix,
				IsDir: true,
			})
		}
	}

	for _, item := range result.Contents {
		if item.Key != nil {
			objects = append(objects, Object{
				Key:          *item.Key,
				Size:         *item.Size,
				LastModified: item.LastModified,
				IsDir:        false,
			})
		}
	}

	c.log.Info("S3 ListObjects returned %d items", len(objects))
	return objects, nil
}

func (c *Client) DownloadObject(ctx context.Context, bucket, key, destPath string) error {
	c.log.Info("S3 DownloadObject called bucket=%q key=%q dest=%q", bucket, key, destPath)
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		c.log.Error("S3 DownloadObject failed: %v", err)
		return err
	}
	defer result.Body.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		c.log.Error("DownloadObject create file failed: %v", err)
		return err
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, result.Body)
	if err != nil {
		c.log.Error("DownloadObject copy failed: %v", err)
		return err
	}
	c.log.Info("DownloadObject completed bytes=%d", n)
	return nil
}
