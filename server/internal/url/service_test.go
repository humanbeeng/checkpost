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

type UrlStoreTest struct{}
type UserStoreTest struct{}

type CtxKey string

const (
	NumUrls CtxKey = "num_expired"
)

const (
	UnknownUser string = "unknown_user"
	FreeUser    string = "free_user"
	ProUser     string = "pro_user"
	BasicUser   string = "basic_user"

	GuestEndpoint string = "abcd"
	FreeEndpoint  string = "freeendpoint"
	ProEndpoint   string = "proendpoint"
	BasicEndpoint string = "basicendpoint"
)

func (us UserStoreTest) GetUserFromUsername(ctx context.Context, username string) (db.User, error) {
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

var _ UrlQuerier = (*UrlStoreTest)(nil)

var userStore = UserStoreTest{}
var urlStore = UrlStoreTest{}

var service = UrlService{
	urlq:  urlStore,
	userq: userStore,
}

func (us UrlStoreTest) GetEndpointRequestCount(ctx context.Context, endpoint string) (db.GetEndpointRequestCountRow, error) {
	return db.GetEndpointRequestCountRow{}, nil
}

func (us UrlStoreTest) GetEndpoint(ctx context.Context, endpoint string) (db.Endpoint, error) {
	return db.Endpoint{}, nil
}

func (us UrlStoreTest) GetEndpointHistory(ctx context.Context, params db.GetEndpointHistoryParams) ([]db.GetEndpointHistoryRow, error) {
	return []db.GetEndpointHistoryRow{}, nil
}

func (us UrlStoreTest) GetNonExpiredEndpointsOfUser(ctx context.Context, userId pgtype.Int8) ([]db.Endpoint, error) {
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
func (us UrlStoreTest) InsertFreeEndpoint(ctx context.Context, params db.InsertFreeEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{
		Endpoint: params.Endpoint,
		ExpiresAt: pgtype.Timestamp{
			Time: time.Now().Add(time.Hour * time.Duration(DefaultExpiryHours)),
		},
	}, nil
}

func (us UrlStoreTest) InsertGuestEndpoint(ctx context.Context, params db.InsertGuestEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{
		Endpoint: GuestEndpoint,
		Plan:     db.PlanGuest,
	}, nil
}

func (us UrlStoreTest) InsertEndpoint(ctx context.Context, arg db.InsertEndpointParams) (db.Endpoint, error) {
	return db.Endpoint{Endpoint: arg.Endpoint}, nil
}

func (us UrlStoreTest) CheckEndpointExists(ctx context.Context, endpoint string) (bool, error) {
	return (endpoint == ProEndpoint || endpoint == GuestEndpoint || endpoint == BasicEndpoint || endpoint == FreeEndpoint), nil
}

// TODO: Remove this
func (us UrlStoreTest) CreateNewRequest(ctx context.Context, params db.CreateNewRequestParams) (db.Request, error) {
	return db.Request{}, nil
}

func (us UrlStoreTest) GetRequestById(ctx context.Context, reqId int64) (db.Request, error) {
	return db.Request{}, nil
}

func TestCheckEndpointExists(t *testing.T) {
	exists, err := service.CheckEndpointExists(context.Background(), ProEndpoint)
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
}

func TestCreateUrlForGuestUser(t *testing.T) {
	endpoint, err := service.CreateGuestUrl(context.TODO())
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint.Endpoint)
}

func TestCreateUrlWhenAlreadyExists(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), ProUser, ProEndpoint)
	assert.NotNil(t, err)
	assert.Empty(t, endpoint)
}

func TestCreateUrlForProUser(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), ProUser, "createurl")
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://createurl.checkpost.io", endpoint.Endpoint)
}

func TestCreateUrlForBasicUser(t *testing.T) {
	endpoint, err := service.CreateUrl(context.TODO(), BasicUser, "createurl")
	assert.Nil(t, err)
	assert.NotEmpty(t, endpoint)
	assert.Equal(t, "https://createurl.checkpost.io", endpoint.Endpoint)
}

func TestCreateUrlWhenFreeUserHasExistingEndpoint(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumUrls, 1)
	url, err := service.CreateUrl(ctx, FreeUser, FreeEndpoint)
	assert.Equal(t, err.Code, http.StatusBadRequest)
	assert.Empty(t, url)
}

func TestCreateUrlWhenBasicUserHasExistingEndpoint(t *testing.T) {
	ctx := context.WithValue(context.TODO(), NumUrls, 1)
	url, err := service.CreateUrl(ctx, BasicUser, FreeEndpoint)
	assert.Error(t, err)
	assert.Equal(t, err.Code, http.StatusBadRequest)
	assert.Empty(t, url)
}

func TestCreateUrlForReservedDomains(t *testing.T) {
	url, err := service.CreateUrl(context.TODO(), ProUser, "dash")
	fmt.Println(err)
	assert.Error(t, err)
	assert.Equal(t, http.StatusConflict, err.Code)
	assert.Empty(t, url)
}

func TestCreateUrlWhenEndpointLessThanFourChars(t *testing.T) {
	url, err := service.CreateUrl(context.TODO(), ProUser, "a")
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Empty(t, url)
}
