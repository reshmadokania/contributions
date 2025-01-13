package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	region          = "us-east-1"    // Change as per your region
	accessKeyID     = "<access key>" // Change as per your AWS Access Key
	secretAccessKey = "<secret key"  // Change as per your AWS Secret Key
)

func main() {
	// File to upload
	filePath := "/Users/rdokania/repository/s3uploadtest/test.txt" // Change as per your file
	bucketName := "intusurgtest"                                   // Change as per your bucket name
	objectKey := "test.txt"                                        // S3 object key (same or different)

	// Load AWS Configuration from environment or hardcoded in credentials
	cfg, err := loadAWSConfig()
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Upload file to S3
	err = uploadFileToS3(client, bucketName, objectKey, filePath)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	log.Println("File uploaded successfully.")
}

// loadAWSConfig loads AWS configuration from environment variables or from hardcoded credentials.
func loadAWSConfig() (aws.Config, error) {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithEndpointResolver(customResolver),
	)
}

// uploadFileToS3 uploads a file to the specified S3 bucket.
func uploadFileToS3(client *s3.Client, bucket, key, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %v", filePath, err)
	}
	defer file.Close()

	// Create a PutObject request to upload the file
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
		//ACL:    types.ObjectCannedACLPublicRead, // Optional: Set object ACL if required
	}

	// Upload file
	_, err = client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}
