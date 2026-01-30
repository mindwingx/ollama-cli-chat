package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	req "ollama-cli/pkg"
	"os"
	"time"
)

const (
	urlHttpSchema  string = "http"
	urlHttpsSchema string = "https"
	baseUri        string = "/"          // GET method
	chatUri        string = "api/chat"   // POST method
	modelsUri      string = "api/tags"   // GET method
	pullUri        string = "api/pull"   // POST method
	deleteUri      string = "api/delete" // DELETE method
)

type (
	ChatRequest struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

type (
	ChatResponse struct {
		Model              string    `json:"model"`
		CreatedAt          time.Time `json:"created_at"`
		Message            Content   `json:"message"`
		Done               bool      `json:"done"`
		DoneReason         string    `json:"done_reason"`
		TotalDuration      int       `json:"total_duration"`
		LoadDuration       int       `json:"load_duration"`
		PromptEvalCount    int       `json:"prompt_eval_count"`
		PromptEvalDuration int       `json:"prompt_eval_duration"`
		EvalCount          int       `json:"eval_count"`
		EvalDuration       int       `json:"eval_duration"`
	}

	Content struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

type (
	ModelResponse struct {
		Models []Model `json:"models"`
	}

	Model struct {
		Name       string       `json:"name"`
		Model      string       `json:"model"`
		ModifiedAt time.Time    `json:"modified_at"`
		Size       int          `json:"size"`
		Digest     string       `json:"digest"`
		Details    ModelDetails `json:"details"`
	}

	ModelDetails struct {
		ParentModel       string   `json:"parent_model"`
		Format            string   `json:"format"`
		Family            string   `json:"family"`
		Families          []string `json:"families"`
		ParameterSize     string   `json:"parameter_size"`
		QuantizationLevel string   `json:"quantization_level"`
	}
)

type (
	ModelRequest struct {
		Model string `json:"model"`
	}

	ModelPullResponse struct {
		Status    string `json:"status,omitempty"`
		Digest    string `json:"digest,omitempty"`
		Total     int64  `json:"total,omitempty"`
		Completed int64  `json:"completed,omitempty"`
		Error     string `json:"error,omitempty"`
	}
)

//

type IOllamaProvider interface {
	Handshake() (string, error)
	SendMessage(msg ChatRequest) (ChatResponse, error)
	GetModels() (ModelResponse, error)
	PullModel(model ModelRequest) (context.Context, <-chan ModelPullResponse, <-chan error)
	DeleteModel(model ModelRequest) error
}

//

type OllamaProvider struct {
	BaseUrl string
}

func NewOllamaProvider() IOllamaProvider {
	baseUrl := os.Getenv("OLLAMA_BASE_URL")

	if len(baseUrl) == 0 {
		baseUrl = "0.0.0.0:11434"
	}

	return &OllamaProvider{BaseUrl: baseUrl}
}

func (o *OllamaProvider) Handshake() (res string, err error) {
	addr := url.URL{}
	{
		addr.Scheme = urlHttpSchema
		addr.Host = o.BaseUrl
		addr.Path = baseUri
	}

	var resp string

	request := req.NewHttp().Ttl(3).Url(addr.String()).Method(http.MethodGet)
	if err = request.DoString(&resp); err != nil {
		return
	}

	res = resp
	return
}

func (o *OllamaProvider) SendMessage(msg ChatRequest) (res ChatResponse, err error) {
	resp := ChatResponse{}

	addr := url.URL{}
	{
		addr.Scheme = urlHttpSchema
		addr.Host = o.BaseUrl
		addr.Path = chatUri
	}

	payload, _ := json.Marshal(msg)

	request := req.NewHttp().Url(addr.String()).Payload(payload).GetResult().Method(http.MethodPost)
	if err = request.DoJson(&resp); err != nil {
		return
	}

	res = resp
	return
}

func (o *OllamaProvider) GetModels() (res ModelResponse, err error) {
	resp := ModelResponse{}

	addr := url.URL{}
	{
		addr.Scheme = urlHttpSchema
		addr.Host = o.BaseUrl
		addr.Path = modelsUri
	}

	request := req.NewHttp().Ttl(3).Url(addr.String()).GetResult().Method(http.MethodGet)
	if err = request.DoJson(&resp); err != nil {
		return
	}

	res = resp
	return
}

func (o *OllamaProvider) PullModel(model ModelRequest) (context.Context, <-chan ModelPullResponse, <-chan error) {
	addr := url.URL{
		Scheme: urlHttpSchema,
		Host:   o.BaseUrl,
		Path:   pullUri,
	}

	payload, _ := json.Marshal(model)

	request := req.NewHttp().Url(addr.String()).Payload(payload).Method(http.MethodPost)
	ctx, respChan, errChan := req.DoStream[ModelPullResponse](request)

	return ctx, respChan, errChan
}

func (o *OllamaProvider) DeleteModel(model ModelRequest) (err error) {
	addr := url.URL{}
	{
		addr.Scheme = urlHttpSchema
		addr.Host = o.BaseUrl
		addr.Path = deleteUri
	}

	payload, _ := json.Marshal(model)

	request := req.NewHttp().Url(addr.String()).Payload(payload).Method(http.MethodDelete)
	if err = request.DoJson(nil); err != nil {
		return
	}

	return
}
