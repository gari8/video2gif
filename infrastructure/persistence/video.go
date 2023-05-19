package persistence

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gari8/video2gif/infrastructure/driver"
	"io"
)

type VideoPersistence struct {
	*driver.BlobStore
	*driver.Ffmpeg
	*redis.Client
}

func NewVideoPersistence(blobStore *driver.BlobStore, ffmpeg *driver.Ffmpeg, rdb *redis.Client) *VideoPersistence {
	return &VideoPersistence{blobStore, ffmpeg, rdb}
}

func (p VideoPersistence) ConvertToGif(output *s3.GetObjectOutput, path string) error {
	return p.Ffmpeg.VideoToGif(output, path)
}

func (p VideoPersistence) UpsertVideo(video io.Reader, objectPath string) error {
	_, err := p.BlobStore.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(p.BlobStore.BucketName),
		Key:    aws.String(objectPath),
		Body:   aws.ReadSeekCloser(video),
	})

	if err != nil {
		return err
	}

	return p.Client.Publish(context.Background(), driver.JobQueueKey, objectPath).Err()
}

func (p VideoPersistence) GetVideo(objectPath string) (*s3.GetObjectOutput, error) {
	return p.BlobStore.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.BlobStore.BucketName),
		Key:    aws.String(objectPath),
	})
}
