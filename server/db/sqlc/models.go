// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "get"
	HttpMethodPost    HttpMethod = "post"
	HttpMethodPut     HttpMethod = "put"
	HttpMethodPatch   HttpMethod = "patch"
	HttpMethodDelete  HttpMethod = "delete"
	HttpMethodOptions HttpMethod = "options"
	HttpMethodHead    HttpMethod = "head"
	HttpMethodTrace   HttpMethod = "trace"
	HttpMethodConnect HttpMethod = "connect"
)

func (e *HttpMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = HttpMethod(s)
	case string:
		*e = HttpMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for HttpMethod: %T", src)
	}
	return nil
}

type NullHttpMethod struct {
	HttpMethod HttpMethod `json:"http_method"`
	Valid      bool       `json:"valid"` // Valid is true if HttpMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullHttpMethod) Scan(value interface{}) error {
	if value == nil {
		ns.HttpMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.HttpMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullHttpMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.HttpMethod), nil
}

type Plan string

const (
	PlanFree  Plan = "free"
	PlanBasic Plan = "basic"
	PlanPro   Plan = "pro"
)

func (e *Plan) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Plan(s)
	case string:
		*e = Plan(s)
	default:
		return fmt.Errorf("unsupported scan type for Plan: %T", src)
	}
	return nil
}

type NullPlan struct {
	Plan  Plan `json:"plan"`
	Valid bool `json:"valid"` // Valid is true if Plan is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPlan) Scan(value interface{}) error {
	if value == nil {
		ns.Plan, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Plan.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPlan) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Plan), nil
}

type Endpoint struct {
	ID        int64              `json:"id"`
	Endpoint  string             `json:"endpoint"`
	UserID    pgtype.Int8        `json:"user_id"`
	Plan      Plan               `json:"plan"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
	IsDeleted pgtype.Bool        `json:"is_deleted"`
}

type FileAttachment struct {
	ID         int64              `json:"id"`
	Uri        string             `json:"uri"`
	EndpointID int64              `json:"endpoint_id"`
	UserID     pgtype.Int8        `json:"user_id"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	IsDeleted  pgtype.Bool        `json:"is_deleted"`
}

type Request struct {
	ID           int64       `json:"id"`
	Uuid         string      `json:"uuid"`
	UserID       pgtype.Int8 `json:"user_id"`
	EndpointID   int64       `json:"endpoint_id"`
	Plan         Plan        `json:"plan"`
	Path         string      `json:"path"`
	ResponseID   pgtype.Int8 `json:"response_id"`
	ResponseTime pgtype.Int4 `json:"response_time"`
	Content      pgtype.Text `json:"content"`
	ContentType  string      `json:"content_type"`
	Method       HttpMethod  `json:"method"`
	// IPv4
	SourceIp     string             `json:"source_ip"`
	ContentSize  int32              `json:"content_size"`
	ResponseCode pgtype.Int4        `json:"response_code"`
	Headers      []byte             `json:"headers"`
	FormData     []byte             `json:"form_data"`
	QueryParams  []byte             `json:"query_params"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	ExpiresAt    pgtype.Timestamptz `json:"expires_at"`
	IsDeleted    pgtype.Bool        `json:"is_deleted"`
}

type Response struct {
	ID           int64              `json:"id"`
	UserID       pgtype.Int8        `json:"user_id"`
	EndpointID   int64              `json:"endpoint_id"`
	ResponseCode int32              `json:"response_code"`
	Content      pgtype.Text        `json:"content"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	IsDeleted    pgtype.Bool        `json:"is_deleted"`
}

type User struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	AvatarUrl string             `json:"avatar_url"`
	Username  string             `json:"username"`
	Plan      Plan               `json:"plan"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	IsDeleted pgtype.Bool        `json:"is_deleted"`
}
