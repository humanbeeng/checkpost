// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: request.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNewRequest = `-- name: CreateNewRequest :one
INSERT INTO
    request (
        user_id,
        endpoint_id,
        PATH,
        response_id,
        CONTENT,
        METHOD,
        source_ip,
        content_size,
        response_code,
        headers,
        query_params,
        expires_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING
    id, user_id, endpoint_id, path, response_id, content, method, source_ip, content_size, response_code, headers, query_params, created_at, expires_at, is_deleted
`

type CreateNewRequestParams struct {
	UserID       pgtype.Int8      `json:"user_id"`
	EndpointID   int64            `json:"endpoint_id"`
	Path         string           `json:"path"`
	ResponseID   pgtype.Int8      `json:"response_id"`
	Content      pgtype.Text      `json:"content"`
	Method       HttpMethod       `json:"method"`
	SourceIp     string           `json:"source_ip"`
	ContentSize  int32            `json:"content_size"`
	ResponseCode pgtype.Int4      `json:"response_code"`
	Headers      []byte           `json:"headers"`
	QueryParams  []byte           `json:"query_params"`
	ExpiresAt    pgtype.Timestamp `json:"expires_at"`
}

func (q *Queries) CreateNewRequest(ctx context.Context, arg CreateNewRequestParams) (Request, error) {
	row := q.db.QueryRow(ctx, createNewRequest,
		arg.UserID,
		arg.EndpointID,
		arg.Path,
		arg.ResponseID,
		arg.Content,
		arg.Method,
		arg.SourceIp,
		arg.ContentSize,
		arg.ResponseCode,
		arg.Headers,
		arg.QueryParams,
		arg.ExpiresAt,
	)
	var i Request
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EndpointID,
		&i.Path,
		&i.ResponseID,
		&i.Content,
		&i.Method,
		&i.SourceIp,
		&i.ContentSize,
		&i.ResponseCode,
		&i.Headers,
		&i.QueryParams,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}

const getEndpointHistory = `-- name: GetEndpointHistory :many
SELECT
    request.id, request.user_id, endpoint_id, path, response_id, content, method, source_ip, content_size, response_code, headers, query_params, request.created_at, request.expires_at, request.is_deleted, endpoint.id, endpoint, endpoint.user_id, plan, endpoint.created_at, endpoint.expires_at, endpoint.is_deleted
FROM
    request
    LEFT JOIN endpoint ON request.endpoint_id = endpoint.id
WHERE
    endpoint.endpoint = $1
    AND request.is_deleted = FALSE
LIMIT
    $2
OFFSET
    $3
`

type GetEndpointHistoryParams struct {
	Endpoint string `json:"endpoint"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type GetEndpointHistoryRow struct {
	ID           int64            `json:"id"`
	UserID       pgtype.Int8      `json:"user_id"`
	EndpointID   int64            `json:"endpoint_id"`
	Path         string           `json:"path"`
	ResponseID   pgtype.Int8      `json:"response_id"`
	Content      pgtype.Text      `json:"content"`
	Method       HttpMethod       `json:"method"`
	SourceIp     string           `json:"source_ip"`
	ContentSize  int32            `json:"content_size"`
	ResponseCode pgtype.Int4      `json:"response_code"`
	Headers      []byte           `json:"headers"`
	QueryParams  []byte           `json:"query_params"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	ExpiresAt    pgtype.Timestamp `json:"expires_at"`
	IsDeleted    pgtype.Bool      `json:"is_deleted"`
	ID_2         pgtype.Int8      `json:"id_2"`
	Endpoint     pgtype.Text      `json:"endpoint"`
	UserID_2     pgtype.Int8      `json:"user_id_2"`
	Plan         NullPlan         `json:"plan"`
	CreatedAt_2  pgtype.Timestamp `json:"created_at_2"`
	ExpiresAt_2  pgtype.Timestamp `json:"expires_at_2"`
	IsDeleted_2  pgtype.Bool      `json:"is_deleted_2"`
}

func (q *Queries) GetEndpointHistory(ctx context.Context, arg GetEndpointHistoryParams) ([]GetEndpointHistoryRow, error) {
	rows, err := q.db.Query(ctx, getEndpointHistory, arg.Endpoint, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEndpointHistoryRow{}
	for rows.Next() {
		var i GetEndpointHistoryRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.EndpointID,
			&i.Path,
			&i.ResponseID,
			&i.Content,
			&i.Method,
			&i.SourceIp,
			&i.ContentSize,
			&i.ResponseCode,
			&i.Headers,
			&i.QueryParams,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.IsDeleted,
			&i.ID_2,
			&i.Endpoint,
			&i.UserID_2,
			&i.Plan,
			&i.CreatedAt_2,
			&i.ExpiresAt_2,
			&i.IsDeleted_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEndpointRequestCount = `-- name: GetEndpointRequestCount :one
SELECT
    COUNT(*) AS total_count,
    COUNT(
        CASE
            WHEN response_code = 200 THEN 1
        END
    ) AS success_count,
    COUNT(
        CASE
            WHEN response_code != 200 THEN 1
        END
    ) AS failure_count
FROM
    request r
    LEFT JOIN endpoint e ON r.endpoint_id = e.id
WHERE
    endpoint = $1
    AND is_deleted = FALSE
`

type GetEndpointRequestCountRow struct {
	TotalCount   int64 `json:"total_count"`
	SuccessCount int64 `json:"success_count"`
	FailureCount int64 `json:"failure_count"`
}

func (q *Queries) GetEndpointRequestCount(ctx context.Context, endpoint string) (GetEndpointRequestCountRow, error) {
	row := q.db.QueryRow(ctx, getEndpointRequestCount, endpoint)
	var i GetEndpointRequestCountRow
	err := row.Scan(&i.TotalCount, &i.SuccessCount, &i.FailureCount)
	return i, err
}

const getRequestById = `-- name: GetRequestById :one
SELECT
    id, user_id, endpoint_id, path, response_id, content, method, source_ip, content_size, response_code, headers, query_params, created_at, expires_at, is_deleted
FROM
    request
WHERE
    id = $1
    AND is_deleted = FALSE
LIMIT
    1
`

func (q *Queries) GetRequestById(ctx context.Context, id int64) (Request, error) {
	row := q.db.QueryRow(ctx, getRequestById, id)
	var i Request
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EndpointID,
		&i.Path,
		&i.ResponseID,
		&i.Content,
		&i.Method,
		&i.SourceIp,
		&i.ContentSize,
		&i.ResponseCode,
		&i.Headers,
		&i.QueryParams,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}
