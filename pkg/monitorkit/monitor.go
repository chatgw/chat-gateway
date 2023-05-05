package monitorkit

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var GPTRequestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "chatgw_gtp_requests_total",
		Help: "Total number of HTTP requests by skey.",
	},
	[]string{"skey"},
)

// Total Tokens = Prompt Tokens + Completion Tokens.
var GTPTokenCont = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "chatgw_gtp_token_total",
		Help: "Total number of gpt requests by skey.",
	},
	[]string{"skey"},
)

var GTPPromptTokens = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "chatgw_gtp_prompt_tokens_total",
		Help: "Total number of gpt requests by skey.",
	},
	[]string{"skey"},
)

var GTPCompletionTokens = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "chatgw_gtp_completion_tokens_total",
		Help: "Total number of gpt requests by skey.",
	},
	[]string{"skey"},
)
