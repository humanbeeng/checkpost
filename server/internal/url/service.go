package url

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/humanbeeng/checkpost/server/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type UrlService struct {
	urlq  UrlQuerier
	userq user.UserQuerier
}

func NewUrlService(urlq UrlQuerier, userq user.UserQuerier) *UrlService {
	return &UrlService{
		urlq:  urlq,
		userq: userq,
	}
}

// TODO: Convert this into checkpost custom error
func NewInternalServerError() *UrlError {
	return &UrlError{
		Code:    http.StatusInternalServerError,
		Message: "Oops! Something went wrong :(",
	}
}

const (
	RandomUrlLength      int = 10
	NumUrlLimitPlanBasic int = 1
	DefaultExpiryHours   int = 4
)

func (s *UrlService) CreateUrl(c context.Context, username string, endpoint string) (db.Endpoint, *UrlError) {
	// Check subdomain length
	if len(endpoint) < 4 {
		return db.Endpoint{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Subdomain should be atleast 4 characters.",
		}
	}

	user, err := s.userq.GetUserFromUsername(c, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Endpoint{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("No user found with username: %s", username),
			}
		}
		return db.Endpoint{}, NewInternalServerError()

	}

	slog.Info("Create url request received", "endpoint", endpoint, "username", username)

	// Check if user has exceeded number of urls that can be generated
	urls, err := s.urlq.GetNonExpiredEndpointsOfUser(c, pgtype.Int8{Int64: user.ID, Valid: true})
	if err != nil {
		slog.Error("Unable to get non expired endpoints", "userId", user.ID, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	switch user.Plan {
	case db.PlanFree:
		{
			if len(urls) >= 1 {
				return db.Endpoint{}, &UrlError{
					Code:    http.StatusBadRequest,
					Message: "Cannot generate more that one url for your current plan. Consider upgrading to Pro.",
				}
			}
			return s.CreateFreeUrl(c, user.ID)
		}
	case db.PlanBasic, db.PlanPro:
		{
			url := fmt.Sprintf("https://%v.checkpost.io", endpoint)

			if user.Plan == db.PlanBasic && len(urls) >= NumUrlLimitPlanBasic {
				return db.Endpoint{}, &UrlError{
					Code:    http.StatusBadRequest,
					Message: "Cannot generate more than one url for your current plan. Consider upgrading to Pro.",
				}
			}

			if _, ok := core.ReservedDomains[endpoint]; ok {
				return db.Endpoint{}, &UrlError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("URL %s is reserved.", url),
				}
			}

			// Check if the requested endpoint already exists
			exists, err := s.urlq.CheckEndpointExists(c, endpoint)
			if err != nil {
				slog.Error("Unable to check if endpoint already exists", "err", err)
				return db.Endpoint{}, NewInternalServerError()
			}
			if exists {
				return db.Endpoint{}, &UrlError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("URL %s already exists", url),
				}
			}

			// endpoint is available
			slog.Info("Creating new pro endpoint", "endpoint", endpoint, "username", username)

			endpointRecord, err := s.urlq.InsertEndpoint(c, db.InsertEndpointParams{
				Endpoint: endpoint,
				UserID:   pgtype.Int8{Int64: user.ID, Valid: true},
				Plan:     user.Plan,

				// Never expires
				ExpiresAt: pgtype.Timestamp{
					Time:             time.Now().Add(time.Hour * 24),
					InfinityModifier: pgtype.Infinity,
					Valid:            true,
				},
			})
			if err != nil {
				slog.Error("Unable to insert new url into db", "endpoint", endpoint, "user", user.ID, "err", err)
				return db.Endpoint{}, NewInternalServerError()
			}

			endpointRecord.Endpoint = url
			slog.Info("Endpoint created and inserted into db", "endpoint", endpoint, "user", user.ID)

			return endpointRecord, nil
		}
	}

	slog.Error("Invalid user plan", "user", username, "plan", user.Plan)
	return db.Endpoint{}, &UrlError{Code: http.StatusBadRequest, Message: "Invalid user plan."}
}

func (s *UrlService) StoreRequestDetails(ctx context.Context, hookReq HookRequest) (db.Request, *UrlError) {

	endpoint := hookReq.Endpoint

	endpointRecord, err := s.urlq.GetEndpoint(ctx, endpoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Request{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("https://%s.checkpost.io is either not created or has expired.", endpoint),
			}
		}
		return db.Request{}, NewInternalServerError()
	}

	slog.Info("Storing request details", "endpoint", endpoint)

	userId := endpointRecord.UserID

	queryBytes, err := json.Marshal(hookReq.Query)
	if err != nil {
		slog.Error("Unable to marshal query params", "err", err)
		return db.Request{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Unable to parse query params.",
		}
	}

	headerBytes, err := json.Marshal(hookReq.Headers)
	if err != nil {
		slog.Error("Unable to marshal headers", "err", err)
		return db.Request{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Unable to parse headers",
		}
	}

	var expiresAt pgtype.Timestamp
	switch endpointRecord.Plan {
	case db.PlanGuest, db.PlanFree:
		{
			expiresAt = pgtype.Timestamp{
				Time:             time.Now().Add(time.Hour * time.Duration(DefaultExpiryHours)),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			}

		}
	case db.PlanPro, db.PlanBasic:
		{
			expiresAt = pgtype.Timestamp{
				Time:             time.Now(),
				InfinityModifier: pgtype.Infinity,
				Valid:            true,
			}
		}
	default:
		{
			return db.Request{}, &UrlError{
				Code:    http.StatusBadRequest,
				Message: "Invalid user plan",
			}
		}
	}

	requestRecord, err := s.urlq.CreateNewRequest(ctx, db.CreateNewRequestParams{
		UserID:     userId,
		EndpointID: endpointRecord.ID,
		Method:     db.HttpMethod(strings.ToLower(hookReq.Method)),
		Content:    pgtype.Text{String: hookReq.Content, Valid: true},
		Path:       hookReq.Path,

		// TODO: Fetch response from configured response
		ResponseCode: pgtype.Int4{Int32: http.StatusOK, Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     hookReq.SourceIp,

		// TODO: Add request body limiting
		ContentSize: int32(len(hookReq.Content)),
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		slog.Error("Unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return db.Request{}, NewInternalServerError()
	}

	slog.Info("Endpoint record created", "endpoint", endpoint, "userId", userId.Int64)

	return requestRecord, nil
}

func (s *UrlService) CreateGuestUrl(c context.Context) (db.Endpoint, *UrlError) {
	slog.Info("Creating random URL")

	randomEndpoint, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", RandomUrlLength)
	if err != nil {
		slog.Error("Unable to generate nano id", "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	randomUrl := fmt.Sprintf("https://%s.checkpost.io", randomEndpoint)

	// Inserting into db assuming that no endpoint with that random url existed. We can add
	// a check later on if needed.
	record, err := s.urlq.InsertGuestEndpoint(c, db.InsertGuestEndpointParams{
		Endpoint: randomEndpoint,
		ExpiresAt: pgtype.Timestamp{
			Time:             time.Now().Add(time.Hour * time.Duration(DefaultExpiryHours)),
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
	})
	if err != nil {
		slog.Error("Unable to insert endpoint", "endpoint", randomUrl, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	slog.Info("Random url generated", "url", randomUrl)
	return record, nil
}

func (s *UrlService) CreateFreeUrl(c context.Context, userId int64) (db.Endpoint, *UrlError) {
	slog.Info("Creating free url", "userId", userId)

	randomEndpoint, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", RandomUrlLength)
	if err != nil {
		slog.Error("Unable to generate nano id", "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	freeUrl := fmt.Sprintf("https://%s.checkpost.io", randomEndpoint)

	// Inserting into db assuming that no endpoint with that random url existed. We can add
	// a check later on if needed.
	endpointRecord, err := s.urlq.InsertFreeEndpoint(c, db.InsertFreeEndpointParams{
		Endpoint: randomEndpoint,
		UserID:   pgtype.Int8{Int64: userId, Valid: true},

		// TODO: Fetch expiry from config
		ExpiresAt: pgtype.Timestamp{
			Time:             time.Now().Add(time.Hour * 24),
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
	})
	if err != nil {
		slog.Error("Unable to insert free endpoint", "endpoint", freeUrl, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	slog.Info("Free url generated", "url", freeUrl)
	return endpointRecord, nil
}

func (s *UrlService) GetEndpointRequestHistory(c context.Context, endpoint string, limit int32, offset int32) ([]Request, *UrlError) {
	slog.Info("Fetch endpoint request history", "endpoint", endpoint)

	// TODO: Add a check to see if the user is authorized to access this endpoint history
	var reqHistory []Request

	reqs, err := s.urlq.GetEndpointHistory(c, db.GetEndpointHistoryParams{Endpoint: endpoint, Limit: limit, Offset: offset})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return reqHistory, nil
		}
		slog.Error("Unable to fetch endpoint request history", "endpoint", endpoint, "err", err)
		return nil, NewInternalServerError()
	}

	for _, req := range reqs {
		rh := Request{
			ID:           req.ID,
			Path:         req.Path,
			Content:      req.Content,
			Method:       req.Method,
			SourceIp:     req.SourceIp,
			ContentSize:  req.ContentSize,
			ResponseCode: req.ResponseCode,
			ExpiresAt:    req.ExpiresAt,
		}

		json.Unmarshal(req.Headers, &rh.Headers)
		json.Unmarshal(req.QueryParams, &rh.QueryParams)

		reqHistory = append(reqHistory, rh)
	}

	return reqHistory, nil
}

func (s *UrlService) GetRequestDetails(c context.Context, reqId int64) (Request, *UrlError) {
	slog.Info("Fetch request details", "reqId", reqId)

	reqRecord, err := s.urlq.GetRequestById(c, reqId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Request{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("No request found for request id: %v", reqId),
			}
		} else {
			slog.Error("Unable to fetch request details", "reqId", reqId, "err", err)
			return Request{}, NewInternalServerError()
		}
	}

	req := Request{
		ID:           reqRecord.ID,
		Path:         reqRecord.Path,
		Method:       reqRecord.Method,
		SourceIp:     reqRecord.SourceIp,
		Content:      reqRecord.Content,
		ContentSize:  reqRecord.ContentSize,
		ResponseCode: reqRecord.ResponseCode,
		CreatedAt:    reqRecord.CreatedAt,
		ExpiresAt:    reqRecord.ExpiresAt,
	}

	json.Unmarshal(reqRecord.Headers, &req.Headers)
	json.Unmarshal(reqRecord.QueryParams, &req.QueryParams)

	return req, nil
}

func (s *UrlService) GetEndpointStats(c context.Context, endpoint string) (EndpointStats, *UrlError) {
	slog.Info("Request endpoint stats", "endpoint", endpoint)

	endpointDetails, err := s.urlq.GetEndpoint(c, endpoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return EndpointStats{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Endpoint %v not found", endpoint),
			}
		}
		slog.Error("Unable to fetch endpoint details", "endpoint", endpoint, "err", err)
		return EndpointStats{}, NewInternalServerError()
	}

	stats, err := s.urlq.GetEndpointRequestCount(c, endpoint)
	if err != nil {
		slog.Error("Unable to fetch endpoint request count", "endpoint", endpoint, "err", err)
		return EndpointStats{}, NewInternalServerError()
	}

	return EndpointStats{
		SuccessCount: stats.SuccessCount,
		FailureCount: stats.FailureCount,
		TotalCount:   stats.TotalCount,
		ExpiresAt:    endpointDetails.ExpiresAt.Time.String(),
		Plan:         string(endpointDetails.Plan),
	}, nil
}

func (s *UrlService) CheckEndpointExists(c context.Context, endpoint string) (bool, *UrlError) {

	slog.Info("Checking if endpoint exists", "endpoint", endpoint)

	exists, err := s.urlq.CheckEndpointExists(c, endpoint)
	if err != nil {
		slog.Error("Unable to check if endpoint exists", "endpoint", endpoint, "err", err)
		return false, NewInternalServerError()
	}
	return exists, nil
}
