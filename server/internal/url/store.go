package url

import (
	"context"
	"strings"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type UrlQuerier interface {
	CheckEndpointExists(ctx context.Context, endpoint string) (bool, error)

	GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error)
	GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error)
	GetUserEndpoints(ctx context.Context, userId int64) ([]db.Endpoint, error)
	GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error)
	GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error)

	InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error)
	InsertEndpoint(ctx context.Context, params db.InsertEndpointParams) (db.Endpoint, error)

	// TODO: Move these to requests querier
	CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error)

	GetRequestById(ctx context.Context, reqId int64) (db.Request, error)
	GetRequestByUUID(ctx context.Context, uuid string) (db.Request, error)

	ExpireRequests(ctx context.Context) error
}

type UrlStore struct {
	q db.Querier
}

func NewUrlStore(q db.Querier) *UrlStore {
	return &UrlStore{
		q: q,
	}
}

func (us UrlStore) GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error) {
	return us.q.GetEndpointRequestCount(ctx, endpoint)
}

func (us UrlStore) CheckEndpointExists(ctx context.Context, endpoint string) (bool, error) {
	return us.q.CheckEndpointExists(ctx, endpoint)
}

func (us UrlStore) GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error) {
	return us.q.GetEndpoint(ctx, endpoint)
}

func (us UrlStore) GetUserEndpoints(ctx context.Context, userId int64) ([]db.Endpoint, error) {
	return us.q.GetUserEndpoints(ctx, pgtype.Int8{Int64: userId, Valid: true})
}

func (us UrlStore) GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error) {
	return us.q.GetEndpointHistory(ctx, params)
}

func (us UrlStore) GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error) {
	return us.q.GetNonExpiredEndpointsOfUser(ctx, userId)
}

func (us UrlStore) InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error) {
	params.Endpoint = strings.ToLower(params.Endpoint)
	return us.q.InsertFreeEndpoint(ctx, params)
}

func (us UrlStore) InsertEndpoint(ctx context.Context, params db.InsertEndpointParams) (db.Endpoint, error) {
	params.Endpoint = strings.ToLower(params.Endpoint)
	return us.q.InsertEndpoint(ctx, params)
}

// TODO: Move this
func (us UrlStore) CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error) {
	return us.q.CreateNewRequest(ctx, params)
}

func (us UrlStore) GetRequestById(ctx context.Context, reqId int64) (db.Request, error) {
	return us.q.GetRequestById(ctx, reqId)
}

func (us UrlStore) GetRequestByUUID(ctx context.Context, uuid string) (db.Request, error) {
	return us.q.GetRequestByUUID(ctx, uuid)
}

func (us UrlStore) ExpireRequests(ctx context.Context) error {
	return us.q.DeleteExpiredRequests(ctx)
}
