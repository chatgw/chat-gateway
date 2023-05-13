package proxymod

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
)

type Proxy struct {
	mux    *chi.Mux
	server *http.Server
}

func New(mux *chi.Mux) *Proxy {
	proxy := &Proxy{
		mux:    mux,
		server: &http.Server{Addr: ":30120", Handler: mux},
	}

	return proxy
}

func (p *Proxy) Start() error {
	p.mux.Route("/", func(r chi.Router) {
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to chat-gateway\n"))
		})
		r.Handle("/metrics", promhttp.Handler())
	})

	// Run the server
	return p.server.ListenAndServe()
}

func (p *Proxy) Stop() error {
	return p.server.Shutdown(context.TODO())
}
