package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	stdurl "net/url"
	"slices"
	"strings"
	"time"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/humanbeeng/checkpost/server/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type EndpointService struct {
	endpointq EndpointQuerier
	userq     user.UserQuerier
}

func NewEndpointService(endpointq EndpointQuerier, userq user.UserQuerier) *EndpointService {
	return &EndpointService{
		endpointq: endpointq,
		userq:     userq,
	}
}

// TODO: Convert this into checkpost custom error
func NewInternalServerError() *EndpointError {
	return &EndpointError{
		Code:    http.StatusInternalServerError,
		Message: "Oops! Something went wrong :(",
	}
}

const (
	RandomEndpointLength int = 10
	DefaultLimitNumUrl   int = 1
	DefaultExpiryHours   int = 6
)

func (s *EndpointService) CreateEndpoint(ctx context.Context, username string, subdomain string) (db.Endpoint, *EndpointError) {
	// Check endpoint length
	if len(subdomain) < 4 || len(subdomain) > 10 {
		return db.Endpoint{}, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: "Endpoint should be 4 to 10 characters.",
		}
	}
	subdomain = strings.ToLower(subdomain)

	endpoint := fmt.Sprintf("https://%v.checkpost.io", subdomain)

	_, err := stdurl.ParseRequestURI(endpoint)
	if err != nil {
		return db.Endpoint{}, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: "Invalid URL",
		}
	}

	// Check reserved endpoints
	if _, ok := core.ReservedSubdomains[subdomain]; ok {
		return db.Endpoint{}, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("URL %s is reserved.", endpoint),
		}
	}

	// Check if the requested endpoint already exists
	exists, err := s.endpointq.CheckEndpointExists(ctx, subdomain)
	if err != nil {
		slog.Error("unable to check if endpoint already exists", "endpoint", subdomain, "username", username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}
	if exists {
		slog.Info("Endpoint exists", "endpoint", endpoint)
		return db.Endpoint{}, &EndpointError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("Endpoint %s already exists", endpoint),
		}
	}

	user, err := s.userq.GetUserFromUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Endpoint{}, &EndpointError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("No user found with username: %s", username),
			}
		}
		slog.Error("unable to get user from username", "username", username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	// Check if user has exceeded number of urls that can be generated
	urls, err := s.endpointq.GetNonExpiredEndpointsOfUser(ctx, pgtype.Int8{Int64: user.ID, Valid: true})
	if err != nil {
		slog.Error("unable to get non expired endpoints", "username", user.Username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	if (user.Plan == db.PlanBasic || user.Plan == db.PlanFree) && len(urls) >= DefaultLimitNumUrl {
		return db.Endpoint{}, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: "Cannot generate more than one endpoint for your current plan. Consider upgrading to Pro.",
		}
	}

	// Check reserved companies. If found, check if the mail is from that organisation
	if _, ok := core.ReservedCompanies[subdomain]; ok {
		if !strings.Contains(strings.ToLower(user.Email), subdomain) || strings.Contains(strings.ToLower(user.Email), "@gmail.com") {
			return db.Endpoint{}, &EndpointError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("You cannot use this endpoint. Please try with mail issued by %s", subdomain),
			}
		}
	}

	slog.Info("Create endpoint request received", "endpoint", subdomain, "username", username, "plan", user.Plan)

	endpointRecord, err := s.endpointq.InsertEndpoint(ctx, db.InsertEndpointParams{
		Endpoint: subdomain,
		UserID:   pgtype.Int8{Int64: user.ID, Valid: true},
		Plan:     user.Plan,

		// Never expires
		ExpiresAt: pgtype.Timestamptz{
			Time:             time.Now().Add(time.Hour * 24),
			InfinityModifier: pgtype.Infinity,
			Valid:            true,
		},
	})
	if err != nil {
		slog.Error("unable to insert new endpoint into db", "endpoint", subdomain, "username", user.Username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	// Send complete endpoint as response
	endpointRecord.Endpoint = endpoint

	slog.Info("Endpoint created", "endpoint", endpoint, "username", user.Username, "plan", user.Plan)

	return endpointRecord, nil
}

func (s *EndpointService) GetUserEndpoints(ctx context.Context, userId int64) ([]Endpoint, *EndpointError) {
	endpointsRec, err := s.endpointq.GetUserEndpoints(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []Endpoint{}, nil
		}
		slog.Error("unable to fetch user endpoints", "userId", userId, "err", err)
		return []Endpoint{}, NewInternalServerError()
	}
	var endpoints []Endpoint

	for _, e := range endpointsRec {
		endpoints = append(endpoints, Endpoint{
			Endpoint:  e.Endpoint,
			ExpiresAt: e.ExpiresAt.Time,
			Plan:      string(e.Plan),
		})
	}
	return endpoints, nil
}

func (s *EndpointService) StoreRequestDetails(ctx context.Context, hookReq HookRequest) (db.Request, *EndpointError) {
	endpoint := hookReq.Endpoint

	endpointRecord, err := s.endpointq.GetEndpoint(ctx, endpoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Request{}, &EndpointError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("https://%s.checkpost.io is either not created or has expired.", endpoint),
			}
		}
		slog.Error("unable to get endpoint details", "endpoint", endpoint, "err", err)
		return db.Request{}, NewInternalServerError()
	}

	slog.InfoContext(ctx, "Storing request details", "endpoint", endpoint, "path", hookReq.Path)

	queryBytes, err := json.Marshal(hookReq.QueryParams)
	if err != nil {
		slog.Error("unable to marshal query params", "err", err)
		return db.Request{}, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: "unable to parse query params.",
		}
	}

	headerBytes, err := json.Marshal(hookReq.Headers)
	if err != nil {
		slog.Error("unable to marshal headers", "err", err)
		return db.Request{}, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: "unable to parse headers",
		}
	}

	var expiresAt pgtype.Timestamptz
	var content pgtype.Text
	var responseCode int
	switch endpointRecord.Plan {
	case db.PlanFree:
		{

			// Checking if request body exceeds 10KB
			if hookReq.ContentSize > 10_000 {
				slog.Warn("Received content that exceeds limit", "plan", endpointRecord.Plan, "received_size", hookReq.ContentSize, "limit", 10_000)
				content = pgtype.Text{Valid: true, String: ""}
				responseCode = http.StatusRequestEntityTooLarge
			} else {
				content = pgtype.Text{Valid: true, String: hookReq.Content}
				responseCode = http.StatusOK
			}

			expiresAt = pgtype.Timestamptz{
				// Use time.Duration for arithmetic operations on time.
				Time:             time.Now().Add(time.Hour * time.Duration(DefaultExpiryHours)),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			}
		}
	case db.PlanPro, db.PlanBasic:
		{
			// Checking if request body exceeds 512KB
			if hookReq.ContentSize > 512_000 {
				slog.Warn("Received content that exceeds limit", "plan", endpointRecord.Plan, "received_size", hookReq.ContentSize, "limit", 512_000)

				content = pgtype.Text{Valid: true, String: ""}
				responseCode = http.StatusRequestEntityTooLarge
			} else {
				content = pgtype.Text{Valid: true, String: hookReq.Content}
				responseCode = http.StatusOK
			}
		}
	default:
		{
			return db.Request{}, &EndpointError{
				Code:    http.StatusBadRequest,
				Message: "Invalid user plan",
			}
		}
	}

	userId := endpointRecord.UserID

	slog.Info("Request code", "code", responseCode)
	requestParams := db.CreateNewRequestParams{
		UserID:      userId,
		EndpointID:  endpointRecord.ID,
		Method:      db.HttpMethod(strings.ToLower(hookReq.Method)),
		Content:     content,
		ContentType: hookReq.ContentType,
		Path:        hookReq.Path,
		Uuid:        hookReq.UUID,

		// TODO: Fetch response from configured response
		ResponseCode: pgtype.Int4{Int32: int32(responseCode), Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     hookReq.SourceIp,

		ContentSize: int32(hookReq.ContentSize),
		ExpiresAt:   expiresAt,
	}

	if strings.Contains(hookReq.ContentType, string(MultipartForm)) || strings.Contains(hookReq.ContentType, string(FormUrlEncoded)) {
		formBytes, err := json.Marshal(hookReq.FormData)
		if err != nil {
			slog.Error("unable to marshal form data", "err", err)
			return db.Request{}, &EndpointError{
				Code:    http.StatusBadRequest,
				Message: "unable to parse form data",
			}
		}
		requestParams.FormData = formBytes
	}

	requestRecord, err := s.endpointq.CreateNewRequest(ctx, requestParams)
	if err != nil {
		slog.Error("unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return db.Request{}, NewInternalServerError()
	}

	slog.Info("Endpoint record created", "endpoint", endpoint, "userId", userId.Int64, "createdAt", requestRecord.CreatedAt)

	return requestRecord, nil
}

func (s *EndpointService) GetEndpointRequestHistory(ctx context.Context, endpoint string, userId int64, limit int32, offset int32) ([]HookRequest, *EndpointError) {
	slog.Info("Fetch endpoint request history", "endpoint", endpoint, "userId", userId)

	var reqHistory []HookRequest

	endpointsRec, err := s.endpointq.GetNonExpiredEndpointsOfUser(ctx, pgtype.Int8{Int64: userId, Valid: true})
	if err != nil {
		return nil, NewInternalServerError()
	}

	var endpoints []string

	for _, e := range endpointsRec {
		endpoints = append(endpoints, e.Endpoint)
	}

	if !slices.Contains(endpoints, endpoint) {
		return reqHistory, &EndpointError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	reqs, err := s.endpointq.GetEndpointHistory(ctx,
		db.GetEndpointHistoryParams{
			Endpoint: endpoint,
			UserID: pgtype.Int8{
				Int64: userId,
				Valid: true,
			},
			Limit:  limit,
			Offset: offset,
		})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return reqHistory, nil
		}
		slog.Error("unable to fetch endpoint request history", "endpoint", endpoint, "userId", userId, "err", err)
		return nil, NewInternalServerError()
	}

	for _, req := range reqs {
		rh := HookRequest{
			Endpoint:     endpoint,
			UUID:         req.Uuid,
			Path:         req.Path,
			Content:      req.Content.String,
			ContentType:  req.ContentType,
			Method:       string(req.Method),
			SourceIp:     req.SourceIp,
			ContentSize:  req.ContentSize,
			ResponseCode: req.ResponseCode.Int32,
			CreatedAt:    req.CreatedAt.Time,
			ExpiresAt:    req.ExpiresAt.Time,
		}

		json.Unmarshal(req.Headers, &rh.Headers)
		json.Unmarshal(req.FormData, &rh.FormData)
		json.Unmarshal(req.QueryParams, &rh.QueryParams)

		reqHistory = append(reqHistory, rh)
	}

	return reqHistory, nil
}

func (s *EndpointService) GetRequestDetails(ctx context.Context, reqId int64) (HookRequest, *EndpointError) {
	slog.Info("Request to fetch request details", "reqId", reqId)

	reqRecord, err := s.endpointq.GetRequestById(ctx, reqId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return HookRequest{}, &EndpointError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("No request found for request id: %v", reqId),
			}
		} else {
			slog.Error("unable to fetch request details", "reqId", reqId, "err", err)
			return HookRequest{}, NewInternalServerError()
		}
	}

	req := HookRequest{
		UUID:         reqRecord.Uuid,
		Path:         reqRecord.Path,
		Method:       string(reqRecord.Method),
		SourceIp:     reqRecord.SourceIp,
		Content:      reqRecord.Content.String,
		ContentType:  reqRecord.ContentType,
		ContentSize:  reqRecord.ContentSize,
		ResponseCode: reqRecord.ResponseCode.Int32,
		CreatedAt:    reqRecord.CreatedAt.Time,
		ExpiresAt:    reqRecord.ExpiresAt.Time,
	}

	json.Unmarshal(reqRecord.Headers, &req.Headers)
	json.Unmarshal(reqRecord.FormData, &req.FormData)
	json.Unmarshal(reqRecord.QueryParams, &req.QueryParams)

	return req, nil
}

func (s *EndpointService) GetEndpointStats(ctx context.Context, endpoint string) (EndpointStats, *EndpointError) {
	endpoint = strings.ToLower(endpoint)
	slog.Info("Request endpoint stats", "endpoint", endpoint)

	endpointDetails, err := s.endpointq.GetEndpoint(ctx, endpoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return EndpointStats{}, &EndpointError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Endpoint %v not found", endpoint),
			}
		}
		slog.Error("unable to fetch endpoint details", "endpoint", endpoint, "err", err)
		return EndpointStats{}, NewInternalServerError()
	}

	stats, err := s.endpointq.GetEndpointRequestCount(ctx, endpoint)
	if err != nil {
		slog.Error("unable to fetch endpoint request count", "endpoint", endpoint, "err", err)
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

type EndpointExists string

const (
	Available        EndpointExists = "Its available. Sign up and make it yours"
	Taken            EndpointExists = "That endpoint is already taken. Try something else?"
	ReservedCompany  EndpointExists = "Endpoint is reserved. But, you can go ahead if you're using mail issued from that organisation."
	ReservedEndpoint EndpointExists = "Endpoint is reserved."
	BadEndpoint      EndpointExists = "Bad endpoint."
	Error            EndpointExists = "Something went wrong."
)

func (s *EndpointService) CheckEndpointExists(ctx context.Context, subdomain string) (EndpointExists, *EndpointError) {
	subdomain = strings.ToLower(subdomain)

	slog.InfoContext(ctx, "Checking if endpoint exists", "endpoint", subdomain)

	if len(subdomain) < 4 || len(subdomain) > 10 {
		return BadEndpoint, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: "Subdomain should be 4 to 10 characters.",
		}
	}

	if _, ok := core.ReservedSubdomains[subdomain]; ok {
		slog.Info("Subdomain is reserved", "subdomain", subdomain)
		return ReservedEndpoint, &EndpointError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Subdomain %s is reserved.", subdomain),
		}
	}

	// Check reserved companies.
	if _, ok := core.ReservedCompanies[subdomain]; ok {
		slog.Info("Subdomain is reserved company", "subdomain", subdomain)
		return ReservedCompany, nil
	}

	exists, err := s.endpointq.CheckEndpointExists(ctx, subdomain)
	if err != nil {
		slog.Error("unable to check if subdomain exists", "subdomain", subdomain, "err", err)
		return Error, NewInternalServerError()
	}

	if exists {
		return Taken, nil
	} else {
		slog.Info("Subdomain is taken", "subdomain", subdomain)
	}

	slog.Info("Subdomain available", "subdomain", subdomain)
	return Available, nil
}

func (s *EndpointService) GetRequestByUUID(ctx context.Context, uuid string) (HookRequest, *EndpointError) {
	slog.Info("Request to fetch request details by uuid", "uuid", uuid)
	reqRecord, err := s.endpointq.GetRequestByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return HookRequest{}, &EndpointError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("No request found for uuid: %v", uuid),
			}
		} else {
			slog.Error("unable to fetch request details", "uuid", uuid, "err", err)
			return HookRequest{}, NewInternalServerError()
		}
	}

	req := HookRequest{
		UUID:         reqRecord.Uuid,
		Path:         reqRecord.Path,
		Method:       string(reqRecord.Method),
		SourceIp:     reqRecord.SourceIp,
		Content:      reqRecord.Content.String,
		ContentSize:  reqRecord.ContentSize,
		ResponseCode: reqRecord.ResponseCode.Int32,
		CreatedAt:    reqRecord.CreatedAt.Time,
		ExpiresAt:    reqRecord.ExpiresAt.Time,
	}

	json.Unmarshal(reqRecord.Headers, &req.Headers)
	json.Unmarshal(reqRecord.QueryParams, &req.QueryParams)

	return req, nil
}

func (s *EndpointService) ExpireRequests(ctx context.Context) error {
	slog.Info("Deleting expired requests", "date", time.Now().Local().String())
	err := s.endpointq.ExpireRequests(ctx)
	if err != nil {
		slog.Error("unable to delete expired requests", "date", time.Now().Local().String(), "err", err)
		return err
	}

	return nil
}
