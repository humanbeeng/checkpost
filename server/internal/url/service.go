package url

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	stdurl "net/url"
	"strings"
	"time"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/humanbeeng/checkpost/server/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	RandomUrlLength    int = 10
	DefaultLimitNumUrl int = 1
	DefaultExpiryHours int = 4
)

func (s *UrlService) CreateUrl(c context.Context, username string, endpoint string) (db.Endpoint, *UrlError) {
	// Check subdomain length
	if len(endpoint) < 4 || len(endpoint) > 10 {
		return db.Endpoint{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Subdomain should be 4 to 10 characters.",
		}
	}
	endpoint = strings.ToLower(endpoint)

	url := fmt.Sprintf("https://%v.checkpost.io", endpoint)

	_, err := stdurl.ParseRequestURI(url)
	if err != nil {
		return db.Endpoint{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Invalid URL",
		}
	}

	// Check reserved subdomains
	if _, ok := core.ReservedSubdomains[endpoint]; ok {
		return db.Endpoint{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("URL %s is reserved.", url),
		}
	}

	// Check if the requested endpoint already exists
	exists, err := s.urlq.CheckEndpointExists(c, endpoint)
	if err != nil {
		slog.Error("unable to check if endpoint already exists", "endpoint", endpoint, "username", username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}
	if exists {
		return db.Endpoint{}, &UrlError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("URL %s already exists", url),
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
		slog.Error("unable to get user from username", "username", username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	// Check if user has exceeded number of urls that can be generated
	urls, err := s.urlq.GetNonExpiredEndpointsOfUser(c, pgtype.Int8{Int64: user.ID, Valid: true})
	if err != nil {
		slog.Error("unable to get non expired endpoints", "username", user.Username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	if (user.Plan == db.PlanBasic || user.Plan == db.PlanFree) && len(urls) >= DefaultLimitNumUrl {
		return db.Endpoint{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Cannot generate more than one url for your current plan. Consider upgrading to Pro.",
		}
	}

	// Check reserved companies. If found, check if the mail is from that organisation
	if _, ok := core.ReservedCompanies[endpoint]; ok {
		if !strings.Contains(strings.ToLower(user.Email), endpoint) || strings.Contains(strings.ToLower(user.Email), "@gmail.com") {
			return db.Endpoint{}, &UrlError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("You cannot use this subdomain. Please try with mail issued by %s", endpoint),
			}
		}
	}

	slog.Info("Create url request received", "endpoint", endpoint, "username", username, "plan", user.Plan)

	endpointRecord, err := s.urlq.InsertEndpoint(c, db.InsertEndpointParams{
		Endpoint: endpoint,
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
		slog.Error("unable to insert new url into db", "endpoint", endpoint, "username", user.Username, "err", err)
		return db.Endpoint{}, NewInternalServerError()
	}

	// Send complete url as response
	endpointRecord.Endpoint = url

	slog.Info("URL created", "url", url, "username", user.Username, "plan", user.Plan)

	return endpointRecord, nil
}

func (s *UrlService) GetUserEndpoints(ctx context.Context, userId int64) ([]Endpoint, *UrlError) {
	endpointsRec, err := s.urlq.GetUserEndpoints(ctx, userId)
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

	var expiresAt pgtype.Timestamptz
	switch endpointRecord.Plan {
	case db.PlanFree:
		{

			// Checking if request body exceeds 10KB
			if hookReq.ContentSize > 10_000 {
				return db.Request{}, &UrlError{
					Code:    http.StatusRequestEntityTooLarge,
					Message: "Content size limit exceeded",
				}
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
				return db.Request{}, &UrlError{
					Code:    http.StatusRequestEntityTooLarge,
					Message: "Content size limit exceeded",
				}
			}
			expiresAt = pgtype.Timestamptz{
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

	userId := endpointRecord.UserID

	queryBytes, err := json.Marshal(hookReq.Query)
	if err != nil {
		slog.Error("unable to marshal query params", "err", err)
		return db.Request{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "unable to parse query params.",
		}
	}

	headerBytes, err := json.Marshal(hookReq.Headers)
	if err != nil {
		slog.Error("unable to marshal headers", "err", err)
		return db.Request{}, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "unable to parse headers",
		}
	}

	requestRecord, err := s.urlq.CreateNewRequest(ctx, db.CreateNewRequestParams{
		UserID:     userId,
		EndpointID: endpointRecord.ID,
		Method:     db.HttpMethod(strings.ToLower(hookReq.Method)),
		Content:    pgtype.Text{String: hookReq.Content, Valid: true},
		Path:       hookReq.Path,
		Uuid:       hookReq.UUID,

		// TODO: Fetch response from configured response
		ResponseCode: pgtype.Int4{Int32: http.StatusOK, Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     hookReq.SourceIp,

		// TODO: Add request body limiting
		ContentSize: int32(hookReq.ContentSize),
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		slog.Error("unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return db.Request{}, NewInternalServerError()
	}

	slog.Info("Endpoint record created", "endpoint", endpoint, "userId", userId.Int64)

	return requestRecord, nil
}

func (s *UrlService) GetEndpointRequestHistory(c context.Context, endpoint string, limit int32, offset int32) ([]Request, *UrlError) {
	slog.Info("Fetch endpoint request history", "endpoint", endpoint)

	var reqHistory []Request

	reqs, err := s.urlq.GetEndpointHistory(c, db.GetEndpointHistoryParams{Endpoint: endpoint, Limit: limit, Offset: offset})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return reqHistory, nil
		}
		slog.Error("unable to fetch endpoint request history", "endpoint", endpoint, "err", err)
		return nil, NewInternalServerError()
	}

	for _, req := range reqs {
		rh := Request{
			ID:           req.ID,
			Path:         req.Path,
			Content:      req.Content.String,
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
			slog.Error("unable to fetch request details", "reqId", reqId, "err", err)
			return Request{}, NewInternalServerError()
		}
	}

	req := Request{
		ID:           reqRecord.ID,
		UUID:         reqRecord.Uuid,
		Path:         reqRecord.Path,
		Method:       reqRecord.Method,
		SourceIp:     reqRecord.SourceIp,
		Content:      reqRecord.Content.String,
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
	endpoint = strings.ToLower(endpoint)
	slog.Info("Request endpoint stats", "endpoint", endpoint)

	endpointDetails, err := s.urlq.GetEndpoint(c, endpoint)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return EndpointStats{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Endpoint %v not found", endpoint),
			}
		}
		slog.Error("unable to fetch endpoint details", "endpoint", endpoint, "err", err)
		return EndpointStats{}, NewInternalServerError()
	}

	stats, err := s.urlq.GetEndpointRequestCount(c, endpoint)
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

type SubdomainExists string

const (
	Available         SubdomainExists = "Its available. Sign up and make it yours"
	Taken             SubdomainExists = "That subdomain is already taken. Try something else?"
	ReservedCompany   SubdomainExists = "Subdomain is reserved. But, you can go ahead if you're using mail issued from that organisation."
	ReservedSubdomain SubdomainExists = "Subdomain is reserved."
	BadSubdomain      SubdomainExists = "Bad subdomain."
	Error             SubdomainExists = "Something went wrong."
)

func (s *UrlService) CheckSubdomainExists(c context.Context, endpoint string) (SubdomainExists, *UrlError) {
	endpoint = strings.ToLower(endpoint)

	slog.Info("Checking if endpoint exists", "endpoint", endpoint)

	if len(endpoint) < 4 || len(endpoint) > 10 {
		return BadSubdomain, &UrlError{
			Code:    http.StatusBadRequest,
			Message: "Subdomain should be 4 to 10 characters.",
		}
	}

	if _, ok := core.ReservedSubdomains[endpoint]; ok {
		return ReservedSubdomain, &UrlError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Subdomain %s is reserved.", endpoint),
		}
	}

	// Check reserved companies. If found, check if the mail is from that organisation
	if _, ok := core.ReservedCompanies[endpoint]; ok {
		return ReservedCompany, nil
	}

	exists, err := s.urlq.CheckEndpointExists(c, endpoint)
	if err != nil {
		slog.Error("unable to check if endpoint exists", "endpoint", endpoint, "err", err)
		return Error, NewInternalServerError()
	}

	if exists {
		fmt.Println("Exists", exists)
		return Taken, nil
	}

	return Available, nil
}

func (s *UrlService) GetRequestByUUID(c context.Context, uuid string) (Request, *UrlError) {

	reqRecord, err := s.urlq.GetRequestByUUID(c, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Request{}, &UrlError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("No request found for uuid: %v", uuid),
			}
		} else {
			slog.Error("unable to fetch request details", "uuid", uuid, "err", err)
			return Request{}, NewInternalServerError()
		}
	}

	req := Request{
		ID:           reqRecord.ID,
		UUID:         reqRecord.Uuid,
		Path:         reqRecord.Path,
		Method:       reqRecord.Method,
		SourceIp:     reqRecord.SourceIp,
		Content:      reqRecord.Content.String,
		ContentSize:  reqRecord.ContentSize,
		ResponseCode: reqRecord.ResponseCode,
		CreatedAt:    reqRecord.CreatedAt,
		ExpiresAt:    reqRecord.ExpiresAt,
	}

	json.Unmarshal(reqRecord.Headers, &req.Headers)
	json.Unmarshal(reqRecord.QueryParams, &req.QueryParams)

	return req, nil
}

func (s *UrlService) ExpireRequests(c context.Context) error {
	slog.Info("Deleting expired requests", "date", time.Now().Local().String())
	err := s.urlq.ExpireRequests(c)
	if err != nil {
		slog.Error("unable to delete expired requests", "date", time.Now().Local().String(), "err", err)
		return err
	}

	return nil
}
