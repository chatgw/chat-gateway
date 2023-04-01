package cachemod

import (
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(NewRedis),
	)
}
