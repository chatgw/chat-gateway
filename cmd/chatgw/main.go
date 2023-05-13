package main

import (
	"context"
	"errors"
	"log"

	"github.com/airdb/chat-gateway/apps/chatgw"
	"github.com/airdb/chat-gateway/modules/proxymod"
	sensitivemod "github.com/airdb/chat-gateway/modules/sensitive"

	"github.com/airdb/chat-gateway/bootstrap"
	telemetrymod "github.com/airdb/chat-gateway/modules/telemetry"
	"github.com/airdb/chat-gateway/pkg/lokikit"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	type invokeDeps struct {
		fx.In

		LokiWriter *lokikit.LokiWriter
		Proxy      *proxymod.Proxy
	}

	app := fx.New(
		telemetrymod.FxOptions(),
		bootstrap.FxOptions(),
		chatgw.FxOptions(),
		sensitivemod.FxOptions(),
		fx.Invoke(func(lc fx.Lifecycle, deps invokeDeps) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go deps.Proxy.Start()
					log.Println("Press Ctrl+C to exit")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					deps.LokiWriter.Shutdown()
					return errors.Join(deps.Proxy.Stop())
				},
			})
		}),
	)

	app.Run()
}
