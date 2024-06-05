package endpoint

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

type MockEndpointStore struct{}
type MockUserStore struct{}

type CtxKey string

const (
	NumEndpoints CtxKey = "num_endpoints"
)

const (
	UnknownUser string = "unknown_user"
	FreeUser    string = "free_user"
	ProUser     string = "pro_user"
	BasicUser   string = "basic_user"

	FreeEndpoint     string = "free-url"
	ProEndpoint      string = "pro-url"
	BasicEndpoint    string = "basic-url"
	UnknownEndpoint  string = "unknown-url"
	ExistingEndpoint string = "nonexist"
)

func (es MockUserStore) GetUserFromUsername(ctx context.Context, username string) (db.User, error) {
	if username == UnknownUser {
		return db.User{}, pgx.ErrNoRows
	} else if username == FreeUser {
		return db.User{
			Username: FreeUser,
			Plan:     db.PlanFree,
			Email:    "freeuser@checkpost.io",
		}, nil
	} else if username == ProUser {
		return db.User{
			Username: username,
			Plan:     db.PlanPro,
			Email:    "prouser@checkpost.io",
		}, nil
	} else if username == BasicUser {
		return db.User{
			Username: username,
			Plan:     db.PlanBasic,
			Email:    "basicuser@checkpost.io",
		}, nil
	}
	return db.User{}, fmt.Errorf("not found")
}

var _ EndpointQuerier = (*MockEndpointStore)(nil)

var userStore = MockUserStore{}
var endpointStore = MockEndpointStore{}

var service = EndpointService{
	endpointq: endpointStore,
	userq:     userStore,
}

func (es MockEndpointStore) GetUserEndpoints(ctx context.Context, userId int64) ([]db.Endpoint, error) {
	return []db.Endpoint{}, nil
}

func (es MockEndpointStore) GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error) {
	if endpoint == BasicEndpoint || endpoint == ProEndpoint || endpoint == FreeEndpoint {
		return db.GetEndpointRequestCountRow{
			SuccessCount: 100,
			FailureCount: 100,
			TotalCount:   200,
		}, nil
	}
	return db.GetEndpointRequestCountRow{}, &EndpointError{
		Code:    http.StatusNotFound,
		Message: "not found",
	}
}

func (es MockEndpointStore) GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error) {
	if endpoint != UnknownEndpoint {
		return db.Endpoint{
			Endpoint: endpoint,
			ExpiresAt: pgtype.Timestamptz{
				Time:             time.Now().Add(time.Hour),
				Valid:            true,
				InfinityModifier: pgtype.Infinity,
			},
			Plan:   db.PlanPro,
			UserID: pgtype.Int8{Int64: 1, Valid: true},
		}, nil
	}
	return db.Endpoint{}, pgx.ErrNoRows
}

func (es MockEndpointStore) ExpireRequests(ctx context.Context) error {
	return nil
}

func (es MockEndpointStore) GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error) {
	return []db.GetEndpointHistoryRow{}, nil
}

func (es MockEndpointStore) GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error) {
	numEndpoints := ctx.Value(NumEndpoints)
	if numEndpoints == nil {
		return []db.Endpoint{}, nil
	}
	var endpoints []db.Endpoint
	for range numEndpoints.(int) {
		endpoints = append(endpoints, db.Endpoint{
			Endpoint: FreeEndpoint,
		})
	}
	return endpoints, nil
}

// TODO: Rename this
func (es MockEndpointStore) InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{
		Endpoint: params.Endpoint,
		ExpiresAt: pgtype.Timestamptz{
			Time: time.Now().Add(time.Hour * time.Duration(DefaultExpiryHours)),
		},
	}, nil
}

func (es MockEndpointStore) InsertEndpoint(ctx context.Context, arg db.InsertEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{Endpoint: arg.Endpoint}, nil
}

func (es MockEndpointStore) CheckEndpointExists(ctx context.Context, endpoint string) (bool, error) {
	return endpoint == ExistingEndpoint, nil
}

// TODO: Move below mocks to request tests
func (es MockEndpointStore) CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error) {
	return db.Request{
		Method:       params.Method,
		QueryParams:  params.QueryParams,
		Headers:      params.Headers,
		Content:      params.Content,
		ContentSize:  params.ContentSize,
		Path:         params.Path,
		ResponseCode: params.ResponseCode,
		SourceIp:     params.SourceIp,
	}, nil
}

func (es MockEndpointStore) GetRequestById(ctx context.Context, reqId int64) (db.Request, error) {
	return db.Request{}, nil
}
func (es MockEndpointStore) GetRequestByUUID(ctx context.Context, uuid string) (db.Request, error) {
	return db.Request{}, nil
}

func TestCheckEndpointExists(t *testing.T) {
	exists, err := service.CheckEndpointExists(context.Background(), ExistingEndpoint)
	assert.Nil(t, err)
	assert.Equal(t, Taken, exists)
}

func TestCreateEndpointForFreeUserWhenUserNotFound(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), UnknownUser, FreeEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, err, &EndpointError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("No user found with username: %s", "unknown_user"),
	})
	assert.Empty(t, endpoint)
}

func TestCreateEndpointForFreeUser(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumEndpoints, 0)
	endpoint, err := service.CreateEndpoint(ctx, FreeUser, FreeEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://free-url.checkpost.io", endpoint.Endpoint)
}

func TestCreateEndpointWhenAlreadyExists(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), ProUser, ExistingEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Code)
	assert.Empty(t, endpoint)
}

func TestCreateEndpointForProUser(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), ProUser, ProEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://pro-url.checkpost.io", endpoint.Endpoint)
}

func TestCreateEndpointForBasicUser(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), BasicUser, BasicEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://basic-url.checkpost.io", endpoint.Endpoint)
}

func TestCreateEndpointWhenFreeUserHasExistingEndpoint(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumEndpoints, 1)
	endpoint, err := service.CreateEndpoint(ctx, FreeUser, FreeEndpoint)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, endpoint)
}

func TestCreateEndpointWhenBasicUserHasExistingEndpoint(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumEndpoints, 1)
	endpoint, err := service.CreateEndpoint(ctx, BasicUser, FreeEndpoint)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, endpoint)
}

func TestCreateEndpointForReservedDomains(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), ProUser, "dash")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, endpoint)
}

func TestCreateEndpointForReservedCompany(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), BasicUser, "google")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, endpoint)
}

func TestCreateEndpointForReservedCompanyWhenUserFromSameOrg(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), BasicUser, "checkpost")
	assert.Nil(t, err)
	assert.Equal(t, endpoint.Endpoint, "https://checkpost.checkpost.io")
}

func TestCreateEndpointWhenEndpointLessThanFourChars(t *testing.T) {
	endpoint, err := service.CreateEndpoint(context.TODO(), ProUser, "a")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, endpoint)
}

func TestGetBasicEndpointStats(t *testing.T) {
	stats, err := service.GetEndpointStats(context.TODO(), BasicEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, stats)
}

func TestGetEndpointStatsUnknownEndpoint(t *testing.T) {
	stats, err := service.GetEndpointStats(context.TODO(), UnknownEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Empty(t, stats)
}

func TestStoreRequestDetails(t *testing.T) {
	hookReq := HookRequest{
		Endpoint: FreeEndpoint,
		Path:     "/",
		Method:   string(db.HttpMethodPost),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		QueryParams: map[string]string{
			"hello": "there",
		},
		SourceIp:     "17.1.1.1",
		Content:      "{\"message\":\"hello world\"}",
		ContentSize:  25,
		ResponseCode: 200,
	}
	req, err := service.StoreRequestDetails(context.TODO(), hookReq)
	assert.Nil(t, err)
	assert.NotEmpty(t, req)

	assert.Equal(t, pgtype.Text{String: hookReq.Content, Valid: true}, req.Content)
	assert.Equal(t, hookReq.Path, req.Path)
	assert.Equal(t, int32(hookReq.ContentSize), req.ContentSize)
	assert.Equal(t, hookReq.Method, string(req.Method))
}

func TestStoreRequestDetailsWhenEndpointNotFound(t *testing.T) {
	hookReq := HookRequest{
		Endpoint: UnknownEndpoint,
		Path:     "/",
		Method:   string(db.HttpMethodPost),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		QueryParams: map[string]string{
			"hello": "there",
		},
		SourceIp:     "17.1.1.1",
		Content:      "{\"message\":\"hello world\"}",
		ContentSize:  25,
		ResponseCode: 200,
	}
	req, err := service.StoreRequestDetails(context.TODO(), hookReq)
	assert.NotNil(t, err)
	assert.Equal(t, err.Code, http.StatusNotFound)
	assert.Empty(t, req)
}
