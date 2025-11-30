package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
)

func (m *MinioClient) UpLoadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("file-%d-%s", time.Now().UnixNano(), header.Filename)
	_, err := m.Client.PutObject(ctx,
		m.Bucket,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}
	return m.PublicUrl(objectName), nil
}

func (m *MinioClient) UploadLogo(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("logo-%d-%s", time.Now().UnixNano(), header.Filename)
	_, err := m.Client.PutObject(ctx,
		m.Bucket,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}
	return m.PublicUrl(objectName), nil
}
