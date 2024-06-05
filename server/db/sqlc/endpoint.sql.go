// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: endpoint.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkEndpointExists = `-- name: CheckEndpointExists :one
SELECT
    EXISTS (
        SELECT
            id, endpoint, user_id, plan, created_at, expires_at, is_deleted
        FROM
            endpoint
        WHERE
            endpoint = $1
            AND expires_at > NOW()
            AND is_deleted = FALSE
        LIMIT
            1
    )
`

func (q *Queries) CheckEndpointExists(ctx context.Context, endpoint string) (bool, error) {
	row := q.db.QueryRow(ctx, checkEndpointExists, endpoint)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getEndpointDetails = `-- name: GetEndpointDetails :one
SELECT
    id, endpoint, user_id, plan, created_at, expires_at, is_deleted
FROM
    "endpoint"
WHERE
    endpoint = $1
    AND is_deleted = FALSE
LIMIT
    1
`

func (q *Queries) GetEndpointDetails(ctx context.Context, endpoint string) (Endpoint, error) {
	row := q.db.QueryRow(ctx, getEndpointDetails, endpoint)
	var i Endpoint
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.UserID,
		&i.Plan,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}

const getNonExpiredEndpointsOfUser = `-- name: GetNonExpiredEndpointsOfUser :many
SELECT
    id, endpoint, user_id, plan, created_at, expires_at, is_deleted
FROM
    "endpoint"
WHERE
    user_id = $1
    AND expires_at > NOW()
    AND is_deleted = FALSE
`

func (q *Queries) GetNonExpiredEndpointsOfUser(ctx context.Context, userID pgtype.Int8) ([]Endpoint, error) {
	rows, err := q.db.Query(ctx, getNonExpiredEndpointsOfUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Endpoint{}
	for rows.Next() {
		var i Endpoint
		if err := rows.Scan(
			&i.ID,
			&i.Endpoint,
			&i.UserID,
			&i.Plan,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.IsDeleted,
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

const getUserEndpoints = `-- name: GetUserEndpoints :many
SELECT
    id, endpoint, user_id, plan, created_at, expires_at, is_deleted
FROM
    "endpoint"
WHERE
    user_id = $1
    AND is_deleted = FALSE
`

func (q *Queries) GetUserEndpoints(ctx context.Context, userID pgtype.Int8) ([]Endpoint, error) {
	rows, err := q.db.Query(ctx, getUserEndpoints, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Endpoint{}
	for rows.Next() {
		var i Endpoint
		if err := rows.Scan(
			&i.ID,
			&i.Endpoint,
			&i.UserID,
			&i.Plan,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.IsDeleted,
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

const insertEndpoint = `-- name: InsertEndpoint :one
INSERT INTO
    endpoint (endpoint, user_id, plan, expires_at)
VALUES
    ($1, $2, $3, $4)
RETURNING
    id, endpoint, user_id, plan, created_at, expires_at, is_deleted
`

type InsertEndpointParams struct {
	Endpoint  string             `json:"endpoint"`
	UserID    pgtype.Int8        `json:"user_id"`
	Plan      Plan               `json:"plan"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) InsertEndpoint(ctx context.Context, arg InsertEndpointParams) (Endpoint, error) {
	row := q.db.QueryRow(ctx, insertEndpoint,
		arg.Endpoint,
		arg.UserID,
		arg.Plan,
		arg.ExpiresAt,
	)
	var i Endpoint
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.UserID,
		&i.Plan,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}

const insertFreeEndpoint = `-- name: InsertFreeEndpoint :one
INSERT INTO
    endpoint (endpoint, user_id, plan, expires_at)
VALUES
    ($1, $2, 'free', $3)
RETURNING
    id, endpoint, user_id, plan, created_at, expires_at, is_deleted
`

type InsertFreeEndpointParams struct {
	Endpoint  string             `json:"endpoint"`
	UserID    pgtype.Int8        `json:"user_id"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) InsertFreeEndpoint(ctx context.Context, arg InsertFreeEndpointParams) (Endpoint, error) {
	row := q.db.QueryRow(ctx, insertFreeEndpoint, arg.Endpoint, arg.UserID, arg.ExpiresAt)
	var i Endpoint
	err := row.Scan(
		&i.ID,
		&i.Endpoint,
		&i.UserID,
		&i.Plan,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.IsDeleted,
	)
	return i, err
}
