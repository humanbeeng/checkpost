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

type UrlError struct {
	Code    int
	Message string
}

func (u *UrlError) Error() string {
	return u.Message
}

func NewInternalServerError() *UrlError {
	return &UrlError{
		Code:    http.StatusInternalServerError,
		Message: "Oops! Something went wrong :(",
	}
}

func (s *UrlService) CreateUrl(c context.Context, username string, endpoint string) (string, *UrlError) {
	// TODO: Min len of 4
	if len(endpoint) < 4 {
		return "", &UrlError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprint("Endpoint should be atleast 4 characters"),
		}
	}

	user, err := s.q.GetUserFromUsername(c, username)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", &UrlError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("No user found with username: %s", username),
		}
	}

	switch user.Plan {
	case db.PlanFree:
		{
			return s.CreateRandomUrl(c)
		}
	case db.PlanNoBrainer, db.PlanPro:
		{
			// TODO: Check number of endpoints limit
			url := fmt.Sprintf("https://%v.checkpost.io", endpoint)

			// Check if the requested endpoint exists
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
	return "", &UrlError{Code: http.StatusBadRequest, Message: "invalid user plan"}
}

func (s *UrlService) StoreRequestDetails(c *fiber.Ctx) *UrlError {
	endpoint := c.Params("endpoint", "")
	if endpoint == "" {
		return &UrlError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Empty endpoint"),
		}
	}

	endpointRecord, err := s.q.GetEndpoint(c.Context(), endpoint)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return &UrlError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("https://%s.checkpost.io is either not created or has expired", endpoint),
		}
	}

	slog.Info("Storing request details", "endpoint", endpoint)

	userId := endpointRecord.UserID

	var req any
	_ = c.BodyParser(&req)

	strBytes, _ := json.Marshal(req)
	body := string(strBytes)
	headers := c.GetReqHeaders()
	ip := c.IP()
	path := c.Params("path", "/")

	method := c.Method()
	query := c.Queries()

	queryBytes, err := json.Marshal(query)
	if err != nil {
		slog.Error("Unable to marshal query params", "err", err)
		return &UrlError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprint("Unable to parse query params"),
		}
	}

	headerBytes, err := json.Marshal(headers)
	if err != nil {
		slog.Error("Unable to marshal headers", "err", err)
		return &UrlError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprint("Unable to parse headers"),
		}
	}

	if _, err := s.q.CreateNewRequest(c.Context(), db.CreateNewRequestParams{
		UserID:       userId,
		EndpointID:   endpointRecord.ID,
		Method:       db.HttpMethod(strings.ToLower(method)),
		Content:      pgtype.Text{String: body, Valid: true},
		Path:         path,
		ResponseCode: pgtype.Int4{Int32: http.StatusOK, Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     ip,
		ContentSize:  int32(len(body)),
	}); err != nil {
		slog.Error("Unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return NewInternalServerError()
	}

	slog.Info("Endpoint record created", "endpoint", endpoint, "userId", userId.Int64)

	return nil
}

func (s *UrlService) CreateRandomUrl(c context.Context) (string, *UrlError) {
	// TODO: length from config ?
	randomEndpoint, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 10)
	if err != nil {
		slog.Error("Unable to generate nano id", "err", err)
		return "", NewInternalServerError()
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
		return "", NewInternalServerError()
	}

	slog.Info("Random url generated", "url", randomUrl)
	return randomUrl, nil
}
