package lokikit

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const LOG_ENTRIES_CHAN_SIZE = 5000

type lokiClientConfig struct {
	PushURL            string
	Labels             map[string]string
	BatchWait          time.Duration
	BatchEntriesNumber int
	Fields             []string
}

// http.Client wrapper for adding new methods, particularly sendJsonReq
type lokiClient struct {
	parent    *http.Client
	beforeDos []func(*http.Request)
}

func newLokiClient() *lokiClient {
	return &lokiClient{
		parent:    &http.Client{},
		beforeDos: make([]func(*http.Request), 0),
	}
}

// A bit more convenient method for sending requests to the HTTP server
func (client *lokiClient) sendJsonReq(
	method, url string, ctype string, reqBody []byte,
) (resp *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", ctype)

	for _, f := range client.beforeDos {
		f(req)
	}

	resp, err = client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBody, nil
}
