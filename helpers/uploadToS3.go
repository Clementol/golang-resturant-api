package helpers

import (
	"bytes"
	"context"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func UploadToS3(fileName string, file multipart.File, fileHeader *multipart.FileHeader) (msg string, url string) {
	godotenv.Load()
	BUCKET := os.Getenv("BUCKET")
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	}))

	svc := s3.New(sess)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	size := fileHeader.Size
	buffer := make([]byte, size)

	file.Read(buffer)

	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(BUCKET),
		ACL:                  aws.String("public-read"),
		Key:                  aws.String("images/" + fileName),
		Body:                 bytes.NewReader([]byte(buffer)),
		ServerSideEncryption: aws.String("AES256"),
	})

	params := &s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String("images/" + fileName),
	}

	if err != nil {
		msg = "unable to upload file" + err.Error()
		return msg, url
	}
	req, _ := svc.GetObjectRequest(params)
	rest.Build(req)
	url = req.HTTPRequest.URL.String()
	return msg, url
}
