package endpoint

import (
	"encoding/json"
	"time"
)

type EndpointStats struct {
	TotalCount   int64  `json:"total_count"`
	SuccessCount int64  `json:"success_count"`
	FailureCount int64  `json:"failure_count"`
	ExpiresAt    string `json:"expires_at"`
	Plan         string `json:"plan"`
}

type EndpointError struct {
	Code    int
	Message string
}

func (u *EndpointError) Error() string {
	return u.Message
}

type ContentType string

const (
	ApplicationJson ContentType = "application/json"
	TextYaml        ContentType = "text/yaml"
	FormUrlEncoded  ContentType = "application/x-www-form-urlencoded"
	MultipartForm   ContentType = "multipart/form-data"
)

type HookRequest struct {
	Endpoint     string              `json:"endpoint"`
	UUID         string              `json:"uuid"`
	Path         string              `json:"path"`
	Headers      map[string][]string `json:"headers"`
	QueryParams  map[string]string   `json:"query_params"`
	FormData     map[string][]string `json:"form_data"`
	Method       string              `json:"method"`
	SourceIp     string              `json:"source_ip"`
	Content      string              `json:"content"`
	ContentType  string              `json:"content_type"`
	ContentSize  int32               `json:"content_size"`
	ResponseCode int32               `json:"response_code"`
	CreatedAt    time.Time           `json:"created_at"`
	ExpiresAt    time.Time           `json:"expires_at"`
}

type Endpoint struct {
	Endpoint  string    `json:"endpoint"`
	ExpiresAt time.Time `json:"expires_at"`
	Plan      string    `json:"plan"`
}

type WSMessage struct {
	Code    int             `json:"code"`
	Payload json.RawMessage `json:"payload"`
	Message string          `json:"message"`
}
