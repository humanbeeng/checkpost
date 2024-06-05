package endpoint

import (
	"context"
	"strings"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type EndpointQuerier interface {
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

type EndpointStore struct {
	q db.Querier
}

func NewEndpointStore(q db.Querier) *EndpointStore {
	return &EndpointStore{
		q: q,
	}
}

func (us EndpointStore) GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error) {
	return us.q.GetEndpointRequestCount(ctx, endpoint)
}

func (us EndpointStore) CheckEndpointExists(ctx context.Context, endpoint string) (bool, error) {
	return us.q.CheckEndpointExists(ctx, endpoint)
}

func (us EndpointStore) GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error) {
	return us.q.GetEndpointDetails(ctx, endpoint)
}

func (us EndpointStore) GetUserEndpoints(ctx context.Context, userId int64) ([]db.Endpoint, error) {
	return us.q.GetUserEndpoints(ctx, pgtype.Int8{Int64: userId, Valid: true})
}

func (us EndpointStore) GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error) {
	return us.q.GetEndpointHistory(ctx, params)
}

func (us EndpointStore) GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error) {
	return us.q.GetNonExpiredEndpointsOfUser(ctx, userId)
}

func (us EndpointStore) InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error) {
	params.Endpoint = strings.ToLower(params.Endpoint)
	return us.q.InsertFreeEndpoint(ctx, params)
}

func (us EndpointStore) InsertEndpoint(ctx context.Context, params db.InsertEndpointParams) (db.Endpoint, error) {
	params.Endpoint = strings.ToLower(params.Endpoint)
	return us.q.InsertEndpoint(ctx, params)
}

// TODO: Move this
func (us EndpointStore) CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error) {
	return us.q.CreateNewRequest(ctx, params)
}

func (us EndpointStore) GetRequestById(ctx context.Context, reqId int64) (db.Request, error) {
	return us.q.GetRequestById(ctx, reqId)
}

func (us EndpointStore) GetRequestByUUID(ctx context.Context, uuid string) (db.Request, error) {
	return us.q.GetRequestByUUID(ctx, uuid)
}

func (us EndpointStore) ExpireRequests(ctx context.Context) error {
	return us.q.DeleteExpiredRequests(ctx)
}
