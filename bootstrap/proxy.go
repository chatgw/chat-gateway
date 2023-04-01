package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
)

type proxyDeps struct {
	fx.In
}

type Proxy struct {
	deps *proxyDeps

	mux    *chi.Mux
	server *http.Server
}

func NewRest(deps proxyDeps) *Proxy {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(render.SetContentType(render.ContentTypeHTML))

	return &Proxy{deps: &deps, mux: mux}
}

func (w *Proxy) Start() error {
	w.mux.HandleFunc("/v1/openai/*", func(w http.ResponseWriter, r *http.Request) {
		client := resty.New()
		uri := fmt.Sprintf("https://api.openai.com/%s", chi.URLParam(r, "*"))

		request := client.R().
			SetAuthToken(os.Getenv("OPENAI_KEY")).
			SetQueryString(r.URL.RawQuery)

		request.SetBody(r.Body)

		resp, err := request.Execute(r.Method, uri)

		if err != nil {
			panic(err)
		}

		w.Write(resp.Body())
	})

	w.server = &http.Server{Addr: "0.0.0.0:30120", Handler: w.mux}

	// Run the server
	return w.server.ListenAndServe()
}

func (w *Proxy) Stop() error {
	return w.server.Shutdown(context.TODO())
}
