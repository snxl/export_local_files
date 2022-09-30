package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var (
	client *s3.S3
	wg     sync.WaitGroup
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACESS_KEY"),
			os.Getenv("AWS_SECRET_KEY"), ""),
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	client = s3.New(sess)
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}

	semaphore := make(chan struct{}, 1000)
	for {
		files, err := dir.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("Error: %v", err)
		}

		wg.Add(1)
		semaphore <- struct{}{}
		go upload(files[0], semaphore)
	}

	wg.Wait()
}

func upload(fileInfo os.FileInfo, semaphore chan struct{}) {
	defer wg.Done()

	fmt.Printf("Upload started: %s\n", fileInfo.Name())

	filepath := fmt.Sprintf("./tmp/%s", fileInfo.Name())

	file, err := os.Open(filepath)
	if err != nil {
		<-semaphore
		fmt.Printf("ERROR: %v", err)
		return
	}
	defer file.Close()

	var fileSize int64 = fileInfo.Size()

	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileInfo.Name()),
		Body:   bytes.NewReader(fileBuffer),
	})
	if err != nil {
		<-semaphore
		fmt.Printf("ERROR upload: %v", err)
		return
	}

	<-semaphore
	fmt.Printf("Upload finished: %s\n", fileInfo.Name())
}
