package handles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/airdb/chat-gateway/modules/openaimod"
	"github.com/airdb/chat-gateway/pkg/monitorkit"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
)

func (deps registerDeps) HandleOpenai(w http.ResponseWriter, r *http.Request) {
	skey := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	token := skey

	logEntry := deps.Logger.
		With("uri", r.URL.String()).
		With("token", token)

	key, err := deps.KeyRepo.First(r.Context(), token)

	logEntry.Error("error key", skey, key)
	if err == nil && key != nil {
		token = os.Getenv("OPENAI_KEY")
	}

	client := resty.New()
	uri := fmt.Sprintf("https://api.openai.com/%s", chi.URLParam(r, "*"))

	request := client.R().
		SetAuthToken(token).
		SetQueryString(r.URL.RawQuery)

	rDumper := bytes.NewBuffer(nil)
	body := io.TeeReader(r.Body, rDumper)
	request.SetBody(body)

	resp, err := request.Execute(r.Method, uri)
	parseBody(logEntry, rDumper.Bytes()).Debug("request body")

	if err != nil {
		panic(err)
	}

	parseBody(logEntry, resp.Body()).Debug("response body")
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
}
