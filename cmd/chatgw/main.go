package main

import (
	"context"
	"errors"
	"log"

	"github.com/airdb/chat-gateway/bootstrap"
	telemetrymod "github.com/airdb/chat-gateway/modules/telemetry"
	"github.com/airdb/chat-gateway/pkg/lokikit"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	type invokeDeps struct {
		fx.In

		LokiWriter *lokikit.LokiWriter
		Rest       *bootstrap.Proxy
	}

	app := fx.New(
		telemetrymod.FxOptions(),
		bootstrap.FxOptions(),
		fx.Invoke(func(lc fx.Lifecycle, deps invokeDeps) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go deps.Rest.Start()
					log.Println("Press Ctrl+C to exit")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					deps.LokiWriter.Shutdown()
					return errors.Join(deps.Rest.Stop())
				},
			})
		}),
	)

	app.Run()
}
