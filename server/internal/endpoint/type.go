package endpoint

import (
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

type HookRequest struct {
	Endpoint     string              `json:"endpoint"`
	UUID         string              `json:"uuid"`
	Path         string              `json:"path"`
	Headers      map[string][]string `json:"headers"`
	QueryParams  map[string]string   `json:"query_params"`
	Method       string              `json:"method"`
	SourceIp     string              `json:"source_ip"`
	Content      string              `json:"content"`
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
