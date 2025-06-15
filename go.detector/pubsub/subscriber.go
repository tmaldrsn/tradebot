package pubsub

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Subscribe(ctx context.Context, rdb *redis.Client, channel string, handler func(string)) {
	sub := rdb.Subscribe(ctx, channel)
	ch := sub.Channel()

	log.Printf("Subscribed to Redis channel: %s", channel)

	go func() {
		for msg := range ch {
			log.Printf("📨 Received message on %s", channel)
			handler(msg.Payload)
		}
	}()
}
