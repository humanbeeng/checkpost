// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

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
	PlanGuest     Plan = "guest"
	PlanFree      Plan = "free"
	PlanNoBrainer Plan = "no_brainer"
	PlanPro       Plan = "pro"
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

type FileAttachment struct {
	ID        int64            `json:"id"`
	Uri       string           `json:"uri"`
	UrlID     int64            `json:"url_id"`
	UserID    pgtype.Int8      `json:"user_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Request struct {
	ID         int64       `json:"id"`
	UserID     pgtype.Int8 `json:"user_id"`
	UrlID      int64       `json:"url_id"`
	ResponseID pgtype.Int8 `json:"response_id"`
	Content    pgtype.Text `json:"content"`
	Method     HttpMethod  `json:"method"`
	// IPv4
	SourceIp     string           `json:"source_ip"`
	ContentSize  int32            `json:"content_size"`
	ResponseCode pgtype.Int4      `json:"response_code"`
	Headers      []byte           `json:"headers"`
	QueryParams  []byte           `json:"query_params"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
}

type Response struct {
	ID           int64            `json:"id"`
	UserID       pgtype.Int8      `json:"user_id"`
	UrlID        int64            `json:"url_id"`
	ResponseCode int32            `json:"response_code"`
	Content      pgtype.Text      `json:"content"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
}

type Url struct {
	ID        int64            `json:"id"`
	Url       string           `json:"url"`
	UserID    pgtype.Int8      `json:"user_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Plan      Plan             `json:"plan"`
}

type User struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Username  string           `json:"username"`
	Plan      Plan             `json:"plan"`
	Email     string           `json:"email"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
