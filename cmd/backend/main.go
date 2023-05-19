package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gari8/video2gif/cmd/backend/docs"
	"github.com/gari8/video2gif/config"
	"github.com/gari8/video2gif/infrastructure/driver"
	"github.com/gari8/video2gif/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Dependency Injection
	cfg := config.Load()
	rdb := driver.NewRedisClient(cfg.Redis.Url)
	awsConfig := driver.NewAwsConfig()
	sess := session.Must(session.NewSession(awsConfig))
	blobStore := driver.NewBlobStore(sess, driver.BucketName)
	fmp := driver.NewFfmpeg(awsConfig)

	repo := persistence.NewVideoPersistence(blobStore, fmp, rdb)
	ctrl := newController(repo)

	// Setup webserver
	app := gin.Default()

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK!"})
	})

	api := app.Group("/api/v1")

	api.POST("/videos", ctrl.UploadVideo)

	docs.SwaggerInfo.Title = "video2gif"
	docs.SwaggerInfo.Description = "video2gif"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.App.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(app.Run(fmt.Sprintf(":%d", cfg.App.Port)))
}

type videoRepo interface {
	UpsertVideo(video io.Reader, objectPath string) error
}

type controller struct {
	videoRepo
}

func newController(videoRepo videoRepo) *controller {
	return &controller{videoRepo}
}

// UploadVideo godoc
// @Summary 動画送信API
// @Tags    Video
// @Accept  mpfd
// @Produce json
// @Param   file formData file true "動画"
// @Success 200  "OK"
// @Failure 500
// @Router  /videos [post]
func (c controller) UploadVideo(ctx *gin.Context) {
	mfh, err := ctx.FormFile("file")
	f, err := mfh.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	if err := c.videoRepo.UpsertVideo(f, mfh.Filename); err != nil {
		log.Fatalln(err)
	}
}
