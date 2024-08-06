//go:build wireinject
// +build wireinject

package mysql

import (
	"context"

	"github.com/google/wire"
)

func InitMySQL(ctx context.Context) (*MySQL, error) {
	wire.Build(
		NewMasterDB,
		NewSlaveDB,
		wire.Struct(new(MySQL), "Master", "Slave"),
	)
	return nil, nil
}
