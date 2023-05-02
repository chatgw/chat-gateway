package bootstrap

import (
	"github.com/airdb/chat-gateway/modules/dbmod"
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(dbmod.NewConn),
		fx.Provide(NewRest),
	)
}
