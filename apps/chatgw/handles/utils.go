package handles

import (
	"encoding/json"

	"golang.org/x/exp/slog"
)

func parseBody(logEntry *slog.Logger, body []byte) *slog.Logger {
	if len(body) == 0 {
		return logEntry
	}
	data := map[string]any{}
	if err := json.Unmarshal(body, &data); err == nil {
		return logEntry.With("body", data)
	} else {
		return logEntry.With("body", string(body)).With("err", err)
	}
}
