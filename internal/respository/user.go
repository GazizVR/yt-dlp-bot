package respository

import (
	"context"
	"fmt"

	"github.com/gazizvr/yt-dlp-bot/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb redis.Client
}

func key(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (s *UserStore) Get(ctx context.Context, id int64) (*domain.User, error) {
	s.rdb.HGetAll(ctx, key(id))
	return nil, nil
}
func (s *UserStore) Set(ctx context.Context, user domain.User) error {
	return nil
}
func (s *UserStore) Delete(ctx context.Context, id int64) error {
	return nil
}
