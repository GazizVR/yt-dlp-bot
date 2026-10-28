package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RequestStore struct {
	rdb *redis.Client
}

func NewRequestStore(rdb *redis.Client) *RequestStore {
	return &RequestStore{
		rdb: rdb,
	}
}

func idKey(id int64) string {
	return fmt.Sprintf("request:%d", id)
}

func (s *RequestStore) Get(ctx context.Context, id int64) (string, error) {
	return s.rdb.Get(ctx, idKey(id)).Result()
}

const secKey string = "request:id:sec"

func (s *RequestStore) Create(ctx context.Context, mediaURL string) (int64, error) {
	createScript := redis.NewScript(`
		local mediaURL = ARGV[1]
		local key = KEYS[1]

		local rawId = redis.call("INCR",key)

		key = "request:" .. "" .. rawId

		redis.call("SET",key,mediaURL)

		return tonumber(rawId)	
	`)
	id, err := createScript.Run(ctx, s.rdb, []string{secKey}, mediaURL).Int64()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *RequestStore) Delete(ctx context.Context, id int64) error {
	return s.rdb.Del(ctx, idKey(id)).Err()
}
