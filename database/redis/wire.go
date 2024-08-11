//go:build wireinject
// +build wireinject

package redis

import (
	"context"

	"github.com/google/wire"
)

func InitRedis(ctx context.Context) (*RedisCluster, error) {
	wire.Build(NewRedisCluster)
	return nil, nil
}
