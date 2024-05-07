package url

import (
	"time"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type EndpointStats struct {
	TotalCount   int64  `json:"total_count"`
	SuccessCount int64  `json:"success_count"`
	FailureCount int64  `json:"failure_count"`
	ExpiresAt    string `json:"expires_at"`
	Plan         string `json:"plan"`
}

type UrlError struct {
	Code    int
	Message string
}

func (u *UrlError) Error() string {
	return u.Message
}

// TODO: Add a request dto.
type Request struct {
	ID      int64         `json:"id"`
	Path    string        `json:"path"`
	Content pgtype.Text   `json:"content"`
	Method  db.HttpMethod `json:"method"`
	UUID    string        `json:"uuid"`

	// IPv4
	SourceIp     string             `json:"source_ip"`
	ContentSize  int32              `json:"content_size"`
	ResponseCode pgtype.Int4        `json:"response_code"`
	Headers      map[string]any     `json:"headers"`
	QueryParams  map[string]any     `json:"query_params"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	ExpiresAt    pgtype.Timestamptz `json:"expires_at"`
}

type HookRequest struct {
	Endpoint     string              `json:"endpoint"`
	UUID         string              `json:"uuid"`
	Path         string              `json:"path"`
	Headers      map[string][]string `json:"headers"`
	Query        map[string]string   `json:"queries"`
	Method       string              `json:"method"`
	SourceIp     string              `json:"source_ip"`
	Content      string              `json:"content"`
	ContentSize  int                 `json:"content_size"`
	ResponseCode int                 `json:"response_code"`
	CreatedAt    time.Time           `json:"created_at"`
	ExpiresAt    time.Time           `json:"expires_at"`
}

type Endpoint struct {
	Endpoint  string    `json:"endpoint"`
	ExpiresAt time.Time `json:"expires_at"`
	Plan      string    `json:"plan"`
}
