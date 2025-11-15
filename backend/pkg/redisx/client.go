package redisx

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func New(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}

const (
	StreamDeclarations = "declarations:requests"
	StreamResults      = "declarations:results"
)

type StreamPublisher struct{ R *redis.Client }

func (p *StreamPublisher) PublishCreate(ctx context.Context, id int64, studentName, declType string) error {
	_, err := p.R.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamDeclarations,
		Values: map[string]any{
			"id":           id,
			"student_name": studentName,
			"decl_type":    declType,
		},
	}).Result()
	return err
}
