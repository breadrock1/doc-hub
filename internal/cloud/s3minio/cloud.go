package s3minio

import (
	"bytes"
	"context"
	"errors"
	"log"
	"time"

	"docs-hub/internal/cloud"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Minio struct {
	config *cloud.CloudConfig
	mc     *minio.Client
}

func New(config *cloud.CloudConfig) *cloud.DocumentHub {
	minioOpts := &minio.Options{
		Creds:  credentials.NewStaticV4(config.Username, config.Password, ""),
		Secure: config.EnableSSL,
	}

	client, err := minio.New(config.Address, minioOpts)
	if err != nil {
		log.Fatalln("failed to connect to minio cloud: ", err.Error())
	}

	s3Minio := &S3Minio{
		config: config,
		mc:     client,
	}

	return &cloud.DocumentHub{Cloud: s3Minio}
}

func (mw *S3Minio) GetBuckets(ctx context.Context) ([]string, error) {
	buckets, err := mw.mc.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketNames := make([]string, len(buckets))
	for index, bucketInfo := range buckets {
		bucketNames[index] = bucketInfo.Name
	}

	return bucketNames, nil
}

func (mw *S3Minio) CreateBucket(ctx context.Context, bucket string) error {
	opts := minio.MakeBucketOptions{}
	return mw.mc.MakeBucket(ctx, bucket, opts)
}

func (mw *S3Minio) RemoveBucket(ctx context.Context, bucket string) error {
	return mw.mc.RemoveBucket(ctx, bucket)
}

func (mw *S3Minio) GetFiles(ctx context.Context, bucket, filePath string) ([]*cloud.StorageItem, error) {
	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    filePath,
		Recursive: false,
	}

	if mw.mc.IsOffline() {
		return nil, errors.New("cloud is offline")
	}

	dirObjects := make([]*cloud.StorageItem, 0)
	for obj := range mw.mc.ListObjects(ctx, bucket, opts) {
		if obj.Err != nil {
			log.Println("failed to get object: ", obj.Err)
			continue
		}

		dirObjects = append(dirObjects, &cloud.StorageItem{
			FileName:      obj.Key,
			DirectoryName: filePath,
			IsDirectory:   len(obj.ETag) == 0,
		})
	}

	return dirObjects, nil
}

func (mw *S3Minio) RemoveFile(ctx context.Context, bucket, filePath string) error {
	opts := minio.RemoveObjectOptions{}
	return mw.mc.RemoveObject(ctx, bucket, filePath, opts)
}

func (mw *S3Minio) UploadFile(ctx context.Context, bucket, filePath string, data bytes.Buffer) error {
	opts := minio.PutObjectOptions{}
	dataLen := int64(data.Len())
	_, err := mw.mc.PutObject(ctx, bucket, filePath, &data, dataLen, opts)
	return err
}

func (mw *S3Minio) CopyFile(ctx context.Context, bucket, srcPath, dstPath string) error {
	srcOpts := minio.CopySrcOptions{Bucket: bucket, Object: srcPath}
	dstOpts := minio.CopyDestOptions{Bucket: bucket, Object: dstPath}
	_, err := mw.mc.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return err
	}

	return nil
}

func (mw *S3Minio) MoveFile(ctx context.Context, bucket, srcPath, dstPath string) error {
	err := mw.CopyFile(ctx, bucket, srcPath, dstPath)
	if err != nil {
		return err
	}

	return mw.RemoveFile(ctx, bucket, srcPath)
}

func (mw *S3Minio) DownloadFile(ctx context.Context, bucket, filePath string) (bytes.Buffer, error) {
	var objBody bytes.Buffer

	opts := minio.GetObjectOptions{}
	obj, err := mw.mc.GetObject(ctx, bucket, filePath, opts)
	if err != nil {
		return objBody, err
	}

	_, err = objBody.ReadFrom(obj)
	if err != nil {
		return objBody, err
	}

	return objBody, nil
}

func (mw *S3Minio) GetShareURL(ctx context.Context, bucket, filePath string, expired time.Duration) (string, error) {
	url, err := mw.mc.PresignedGetObject(ctx, bucket, filePath, expired, map[string][]string{})
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (mw *S3Minio) UploadExpired(ctx context.Context, bucket, filePath string, expired time.Time, data bytes.Buffer) error {
	opts := minio.PutObjectOptions{
		Expires: expired,
	}

	dataLen := int64(data.Len())
	_, err := mw.mc.PutObject(ctx, bucket, filePath, &data, dataLen, opts)
	return err

}
