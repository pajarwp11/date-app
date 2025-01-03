package users

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Users interface {
	SetViewedUser(ctx context.Context, key string, value []string) error
	DeleteRedisKey(ctx context.Context, key string) error
	GetViewedUser(ctx context.Context, key string) ([]string, error)
}
type usersRepository struct {
	RDB *redis.Client
}

func NewUsersRepository(rdb *redis.Client) Users {
	return &usersRepository{
		RDB: rdb,
	}
}

func (u *usersRepository) SetViewedUser(ctx context.Context, key string, value []string) error {
	serializedValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not serialize value: %v", err)
	}

	now := time.Now()

	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	durationUntilMidnight := midnight.Sub(now)

	err = u.RDB.Set(ctx, key, serializedValue, durationUntilMidnight).Err()
	if err != nil {
		return fmt.Errorf("could not set key %s: %v", key, err)
	}
	return nil
}

func (u *usersRepository) GetViewedUser(ctx context.Context, key string) ([]string, error) {
	serializedValue, err := u.RDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, fmt.Errorf("could not get key %s: %v", key, err)
	}

	var value []string
	err = json.Unmarshal(serializedValue, &value)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize value: %v", err)
	}

	return value, nil
}

func (u *usersRepository) DeleteRedisKey(ctx context.Context, key string) error {
	_, err := u.RDB.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("could not delete key %s: %v", key, err)
	}
	return nil
}
