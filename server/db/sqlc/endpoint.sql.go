// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: endpoint.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNewEndpoint = `-- name: CreateNewEndpoint :one
insert into
  endpoint (endpoint, user_id, plan)
values
($1, $2, $3) returning id, endpoint, user_id, created_at, plan
`

type CreateNewEndpointParams struct {
	Endpoint string      `json:"endpoint"`
	UserID   pgtype.Int8 `json:"user_id"`
	Plan     Plan        `json:"plan"`
}

func (q *Queries) CreateNewEndpoint(ctx context.Context, arg CreateNewEndpointParams) (Endpoint, error) {
	row := q.db.QueryRow(ctx, createNewEndpoint, arg.Endpoint, arg.UserID, arg.Plan)
	var i Endpoint
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.UserID,
		&i.CreatedAt,
		&i.Plan,
	)
	return i, err
}

const getEndpoint = `-- name: GetEndpoint :one
select id, endpoint, user_id, created_at, plan from "endpoint" where endpoint = $1 limit 1
`

func (q *Queries) GetEndpoint(ctx context.Context, endpoint string) (Endpoint, error) {
	row := q.db.QueryRow(ctx, getEndpoint, endpoint)
	var i Endpoint
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.UserID,
		&i.CreatedAt,
		&i.Plan,
	)
	return i, err
}