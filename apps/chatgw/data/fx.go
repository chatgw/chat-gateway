package data

import (
	"github.com/airdb/chat-gateway/apps/chatgw/data/repos"
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			repos.NewUserRepo,
			repos.NewKeyRepo,
		),
	)
}
