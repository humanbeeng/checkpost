package url

import (
	"context"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type UrlQuerier interface {
	GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error)
	GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error)
	GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error)
	GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error)

	InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error)
	InsertGuestEndpoint(ctx context.Context, params db.InsertGuestEndpointParams) (db.Endpoint, error)
	InsertEndpoint(ctx context.Context, params db.InsertEndpointParams) (db.Endpoint, error)

	CheckEndpointExists(ctx context.Context, endpoint string) (bool, error)

	// TODO: Remove this
	CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error)
	GetRequestById(ctx context.Context, reqId int64) (db.Request, error)
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

func (us UrlStore) GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error) {
	return us.q.GetEndpointHistory(ctx, params)
}

func (us UrlStore) GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error) {
	return us.q.GetNonExpiredEndpointsOfUser(ctx, userId)
}

func (us UrlStore) InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error) {
	return us.q.InsertFreeEndpoint(ctx, params)
}

func (us UrlStore) InsertGuestEndpoint(ctx context.Context, params db.InsertGuestEndpointParams) (db.Endpoint, error) {
	return us.q.InsertGuestEndpoint(ctx, params)
}

func (us UrlStore) InsertEndpoint(ctx context.Context, params db.InsertEndpointParams) (db.Endpoint, error) {
	return us.q.InsertEndpoint(ctx, params)
}

// TODO: Move this
func (us UrlStore) CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error) {
	return us.q.CreateNewRequest(ctx, params)
}

func (us UrlStore) GetRequestById(ctx context.Context, reqId int64) (db.Request, error) {
	return us.q.GetRequestById(ctx, reqId)
}
