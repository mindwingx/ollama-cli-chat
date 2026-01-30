package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	AcceptHeader    = "Accept"
	ApplicationJson = "application/json"
)

type Http struct {
	method  string
	url     string
	headers map[string]string
	ttl     time.Duration
	body    []byte
	client  *http.Client
	result  bool
}

func NewHttp() *Http {
	return &Http{
		headers: make(map[string]string),
		client:  &http.Client{},
	}
}

func (req *Http) Method(method string) *Http {
	req.method = method
	return req
}

func (req *Http) Url(url string) *Http {
	req.url = url
	return req
}

func (req *Http) Header(key, value string) *Http {
	req.headers[key] = value
	return req
}

func (req *Http) Ttl(t int) *Http {
	req.ttl = time.Duration(t) * time.Second
	return req
}

func (req *Http) Payload(body []byte) *Http {
	req.body = body
	return req
}

func (req *Http) GetResult() *Http {
	req.result = true
	return req
}

// DoString customized to get the Ollama handshake response as string
func (req *Http) DoString(res interface{}) (err error) {
	resp, err := request(req)

	if err != nil {
		return err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			return
		}
	}()

	// handle the string assignment properly
	val, _ := io.ReadAll(resp.Body)
	if r, ok := res.(*string); ok {
		*r = string(val)
	}

	return
}

func (req *Http) DoJson(res interface{}) (err error) {
	resp, err := request(req)

	if err != nil {
		return err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			return
		}
	}()

	if req.result == true {
		if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
			return
		}
	}

	return
}

// DoStream NDJSON/Streaming JSON
func DoStream[T any](req *Http) (context.Context, <-chan T, <-chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	respChan := make(chan T)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		res, reqErr := request(req)
		if reqErr != nil {
			errChan <- reqErr
			return
		}

		defer func() {
			if closeErr := res.Body.Close(); closeErr != nil {
				errChan <- closeErr
				return
			}
		}()

		for {
			var resp T
			if decodeErr := json.NewDecoder(res.Body).Decode(&resp); decodeErr != nil {
				if decodeErr == io.EOF {
					cancel()
					return
				}

				errChan <- errors.New("process failed. try again")
				return
			}

			respChan <- resp
		}
	}()

	return ctx, respChan, errChan
}

// HELPERS

func request(req *Http) (resp *http.Response, err error) {
	var payload *bytes.Buffer

	if len(req.body) > 0 {
		payload = bytes.NewBuffer(req.body)
	} else {
		payload = &bytes.Buffer{}
	}

	httpReq, err := http.NewRequest(req.method, req.url, payload)
	if err != nil {
		return
	}

	httpReq.Header.Set(AcceptHeader, ApplicationJson)

	if len(req.headers) > 0 {
		for k, v := range req.headers {
			httpReq.Header.Set(k, v)
		}
	}

	if req.ttl != 0 {
		req.client.Timeout = req.ttl
	}

	return req.client.Do(httpReq)
}
