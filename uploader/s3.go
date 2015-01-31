package uploader

import (
	//	"io"
	"log"
	"os"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/s3"

	"img/config"
)

//	move the file to s3
func MoveToS3(filePath, fileKey string) (err error) {

	creds := aws.Creds(*config.AWS_ACCESS_KEY_ID, *config.AWS_SECRET_ACCESS_KEY, "")
	cli := s3.New(creds, "us-east-1", nil)

	log.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	stat, err := file.Stat()
	if err != nil {
		return
	}

	putReq := s3.PutObjectRequest{
		ACL:           aws.String(s3.BucketCannedACLPublicRead),
		Body:          file,
		Bucket:        aws.String(*config.S3_BUCKET),
		Key:           aws.String(fileKey),
		ContentLength: aws.Long(stat.Size()),
	}

	resp, err := cli.PutObject(&putReq)
	if err != nil {
		return
	}

	return
}
