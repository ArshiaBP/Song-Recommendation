package configs

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"io"
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

func DownloadFile(fileName string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("App .env file not found")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
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

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destinationKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, destinationKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", destinationKey, err)
	}
	_, err = file.Write(body)
	return err
}
