//go:build wireinject
// +build wireinject

package redis

import (
	"context"

	"github.com/google/wire"
)

func InitRedis(ctx context.Context) (*Redis, error) {
	wire.Build(NewRedis)
	return nil, nil
}
