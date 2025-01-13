# S3 File Upload and Integrity Verification

This repository provides three Go programs demonstrating different ways to interact with Amazon S3 for file uploads, resumable multipart uploads, and integrity verification of uploaded/downloaded files. These examples utilize the AWS SDK for Go.

## Prerequisites

Before running the examples, ensure the following:

- You have AWS credentials (Access Key ID and Secret Access Key) with permissions to access S3.
- You have an S3 bucket where files can be uploaded.
- Go is installed on your machine.
- Replace placeholder values like `<access key>` and `<secret key>` with your actual AWS credentials.

## Files Overview

### 1. **`basic_upload.go`**

This program demonstrates a simple file upload to an S3 bucket.

#### Key Features:
- Uploads a file to a specified S3 bucket.
- Uses AWS SDK v2 for Go.
- Loads AWS configuration from either hardcoded credentials or environment variables.

#### Usage:
1. Update the constants `region`, `accessKeyID`, `secretAccessKey`, `filePath`, and `bucketName` with your configuration.
2. Run the program:
   ```bash
   go run basic_upload.go
   ```
3. Verify the file has been uploaded to your S3 bucket.

### 2. **`resumable_multipart_upload.go`**

This program implements a resumable multipart upload for large files.

#### Key Features:
- Splits large files into chunks and uploads them as parts.
- Supports resuming uploads by saving the `uploadId` and already uploaded parts.
- Uses AWS SDK v1 for Go.

#### Usage:
1. Update the variables `bucketName`, `filePath`, and `key` with your configuration.
2. Run the program:
   ```bash
   go run resumable_multipart_upload.go
   ```
3. If interrupted, the program will check previously uploaded parts and resume the upload.

### 3. **`integrity_verification.go`**

This program verifies the integrity of a file by comparing the MD5 hash of the original file with the hash of the downloaded file.

#### Key Features:
- Downloads a file from S3 (method is provided but commented out for customization).
- Calculates and compares MD5 hashes of the original and downloaded files to ensure integrity.

#### Usage:
1. Update the variables `originalFilePath` and `downloadedFilePath` with paths to the original and downloaded files.
2. Uncomment and customize the `downloadFileFromS3` method if needed.
3. Run the program:
   ```bash
   go run integrity_verification.go
   ```
4. The program will display whether the files match based on their MD5 hashes.

## How to Run

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Install dependencies (if any):
   ```bash
   go mod tidy
   ```

3. Run the desired Go file:
   ```bash
   go run <filename>.go
   ```

## Notes

- Ensure your AWS credentials are correctly configured.
- The multipart upload program can save the `uploadId` in a database or file for better resumability.
- Modify and extend the programs as needed to suit your requirements.

## License

This repository is licensed under the MIT License. See the `LICENSE` file for details.

