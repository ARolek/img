package uploader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arolek/tools"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/arolek/img/config"
)

const region = "us-east-1"

//	move the file to s3
func MoveToS3(filePath, fileKey string) (err error) {
	//	creds := aws.Creds(*config.AWS_ACCESS_KEY_ID, *config.AWS_SECRET_ACCESS_KEY, "")

	svc := s3.New(
		session.New(),
		&aws.Config{Region: aws.String(region)},
	)

	//	cli := s3.New(creds, "us-east-1", nil)

	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	stat, err := file.Stat()
	if err != nil {
		return
	}

	putReq := s3.PutObjectInput{
		ACL:           aws.String(s3.BucketCannedACLPublicRead),
		Body:          file,
		Bucket:        aws.String(*config.S3_BUCKET),
		Key:           aws.String(fileKey),
		ContentLength: aws.Int64(stat.Size()),
	}

	resp, err := svc.PutObject(&putReq)
	if err != nil {
		return
	}

	//	TODO: manage errors
	//	log.Println("s3 put response: ", resp)

	return
}

//	download from s3 based on bucket and file key.
//	return an error or the tmpFile location
func FetchFromS3(key, bucket string) (tmpFilePath string, err error) {
	//	build our s3 url
	url := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, key)

	//	our temp file path
	tmpFilePath = filepath.Join(*config.FS_TEMP, RandomHash())

	//	downlod from s3 to our temp file location
	if err = tools.HTTPdownload(url, tmpFilePath); err != nil {
		return
	}

	return
}
