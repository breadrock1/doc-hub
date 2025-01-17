package cloud

import (
	"bytes"
	"context"
	"time"
)

type DocumentHub struct {
	Cloud ICloud
}

type ICloud interface {
	IBucket
	IDocument
	IShare
	IExpired
}

type IBucket interface {
	GetBuckets(ctx context.Context) ([]string, error)
	CreateBucket(ctx context.Context, bucket string) error
	RemoveBucket(ctx context.Context, bucket string) error
	IsBucketExist(ctx context.Context, bucket string) (bool, error)
}

type IDocument interface {
	GetFiles(ctx context.Context, bucket, filePath string) ([]*StorageItem, error)
	CopyFile(ctx context.Context, bucket, srcPath, dstPath string) error
	MoveFile(ctx context.Context, bucket, srcPath, dstPath string) error
	RemoveFile(ctx context.Context, bucket, filePath string) error
	UploadFile(ctx context.Context, bucket, filePath string, data bytes.Buffer) error
	DownloadFile(ctx context.Context, bucket, filePath string) (bytes.Buffer, error)
}

type IShare interface {
	GetShareURL(ctx context.Context, bucket, filePath string, expired time.Duration) (string, error)
}

type IExpired interface {
	UploadExpired(ctx context.Context, bucket, filePath string, expired time.Time, data bytes.Buffer) error
}
