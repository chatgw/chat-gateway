package telemetrymod

import (
	"os"

	"github.com/airdb/chat-gateway/pkg/logkit"
	"github.com/airdb/chat-gateway/pkg/lokikit"
	"github.com/gofrs/uuid"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(func() (*lokikit.LokiWriter, error) {
			instanceId, err := uuid.NewV6()
			if err != nil {
				return nil, err
			}
			writer, err := lokikit.NewLokiWriter(
				os.Getenv("LOKI_URL"), 0,
				lokikit.WithBasicAuth(
					os.Getenv("LOKI_USER"),
					os.Getenv("LOKI_PASSWORD"),
				),
				lokikit.WithLabels(map[string]string{
					"service":  "chatgw",
					"instance": instanceId.String(),
				}),
				lokikit.WithFields([]string{
					"level", "requestID", "user", "command",
				}),
			)
			if err != nil {
				return nil, err
			}
			return writer, nil
		}),
		fx.Provide(func(writer *lokikit.LokiWriter) (*slog.Logger, error) {
			logger, err := logkit.New(nil)
			if handler, ok := logger.Handler().(*logkit.TeeHandler); ok {
				var leveler logkit.LevelEnablerFunc = func(level slog.Level) bool {
					return true
				}
				handler.AppendHandlers(logkit.NewLevelHandler(
					leveler, handler.HandlerOptions().NewJSONHandler(writer),
				))
			}
			return logger, err
		}),
	)
}
