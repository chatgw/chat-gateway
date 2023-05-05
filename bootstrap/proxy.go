package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/airdb/chat-gateway/modules/openaimod"
	sensitivemod "github.com/airdb/chat-gateway/modules/sensitive"
	"github.com/airdb/chat-gateway/pkg/monitorkit"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type proxyDeps struct {
	fx.In

	Logger  *slog.Logger
	Checker *sensitivemod.Checker
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

func (p *Proxy) Start() error {
	p.mux.Route("/", func(r chi.Router) {
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to chat-gateway\n"))
		})
		r.Handle("/metrics", promhttp.Handler())
	})

	p.mux.Route("/v1", func(r chi.Router) {
		r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			log := p.deps.Logger.With("uri", r.URL.String())

			log.Debug("ping")
			fmt.Fprintf(w, "pong\n")
		})
		r.HandleFunc("/sensitive", func(w http.ResponseWriter, r *http.Request) {
			log := p.deps.Logger.With("uri", r.URL.String())

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
			result := p.deps.Checker.HasSense([]byte(search))
			fmt.Fprintf(w, "check result:"+fmt.Sprintf("%v", result))

		})
		r.HandleFunc("/openai/*", func(w http.ResponseWriter, r *http.Request) {
			skey := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			token := skey
			if len(skey) < 40 {
				validSkeys := os.Getenv("CHATGW_TOKEN")
				if strings.Contains(validSkeys, skey) {
					token = os.Getenv("OPENAI_KEY")
				}
			}

			logEntry := p.deps.Logger.
				With("uri", r.URL.String()).
				With("token", token)

			client := resty.New()
			uri := fmt.Sprintf("https://api.openai.com/%s", chi.URLParam(r, "*"))

			request := client.R().
				SetAuthToken(token).
				SetQueryString(r.URL.RawQuery)

			rDumper := bytes.NewBuffer(nil)
			body := io.TeeReader(r.Body, rDumper)
			request.SetBody(body)

			resp, err := request.Execute(r.Method, uri)
			p.parseBody(logEntry, rDumper.Bytes()).Debug("request body")

			if err != nil {
				panic(err)
			}

			p.parseBody(logEntry, resp.Body()).Debug("response body")
			w.Write(resp.Body())

			monitorkit.GPTRequestCount.WithLabelValues(skey).Inc()

			var chatGPTResp openaimod.ChatGPTResp

			err = json.Unmarshal(resp.Body(), &chatGPTResp)
			if err != nil {
				log.Println()
				return
			}

			log.Println("token count", chatGPTResp.Usage.TotalTokens)

			monitorkit.GTPTokenCont.WithLabelValues(skey).Add(float64(chatGPTResp.Usage.TotalTokens))
			monitorkit.GTPPromptTokens.WithLabelValues(skey).Add(float64(chatGPTResp.Usage.PromptTokens))
			monitorkit.GTPCompletionTokens.WithLabelValues(skey).Add(float64(chatGPTResp.Usage.CompletionTokens))
		})
		r.HandleFunc("/azure", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("waiting for implement\n"))
		})

	})

	p.server = &http.Server{Addr: ":30120", Handler: p.mux}

	// Run the server
	return p.server.ListenAndServe()
}

func (p *Proxy) Stop() error {
	return p.server.Shutdown(context.TODO())
}

func (p *Proxy) parseBody(logEntry *slog.Logger, body []byte) *slog.Logger {
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
