package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type CloudflareService struct{}

// UploadCloudflare uploads the file to Cloudflare R2 and returns the file URL
func (service *CloudflareService) UploadCloudflare(ctx context.Context, file *multipart.FileHeader) (string, *apierror.ApiError) {
	minioClient, err := minio.New(config.R2Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.R2Config.AccessKeyID, config.R2Config.SecretAccessKey, ""),
		Secure: true,
	})	

	if err != nil {
		return "", apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to initialize R2 client: " + err.Error()).
			Build()
	}
	

    // Open the file
    fileContent, err := file.Open()
    if err != nil {
        return "", apierror.NewApiErrorBuilder().
            WithStatus(http.StatusBadRequest).
            WithMessage("Failed to open the file").
            Build()
    }
    defer fileContent.Close()

    // Upload the file to Cloudflare R2
	filePath := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), file.Filename)
	_, err = minioClient.PutObject(ctx, config.R2Config.BucketName, filePath, fileContent, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to upload file to Cloudflare R2").
			Build()
	}

	fileURL := fmt.Sprintf("https://pub-736ef3be77f045e8ba550ae958fe7e1b.r2.dev/%s", filePath)
	return fileURL, nil

}
