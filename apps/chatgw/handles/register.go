package handles

import (
	"fmt"
	"net/http"

	"github.com/airdb/chat-gateway/apps/chatgw/data/repos"
	sensitivemod "github.com/airdb/chat-gateway/modules/sensitive"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type registerDeps struct {
	fx.In

	Mux     *chi.Mux
	Logger  *slog.Logger
	KeyRepo *repos.KeyRepo
	Checker *sensitivemod.Checker
}

func Register(deps registerDeps) {
	deps.Mux.Route("/v1", func(r chi.Router) {
		r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			log := deps.Logger.With("uri", r.URL.String())

			log.Debug("ping")
			fmt.Fprintf(w, "pong\n")
		})
		r.HandleFunc("/sensitive", func(w http.ResponseWriter, r *http.Request) {
			log := deps.Logger.With("uri", r.URL.String())

			search := r.URL.Query().Get("s")
			if search == "" {
				w.Write([]byte("s(search) is empty"))
				return
			}
			/*
				defer r.Body.Close()
				body, _ := io.ReadAll(r.Body)
			*/
			log.Debug("Get :" + search)
			result := deps.Checker.HasSense([]byte(search))
			fmt.Fprintf(w, "check result:"+fmt.Sprintf("%v", result))
		})
		r.HandleFunc("/openai/*", deps.HandleOpenai)
		r.HandleFunc("/azure", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("waiting for implement\n"))
		})
	})
}
