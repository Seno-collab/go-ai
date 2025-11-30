package storage

import (
	"context"
	"fmt"
	"go-ai/internal/config"
	"go-ai/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient() *MinioClient {
	logger := logger.NewLogger().With().Str("component", "minio").Logger()
	config, _ := config.LoadConfig()
	endpoint := fmt.Sprintf("%s:%s", config.MinioEndPoint, config.MinioPort)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect MinIO")
	}
	ctx := context.Background()
	exists, _ := client.BucketExists(ctx, config.Bucket)
	if !exists {
		if err := client.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{}); err != nil {
			logger.Fatal().Err(err).Msg("Cannot create bucket")
		}
		logger.Info().Str("bucket", config.Bucket).Msg("Bucket created")
	}
	return &MinioClient{
		Client: client,
		Bucket: config.Bucket,
	}
}

func (m *MinioClient) PublicUrl(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s",
		m.Client.EndpointURL().Host,
		m.Bucket,
		objectName,
	)
}
