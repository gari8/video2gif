package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gari8/video2gif/config"
	"github.com/gari8/video2gif/infrastructure/driver"
	"github.com/gari8/video2gif/infrastructure/persistence"
	"golang.org/x/net/context"
)

func main() {
	cfg := config.Load()
	rdb := driver.NewRedisClient(cfg.Redis.Url)
	awsConfig := driver.NewAwsConfig()
	sess := session.Must(session.NewSession(awsConfig))
	blobStore := driver.NewBlobStore(sess, driver.BucketName)
	fmp := driver.NewFfmpeg(awsConfig)
	serv := newServer(rdb)
	vr := persistence.NewVideoPersistence(blobStore, fmp, rdb)
	hdl := newHandler(vr)
	serv.Run(func(ctx context.Context, payload string) error {
		return hdl.ConvertAndStoreData(driver.ConvertedBucketName, payload)
	})
}
