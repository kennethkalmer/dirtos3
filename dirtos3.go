package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

func main() {
	AWSAuth := aws.Auth{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	region := aws.USWest2

	connection := s3.New(AWSAuth, region)

	bucket := connection.Bucket(os.Getenv("S3_BUCKET_NAME"))
	s3prefix := os.Getenv("S3_PREFIX")

	source := flag.String("source", ".", "Directory to upload")
	flag.Parse()

	files := dir(*source)

	fmt.Printf("Attempting upload of %s\n\n", files)

	for _, path := range files {
		upload(connection, bucket, s3prefix, path)
	}
}

func upload(connection *s3.S3, bucket *s3.Bucket, prefix string, path string) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	bytes := make([]byte, fileSize)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	filetype := http.DetectContentType(bytes)

	s3path := fmt.Sprintf("%s/%s", prefix, path)

	const fileChunk = 5242880

	// Multipart upload for objects larger than 5 MB
	if fileSize > fileChunk {
		multi, err := bucket.InitMulti(s3path, filetype, s3.ACL("public-read"))

		if err != nil {
			fmt.Println(err)
			return
		}

		parts, err := multi.PutAll(file, fileChunk)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = multi.Complete(parts)

		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		err = bucket.Put(s3path, bytes, filetype, s3.ACL("public-read"))

		if err != nil {
			fmt.Println(err)
			return
		}

	}

	fmt.Printf("Uploaded %s (%v bytes)\n", path, fileSize)
}

func dir(thepath string) []string {
	var files []string

	filepath.Walk(thepath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here
			return nil       // but continue walking
		}

		if !!fi.IsDir() {
			return nil // ignore directories
		}

		files = append(files, path)
		return nil
	})

	return files
}
