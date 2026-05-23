package message

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gliedabrennung/sedna/internal/entity"
	"github.com/redis/go-redis/v9"
)

func TestRedisCache_Integration(t *testing.T) {
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if redisAddr == "" {
		t.Skip("TEST_REDIS_ADDR not set, skipping redis integration test")
	}

	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("could not connect to redis: %v", err)
	}
	defer client.Close()

	cache := NewRedisCache(client)
	chatID := "test:redis:1"

	client.Del(ctx, "chat:"+chatID+":cache")

	msg := &entity.Message{
		ChatID:    chatID,
		MessageID: "msg-1",
		FromID:    1,
		ToID:      2,
		Content:   "hello cache",
		CreatedAt: time.Now(),
	}

	cache.CacheMessage(ctx, msg)

	cached, ok := cache.GetCachedHistory(ctx, chatID, 10)
	if !ok {
		t.Fatal("expected to find cached history")
	}
	if len(cached) != 1 {
		t.Fatalf("expected 1 message, got %d", len(cached))
	}
	if cached[0].MessageID != "msg-1" {
		t.Errorf("expected msg-1, got %s", cached[0].MessageID)
	}

	msgs := []*entity.Message{
		{ChatID: chatID, MessageID: "msg-2", Content: "second"},
		{ChatID: chatID, MessageID: "msg-3", Content: "third"},
	}
	cache.WarmUpCache(ctx, chatID, msgs)

	cached2, ok := cache.GetCachedHistory(ctx, chatID, 10)
	if !ok || len(cached2) != 2 {
		t.Fatalf("expected 2 messages after warmup, got %d", len(cached2))
	}
	if cached2[0].MessageID != "msg-2" {
		t.Errorf("expected msg-2 to be newest, got %s", cached2[0].MessageID)
	}
}

func TestRedisCache_PubSub(t *testing.T) {
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if redisAddr == "" {
		t.Skip("TEST_REDIS_ADDR not set, skipping redis pubsub test")
	}

	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()
	cache := NewRedisCache(client)
	chatID := "test:pubsub:1"

	ch, unsubscribe, err := cache.Subscribe(ctx, chatID)
	if err != nil {
		t.Fatalf("subscribe error: %v", err)
	}
	defer unsubscribe()

	msg := &entity.Message{
		ChatID:    chatID,
		MessageID: "pubsub-1",
		Content:   "realtime",
	}

	time.Sleep(100 * time.Millisecond)

	cache.Publish(ctx, msg)

	select {
	case got := <-ch:
		if got.MessageID != "pubsub-1" {
			t.Errorf("expected pubsub-1, got %s", got.MessageID)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for pubsub message")
	}
}
