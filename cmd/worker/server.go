package main

import (
	"context"
	"github.com/gari8/video2gif/infrastructure/driver"
	"log"
	"os"
)

type PubSubClient interface {
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

// redis.Clientを隠蔽してサーバー起動するところ
type server struct {
	PubSubClient
}

func newServer(rdb PubSubClient) *server {
	return &server{rdb}
}

func (s server) Run(f func(ctx context.Context, payload string) error) {
	ctx := context.Background()
	pubSub := s.PubSubClient.Subscribe(ctx, driver.JobQueueKey)
	defer pubSub.Close()
	ch := pubSub.Channel()
	for msg := range ch {
		log.Println("...received")
		err := f(ctx, msg.Payload)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("success")
		}
	}
	os.Exit(0)
}
