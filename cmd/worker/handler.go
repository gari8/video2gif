package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gari8/video2gif/infrastructure/driver"
	"io"
	"path/filepath"
)

type videoRepo interface {
	ConvertToGif(output *s3.GetObjectOutput, path string) error
	GetVideo(objectPath string) (*s3.GetObjectOutput, error)
	UpsertVideo(video io.Reader, objectPath string) error
}

type handler struct {
	videoRepo
}

func newHandler(vr videoRepo) *handler {
	return &handler{vr}
}

func (h handler) ConvertAndStoreData(convertedPathPrefix, path string) error {
	output, err := h.videoRepo.GetVideo(path)
	if err != nil {
		return err
	}

	if err := h.videoRepo.ConvertToGif(output, filepath.Join(driver.BucketName, convertedPathPrefix, path)); err != nil {
		return err
	}
	return nil
}
