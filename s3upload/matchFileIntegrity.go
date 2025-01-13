package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	//"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
)

func main() {
	// Replace with your details
	//bucketName := "intusurgtest"
	//fileKey := "uploadfilelarge20mb.pdf"
	originalFilePath := "/Users/rdokania/repository/s3uploadtest/uploadfilelarge20mb.pdf"
	downloadedFilePath := "/Users/rdokania/repository/s3uploadtest/download/uploadfilelarge20mb.pdf"

	// Initialize AWS SDK config and S3 client
	/*cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		return
	}
	client := s3.NewFromConfig(cfg)

	// Download file from S3
	err = downloadFileFromS3(context.TODO(), client, bucketName, fileKey, downloadedFilePath)
	if err != nil {
		fmt.Printf("Failed to download file: %v\n", err)
		return
	}*/

	// Verify MD5 hashes
	isMatch, err := verifyFileIntegrity(originalFilePath, downloadedFilePath)
	if err != nil {
		fmt.Printf("Failed to verify file integrity: %v\n", err)
		return
	}

	if isMatch {
		fmt.Println("MD5 hashes match. File integrity verified.")
	} else {
		fmt.Println("MD5 hashes do not match. File may be corrupted.")
	}
}

func downloadFileFromS3(ctx context.Context, client *s3.Client, bucketName, key, outputPath string) error {
	// Open the output file for writing
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download the file from S3
	_, err = client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}, func(o *s3.Options) {
		o.HTTPClient = nil // Optional: Customize HTTP client if needed
	})
	if err != nil {
		return fmt.Errorf("failed to download file from S3: %w", err)
	}

	return nil
}

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5 hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyFileIntegrity(originalFilePath, downloadedFilePath string) (bool, error) {
	// Calculate MD5 of the original file
	originalMD5, err := calculateMD5(originalFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to calculate MD5 for original file: %w", err)
	}

	// Calculate MD5 of the downloaded file
	downloadedMD5, err := calculateMD5(downloadedFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to calculate MD5 for downloaded file: %w", err)
	}

	fmt.Printf("Original MD5:   %s\n", originalMD5)
	fmt.Printf("Downloaded MD5: %s\n", downloadedMD5)

	return originalMD5 == downloadedMD5, nil
}

