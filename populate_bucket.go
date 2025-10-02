package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Load default configuration (reads from environment variables, ~/.aws/config, etc.)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Open the file you want to upload
	file, err := os.Open("file.html")
	if err != nil {
		log.Fatalf("unable to open file %q, %v", "file.html", err)
	}
	defer file.Close()

	// Create an uploader and upload the file
	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("PUT YOUR BUCKET IDENTIFIER NAME HERE "),
		Key:    aws.String("file.html"),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("unable to upload %q to %q, %v", "file.html", "PUT YOUR BUCKET IDENTIFIER NAME HERE ", err)
	}

	log.Printf("successfully uploaded file to %s\n", result.Location)
}
