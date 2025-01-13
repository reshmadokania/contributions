package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func resumableMultipartUpload(chunkSize int64, useAccelerate bool) error {

	bucketName := "intusurgtest"
	filePath := "/Users/rdokania/repository/s3uploadtest/uploadfilelarge20mb.pdf"
	key := "uploadfilelarge20mb.pdf"
	// Configure the session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"<accesskey>", // Access Key ID
			"<secretkey>", // Secret Access Key
			"",            // Token (optional, for temporary credentials)
		),
		// Use acceleration endpoint if enabled
		//Endpoint: aws.String("https://" + bucketName + ".s3-accelerate.amazonaws.com"),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	// Create S3 service client
	svc := s3.New(sess)
	//uploadID := "yRpPFX9NiVaO9QGLTzdcLdStIKRW9fZzPBwnmBFqVdiqcxEm7AMnvfBvlkhQvPCxSOmBFDFSsrXe_GJ.frPxMtSepVDyl2yb6CI2vMBZ.XSOzNKtlZkfTcUOzrECCBSh"

	uploadID := nil //save this to db once this is generated
	// todo save uploadId in db. will be used
	//for future if this needs to be resumed again

	// Check if uploadId is nil and initialize if needed
	log.Printf("Multipart Upload ID: %s\n", uploadID)
	if uploadID == "" {
		createOutput, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("failed to initialize multipart upload: %w", err)
		}
		uploadID = *createOutput.UploadId
		log.Printf("Initialized new multipart upload with UploadId: %s", uploadID)
	}

	var completedParts []*s3.CompletedPart
	partNumber := int64(1)
	// Check for previously uploaded parts (for resumability)
	listPartsResp, _ := svc.ListParts(&s3.ListPartsInput{
		Bucket:   aws.String(bucketName),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})
	uploadedParts := make(map[int64]string)
	for _, part := range listPartsResp.Parts {
		uploadedParts[*part.PartNumber] = *part.ETag
		log.Printf("Part %d already uploaded with ETag: %s", part.PartNumber, *part.ETag)
	}
	// Read the file in chunks and upload each part
	buffer := make([]byte, chunkSize)
	for {
		// Skip already uploaded parts
		if _, exists := uploadedParts[partNumber]; exists {
			log.Printf("Part %d already uploaded. Skipping...\n", partNumber)
			/*completedParts = append(completedParts, &s3.CompletedPart{
				ETag:        aws.String(uploadedParts[partNumber]),
				PartNumber: aws.Int64(partNumber),
			})*/
			partNumber++
			continue
		}

		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			// Abort the upload on error
			_, _ = svc.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucketName),
				Key:      aws.String(key),
				UploadId: aws.String(uploadID),
			})
			return fmt.Errorf("error reading file: %w", err)
		}
		if bytesRead == 0 {
			break
		}

		// Upload part
		partResp, err := svc.UploadPart(&s3.UploadPartInput{
			Bucket:     aws.String(bucketName),
			Key:        aws.String(key),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int64(partNumber),
			Body:       bytes.NewReader(buffer[:bytesRead]),
		})
		if err != nil {
			// Abort the upload on error
			_, _ = svc.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucketName),
				Key:      aws.String(key),
				UploadId: aws.String(uploadID),
			})
			return fmt.Errorf("error uploading part %d: %w", partNumber, err)
		}
		//log.Printf("Uploaded part %d. Upload ID: %s\n", partNumber, uploadID)
		completedParts = append(completedParts, &s3.CompletedPart{
			ETag:       partResp.ETag,
			PartNumber: aws.Int64(partNumber),
		})
		log.Printf("Uploaded part %d with ETag: %s", partNumber, *partResp.ETag)
		partNumber++
	}
	// Complete the multipart upload
	_, err = svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucketName),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}
	fmt.Printf("Successfully uploaded %s to bucket %s\n", key, bucketName)
	return nil
}
func main() {
	err := resumableMultipartUpload(5*1024*1024, true) // 10MB chunk size
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
