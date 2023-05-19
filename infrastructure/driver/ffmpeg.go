package driver

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"strings"
)

type Ffmpeg struct {
	*aws.Config
}

func NewFfmpeg(awsConfig *aws.Config) *Ffmpeg {
	return &Ffmpeg{awsConfig}
}

func (d Ffmpeg) VideoToGif(output *s3.GetObjectOutput, filepath string) error {
	f, err := os.OpenFile("tmp/tmp.mp4", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, output.Body); err != nil {
		return err
	}
	paths := strings.Split(filepath, ".")
	s3Path := fmt.Sprintf("s3://%s", fmt.Sprintf("%s.%s", paths[0], "gif"))
	return ffmpeg.
		Input("tmp/tmp.mp4", ffmpeg.KwArgs{"ss": "00:00:01"}).
		Output(s3Path, ffmpeg.KwArgs{
			"aws_config": d.Config,
			"s":          "320x240",
			"pix_fmt":    "rgb24",
			"t":          "3",
			"r":          "3",
			"format":     "gif",
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
}
