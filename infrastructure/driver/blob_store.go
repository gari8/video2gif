package driver

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

const BucketName = "my-gifs"
const ConvertedBucketName = "converted"

type BlobStore struct {
	*s3.S3
	BucketName string
}

func NewBlobStore(sess *session.Session, bucketName string) *BlobStore {
	sss := s3.New(sess)
	_, err := sss.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		_, e := sss.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if e != nil {
			log.Fatalln(e)
		}
	}
	return &BlobStore{sss, bucketName}
}
