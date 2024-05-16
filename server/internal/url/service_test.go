package url

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

type MockUrlStore struct{}
type MockUserStore struct{}

type CtxKey string

const (
	NumUrls CtxKey = "num_expired"
)

const (
	UnknownUser string = "unknown_user"
	FreeUser    string = "free_user"
	ProUser     string = "pro_user"
	BasicUser   string = "basic_user"

	FreeEndpoint     string = "freeurl"
	ProEndpoint      string = "prourl"
	BasicEndpoint    string = "basicurl"
	UnknownEndpoint  string = "unknownurl"
	ExistingEndpoint string = "nonexist"
)

func (us MockUserStore) GetUserFromUsername(ctx context.Context, username string) (db.User, error) {
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

var _ UrlQuerier = (*MockUrlStore)(nil)

var userStore = MockUserStore{}
var urlStore = MockUrlStore{}

var service = UrlService{
	urlq:  urlStore,
	userq: userStore,
}

func (us MockUrlStore) GetUserEndpoints(ctx context.Context, userId int64) ([]db.Endpoint, error) {
	return []db.Endpoint{}, nil
}

func (us MockUrlStore) GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error) {
	if endpoint == BasicEndpoint || endpoint == ProEndpoint || endpoint == FreeEndpoint {
		return db.GetEndpointRequestCountRow{
			SuccessCount: 100,
			FailureCount: 100,
			TotalCount:   200,
		}, nil
	}
	return db.GetEndpointRequestCountRow{}, &UrlError{
		Code:    http.StatusNotFound,
		Message: "not found",
	}
}

func (us MockUrlStore) GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error) {
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

func (us MockUrlStore) ExpireRequests(ctx context.Context) error {
	return nil
}

func (us MockUrlStore) GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error) {
	return []db.GetEndpointHistoryRow{}, nil
}

func (us MockUrlStore) GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error) {
	numUrls := ctx.Value(NumUrls)
	if numUrls == nil {
		return []db.Endpoint{}, nil
	}
	var urls []db.Endpoint
	for range numUrls.(int) {
		urls = append(urls, db.Endpoint{
			Endpoint: FreeEndpoint,
		})
	}
	return urls, nil
}

// TODO: Rename this
func (us MockUrlStore) InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{
		Endpoint: params.Endpoint,
		ExpiresAt: pgtype.Timestamptz{
			Time: time.Now().Add(time.Hour * time.Duration(DefaultExpiryHours)),
		},
	}, nil
}

func (us MockUrlStore) InsertEndpoint(ctx context.Context, arg db.InsertEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{Endpoint: arg.Endpoint}, nil
}

func (us MockUrlStore) CheckEndpointExists(ctx context.Context, endpoint string) (bool, error) {
	return endpoint == ExistingEndpoint, nil
}

// TODO: Move below mocks to request tests
func (us MockUrlStore) CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error) {
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

func (us MockUrlStore) GetRequestById(ctx context.Context, reqId int64) (db.Request, error) {
	return db.Request{}, nil
}
func (us MockUrlStore) GetRequestByUUID(ctx context.Context, uuid string) (db.Request, error) {
	return db.Request{}, nil
}

func TestCheckEndpointExists(t *testing.T) {
	exists, err := service.CheckEndpointExists(context.Background(), ExistingEndpoint)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestCreateUrlForFreeUserWhenUserNotFound(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), UnknownUser, FreeEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, err, &UrlError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("No user found with username: %s", "unknown_user"),
	})
	assert.Empty(t, endpoint)
}

func TestCreateUrlForFreeUser(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumUrls, 0)
	endpoint, err := service.CreateUrl(ctx, FreeUser, FreeEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://freeurl.checkpost.io", endpoint.Endpoint)
}

func TestCreateUrlWhenAlreadyExists(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), ProUser, ExistingEndpoint)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, err.Code)
	assert.Empty(t, endpoint)
}

func TestCreateUrlForProUser(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), ProUser, ProEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://prourl.checkpost.io", endpoint.Endpoint)
}

func TestCreateUrlForBasicUser(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), BasicUser, BasicEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://basicurl.checkpost.io", endpoint.Endpoint)
}

func TestCreateUrlWhenFreeUserHasExistingEndpoint(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumUrls, 1)
	url, err := service.CreateUrl(ctx, FreeUser, FreeEndpoint)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, url)
}

func TestCreateUrlWhenBasicUserHasExistingEndpoint(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumUrls, 1)
	url, err := service.CreateUrl(ctx, BasicUser, FreeEndpoint)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, url)
}

func TestCreateUrlForReservedDomains(t *testing.T) {
	url, err := service.CreateUrl(context.TODO(), ProUser, "dash")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, url)
}

func TestCreateUrlForReservedCompany(t *testing.T) {
	url, err := service.CreateUrl(context.TODO(), BasicUser, "google")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, url)
}

func TestCreateUrlForReservedCompanyWhenUserFromSameOrg(t *testing.T) {
	url, err := service.CreateUrl(context.TODO(), BasicUser, "checkpost")
	assert.Nil(t, err)
	assert.Equal(t, url.Endpoint, "https://checkpost.checkpost.io")
}

func TestCreateUrlWhenEndpointLessThanFourChars(t *testing.T) {
	url, err := service.CreateUrl(context.TODO(), ProUser, "a")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, url)
}

func TestGetBasicEndpointStats(t *testing.T) {
	stats, err := service.GetEndpointStats(context.TODO(), BasicEndpoint)
	assert.Nil(t, err)
	assert.NotEmpty(t, stats)
}

func TestGetEndpointStatsUnknownUrl(t *testing.T) {
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
		Query: map[string]string{
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
		Query: map[string]string{
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
