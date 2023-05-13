package bootstrap

import (
	"github.com/airdb/chat-gateway/apps/chatgw/handles"
	"github.com/airdb/chat-gateway/modules/dbmod"
	"github.com/airdb/chat-gateway/modules/proxymod"
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		proxymod.FxOptions(),
		fx.Provide(dbmod.NewConn),
		fx.Invoke(handles.Register),
	)
}
