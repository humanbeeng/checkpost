// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CheckEndpointExists(ctx context.Context, endpoint string) (bool, error)
	CreateNewRequest(ctx context.Context, arg CreateNewRequestParams) (Request, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteExpiredRequests(ctx context.Context) error
	DeleteUser(ctx context.Context, id int64) error
	GetEndpoint(ctx context.Context, endpoint string) (Endpoint, error)
	GetEndpointHistory(ctx context.Context, arg GetEndpointHistoryParams) ([]GetEndpointHistoryRow, error)
	GetEndpointRequestCount(ctx context.Context, endpoint string) (GetEndpointRequestCountRow, error)
	GetNonExpiredEndpointsOfUser(ctx context.Context, userID pgtype.Int8) ([]Endpoint, error)
	GetRequestById(ctx context.Context, id int64) (Request, error)
	GetRequestByUUID(ctx context.Context, uuid string) (Request, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserEndpoints(ctx context.Context, userID pgtype.Int8) ([]Endpoint, error)
	GetUserFromEmail(ctx context.Context, email string) (User, error)
	GetUserFromUsername(ctx context.Context, username string) (User, error)
	InsertEndpoint(ctx context.Context, arg InsertEndpointParams) (Endpoint, error)
	InsertFreeEndpoint(ctx context.Context, arg InsertFreeEndpointParams) (Endpoint, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
}

var _ Querier = (*Queries)(nil)
