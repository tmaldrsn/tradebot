package pubsub

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Subscribe(ctx context.Context, rdb *redis.Client, channel string, handler func(payload string, rdb *redis.Client)) {
	sub := rdb.Subscribe(ctx, channel)
	ch := sub.Channel()

	log.Printf("Subscribed to Redis channel: %s", channel)

	go func() {
		defer sub.Close() // ensure cleanup
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return // channel closed
				}
				log.Printf("📨 Received message on %s", channel)
				handler(msg.Payload, rdb)
			}
		}
	}()
}
