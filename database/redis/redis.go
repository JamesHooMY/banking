package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisCluster struct {
	Client *redis.ClusterClient
}

func NewRedisCluster(ctx context.Context) (*RedisCluster, error) {
	addrs := viper.GetStringSlice("redis.cluster.addrs")
	password := viper.GetString("redis.cluster.password")

	client, err := initRedisClusterClient(ctx, addrs, password)
	if err != nil {
		return nil, err
	}

	return &RedisCluster{Client: client}, nil
}

func initRedisClusterClient(ctx context.Context, addrs []string, password string) (*redis.ClusterClient, error) {
	var client *redis.ClusterClient

	err := retry(ctx, func() error {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: password,
		})

		// Ping the Redis server to check the connection
		return client.Ping(ctx).Err()
	}, 5, 2*time.Second)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func retry(ctx context.Context, action func() error, attempts int, sleep time.Duration) error {
	for i := 0; i < attempts; i++ {
		err := action()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(sleep)
		}
	}
	return nil
}
