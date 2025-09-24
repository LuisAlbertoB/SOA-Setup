package services

import (
	"context"
	"log"
	"mime/multipart"
	"time"

	appConfig "rob-api-go/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client     *s3.Client
	Presign    *s3.PresignClient
	BucketName string
}

var s3Svc *S3Service

// NewS3Service inicializa y devuelve un singleton del servicio S3.
func NewS3Service() *S3Service {
	if s3Svc != nil {
		return s3Svc
	}

	region := appConfig.Config.AWSRegion
	bucketName := appConfig.Config.AWSS3BucketName

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("No se pudo cargar la configuración de AWS: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	s3Svc = &S3Service{
		Client:     client,
		Presign:    presignClient,
		BucketName: bucketName,
	}
	return s3Svc
}

// UploadFile sube un archivo a S3.
func (s *S3Service) UploadFile(fileHeader *multipart.FileHeader, key string) (*s3.PutObjectOutput, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
}

// DeleteFile elimina un archivo de S3.
func (s *S3Service) DeleteFile(key string) (*s3.DeleteObjectOutput, error) {
	return s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
}

// GetSignedURL genera una URL firmada y temporal para acceder a un archivo.
func (s *S3Service) GetSignedURL(key string) (string, error) {
	if key == "" {
		return "", nil
	}
	request, err := s.Presign.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * time.Minute) // La URL expira en 15 minutos
	})
	if err != nil {
		return "", err
	}
	return request.URL, nil
}
