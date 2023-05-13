package chatgw

import (
	"github.com/airdb/chat-gateway/apps/chatgw/data"
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		data.FxOptions(),
		fx.Invoke(data.Migrate),
	)
}
