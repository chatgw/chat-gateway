package openaimod

import (
	"os"

	"github.com/hanyuancheung/gpt-go"
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				Key: os.Getenv("OPENAI_KEY"),
			}
		}),
		fx.Provide(func(cfg *Config) (gpt.Client, error) {
			return gpt.NewClient(cfg.Key), nil
		}),
		fx.Provide(NewChatGpt),
	)
}
