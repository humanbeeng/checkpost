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

	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/checkpost/server/config"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type UrlService struct {
	q      db.Querier
	config *config.AppConfig
}

func NewUrlService(q db.Querier, config *config.AppConfig) *UrlService {
	return &UrlService{q: q, config: config}
}

// TODO: Convert this into checkpost custom error
func NewInternalServerError() *UrlError {
	return &UrlError{
		Code:    http.StatusInternalServerError,
		Message: "Oops! Something went wrong :(",
	}
}

const (
	RandomUrlLength          int = 10
	NumUrlLimitPlanNoBrainer int = 1
)

func (s *UrlService) CreateUrl(c context.Context, username string, endpoint string) (string, *UrlError) {
	// Check subdomain length
	if len(endpoint) < 4 {
		return "", &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Subdomain should be atleast 4 characters.",
		}
	}

	user, err := s.q.GetUserFromUsername(c, username)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", &UrlError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("No user found with username: %s", username),
		}
	}
	// Check if user has exceeded number of urls that can be generated
	urls, err := s.q.GetNonExpiredEndpointsOfUser(c, pgtype.Int8{Int64: user.ID, Valid: true})
	if err != nil {
		slog.Error("Unable to get non expired endpoints", "userId", user.ID, "err", err)
		return "", NewInternalServerError()
	}

	switch user.Plan {
	case db.PlanFree:
		{
			if len(urls) >= 1 {
				return "", &UrlError{
					Code:    http.StatusBadRequest,
					Message: "Cannot generate more that one url for your current plan. Consider upgrading to Pro.",
				}
			}
			return s.CreateFreeUrl(c, user.ID)
		}
	case db.PlanNoBrainer, db.PlanPro:
		{
			url := fmt.Sprintf("https://%v.checkpost.io", endpoint)

			if user.Plan == db.PlanNoBrainer && len(urls) >= NumUrlLimitPlanNoBrainer {
				return "", &UrlError{
					Code:    http.StatusBadRequest,
					Message: "Cannot generate more than one url for your current plan. Consider upgrading to Pro.",
				}
			}

			if _, ok := core.ReservedDomains[endpoint]; ok {
				return "", &UrlError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("URL %s is reserved.", url),
				}
			}

			// Check if the requested endpoint already exists
			exists, err := s.q.CheckEndpointExists(c, endpoint)
			if err != nil {
				slog.Error("Unable to check if endpoint already exists", "err", err)
				return "", NewInternalServerError()
			}
			if exists {
				return "", &UrlError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("URL %s already exists", url),
				}
			}

			// endpoint is available
			slog.Info("Creating new pro endpoint", "endpoint", endpoint, "username", username)

			_, err = s.q.CreateNewEndpoint(c, db.CreateNewEndpointParams{
				Endpoint: endpoint,
				UserID:   pgtype.Int8{Int64: user.ID, Valid: true},
				Plan:     user.Plan,
				ExpiresAt: pgtype.Timestamp{
					// TODO: Change this
					Time:             time.Now().Add(time.Hour * 24),
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
			})
			if err != nil {
				slog.Error("Unable to insert new url into db", "endpoint", endpoint, "user", user.ID, "err", err)
				return "", NewInternalServerError()
			}

			slog.Info("Endpoint created and inserted into db", "endpoint", endpoint, "user", user.ID)

			return url, nil
		}
	}

	slog.Error("Invalid user plan", "user", username, "plan", user.Plan)
	return "", &UrlError{Code: http.StatusBadRequest, Message: "Invalid user plan."}
}

func (s *UrlService) StoreRequestDetails(c *fiber.Ctx) (HookRequest, *UrlError) {

	endpoint := c.Params("endpoint", "")
	if endpoint == "" {
		return HookRequest{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Empty endpoint",
		}
	}

	endpointRecord, err := s.q.GetEndpoint(c.Context(), endpoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return HookRequest{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("https://%s.checkpost.io is either not created or has expired.", endpoint),
			}
		}
		return HookRequest{}, NewInternalServerError()
	}

	slog.Info("Storing request details", "endpoint", endpoint)

	userId := endpointRecord.UserID

	var req any
	_ = c.BodyParser(&req)

	strBytes, _ := json.Marshal(req)
	body := string(strBytes)
	// Note: key is string and value is []string
	headers := c.GetReqHeaders()
	ip := c.IP()
	path := c.Params("path", "/")

	method := c.Method()
	query := c.Queries()

	queryBytes, err := json.Marshal(query)
	if err != nil {
		slog.Error("Unable to marshal query params", "err", err)
		return HookRequest{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Unable to parse query params.",
		}
	}

	headerBytes, err := json.Marshal(headers)
	if err != nil {
		slog.Error("Unable to marshal headers", "err", err)
		return HookRequest{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Unable to parse headers.",
		}
	}

	cr, err := s.q.CreateNewRequest(c.Context(), db.CreateNewRequestParams{
		UserID:     userId,
		EndpointID: endpointRecord.ID,
		Method:     db.HttpMethod(strings.ToLower(method)),
		Content:    pgtype.Text{String: body, Valid: true},
		Path:       path,

		// TODO: Fetch response from configured response
		ResponseCode: pgtype.Int4{Int32: http.StatusOK, Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     ip,

		// TODO: Add request body limiting
		ContentSize: int32(len(body)),
	})
	if err != nil {
		slog.Error("Unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return HookRequest{}, NewInternalServerError()
	}

	hookReq := HookRequest{
		Endpoint:     endpoint,
		Path:         path,
		Headers:      headers,
		Query:        query,
		SourceIp:     ip,
		Content:      body,
		ContentSize:  len(body),
		ResponseCode: http.StatusOK,
		CreatedAt:    cr.CreatedAt.Time,
	}

	slog.Info("Endpoint record created", "endpoint", endpoint, "userId", userId.Int64)

	return hookReq, nil
}

func (s *UrlService) CreateGuestUrl(c context.Context, user *db.User) (string, string, *UrlError) {
	slog.Info("Creating random URL")

	randomEndpoint, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", RandomUrlLength)
	if err != nil {
		slog.Error("Unable to generate nano id", "err", err)
		return "", "", NewInternalServerError()
	}

	randomUrl := fmt.Sprintf("https://%s.checkpost.io", randomEndpoint)

	// Inserting into db assuming that no endpoint with that random url existed. We can add
	// a check later on if needed.
	if _, err := s.q.CreateNewGuestEndpoint(c, db.CreateNewGuestEndpointParams{
		Endpoint: randomEndpoint,
		// TODO: Fetch expiry from config
		ExpiresAt: pgtype.Timestamp{
			Time:             time.Now().Add(time.Hour * 24),
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
	}); err != nil {
		slog.Error("Unable to insert endpoint", "endpoint", randomUrl, "err", err)
		return "", "", NewInternalServerError()
	}

	slog.Info("Random url generated", "url", randomUrl)
	return randomUrl, string(db.PlanGuest), nil
}

func (s *UrlService) CreateFreeUrl(c context.Context, userId int64) (string, *UrlError) {
	slog.Info("Creating free url", "userId", userId)

	randomEndpoint, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", RandomUrlLength)
	if err != nil {
		slog.Error("Unable to generate nano id", "err", err)
		return "", NewInternalServerError()
	}

	freeUrl := fmt.Sprintf("https://%s.checkpost.io", randomEndpoint)

	// Inserting into db assuming that no endpoint with that random url existed. We can add
	// a check later on if needed.
	if _, err := s.q.CreateNewFreeUrl(c, db.CreateNewFreeUrlParams{
		Endpoint: randomEndpoint,
		UserID:   pgtype.Int8{Int64: userId, Valid: true},

		// TODO: Fetch expiry from config
		ExpiresAt: pgtype.Timestamp{
			Time:             time.Now().Add(time.Hour * 24),
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
	}); err != nil {
		slog.Error("Unable to insert free endpoint", "endpoint", freeUrl, "err", err)
		return "", NewInternalServerError()
	}

	slog.Info("Free url generated", "url", freeUrl)
	return freeUrl, nil
}

func (s *UrlService) GetEndpointRequestHistory(c context.Context, endpoint string, limit int32, offset int32) ([]Request, *UrlError) {
	slog.Info("Fetch endpoint request history", "endpoint", endpoint)

	// TODO: Add a check to see if the user is authorized to access this endpoint history
	var reqHistory []Request

	reqs, err := s.q.GetEndpointHistory(c, db.GetEndpointHistoryParams{Endpoint: endpoint, Limit: limit, Offset: offset})
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
		}

		json.Unmarshal(req.Headers, &rh.Headers)
		json.Unmarshal(req.QueryParams, &rh.QueryParams)

		reqHistory = append(reqHistory, rh)
	}

	return reqHistory, nil
}

func (s *UrlService) GetRequestDetails(c context.Context, reqId int64) (Request, *UrlError) {
	slog.Info("Fetch request details", "reqId", reqId)

	reqRecord, err := s.q.GetRequestById(c, reqId)
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
	}

	json.Unmarshal(reqRecord.Headers, &req.Headers)
	json.Unmarshal(reqRecord.QueryParams, &req.QueryParams)

	return req, nil
}

func (s *UrlService) GetEndpointStats(c context.Context, endpoint string) (EndpointStats, *UrlError) {
	slog.Info("Request endpoint stats", "endpoint", endpoint)

	endpointDetails, err := s.q.GetEndpoint(c, endpoint)
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

	stats, err := s.q.GetEndpointRequestCount(c, endpoint)
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

	exists, err := s.q.CheckEndpointExists(c, endpoint)
	if err != nil {
		slog.Error("Unable to check if endpoint exists", "endpoint", endpoint, "err", err)
		return false, NewInternalServerError()
	}
	return exists, nil
}
