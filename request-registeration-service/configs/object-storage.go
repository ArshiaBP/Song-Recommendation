package configs

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func UploadFile(fileContent *bytes.Reader, fileName string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("App .env file not found")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return err
	}

	// Define AWS credentials and bucket information
	awsAccessKeyID := os.Getenv("ACCESS_KEY")
	awsSecretAccessKey := os.Getenv("SECRET_KEY")
	endpoint := os.Getenv("ENDPOINT")
	bucketName := os.Getenv("BUCKET_NAME")

	// Initialize S3 client with custom configuration
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     awsAccessKeyID,
			SecretAccessKey: awsSecretAccessKey,
		}, nil
	})

	cfg.BaseEndpoint = aws.String(endpoint)

	client := s3.NewFromConfig(cfg)

	// Specify the destination key in the bucket
	destinationKey := "uploads/" + fileName

	// Use the S3 client to upload the file
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destinationKey),
		Body:   fileContent,
	})

	return err
}
