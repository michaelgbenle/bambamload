package uploadService

import (
	"bambamload/logger"
	"bambamload/utils"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UploadService struct {
	ApplicationKeyID string
	ApplicationKey   string
	BucketName       string
	BucketRegion     string
}

func NewUploadService() *UploadService {
	return &UploadService{
		ApplicationKeyID: os.Getenv("B2_APPLICATION_KEY_ID"),
		ApplicationKey:   os.Getenv("B2_APPLICATION_KEY"),
		BucketName:       os.Getenv("B2_BUCKET_NAME"),
		BucketRegion:     os.Getenv("B2_BUCKET_REGION"),
	}
}

func (us *UploadService) GeneratePresignedDownloadURL(ctx context.Context, client *s3.Client, bucket, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry)) // e.g. 1 * time.Hour
	if err != nil {
		logger.Logger.Errorf("Failed to generate presigned download URL for bucket %s, %s", bucket, err)
		return "", err
	}
	return req.URL, nil
}

func (us *UploadService) Upload(file multipart.File, fileName string) (string, error) {
	defer file.Close()
	ctx := context.Background()
	// Custom resolver for B2 endpoint
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               fmt.Sprintf("https://%s", us.BucketRegion),
			SigningRegion:     "eu-central-003", // usually matches the pod region
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			us.ApplicationKeyID,
			us.ApplicationKey,
			"", // session token – not used for B2
		)),
		config.WithEndpointResolverWithOptions(resolver),
		config.WithRegion("eu-central-003"), // required even if fake
	)
	if err != nil {
		logger.Logger.Errorf("[Upload]cannot load config: %v", err)
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 100 * 1024 * 1024 // 100 MiB parts (B2 min 5 MiB, max 5 GiB)
	})

	name, ext := utils.SplitFileName(fileName)
	extType := utils.ExtensionToContentType[ext]
	if extType == "" {
		logger.Logger.Errorf("[Upload]cannot parse file content type: %s", fileName)
		return "", fmt.Errorf("cannot parse file content type: %s", fileName)
	}

	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(us.BucketName),
		Key:         aws.String(fmt.Sprintf("%s/%s", name, time.Now().Format("20060102150405"))),
		Body:        file,
		ContentType: aws.String(extType),
	})
	if err != nil {
		//log.Fatalf("multipart upload failed: %v", err)
		logger.Logger.Errorf("[Upload]cannot upload object: %v", err)
		return "", err
	}

	key := ""
	if result.Key != nil {
		key = *result.Key
	}

	downloadURL, err := us.GeneratePresignedDownloadURL(ctx, client, us.BucketName, key, 24*time.Hour*6)
	if err != nil {
		logger.Logger.Errorf("[Upload]cannot generate download URL for bucket %s, %s", us.BucketName, err)
		return "", err
	}

	//fmt.Printf("Upload successful!\nLocation: %s\n", result.Location)
	//fmt.Printf("Download URL: %s\n", downloadURL)
	return downloadURL, nil
}
