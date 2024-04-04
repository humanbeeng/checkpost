package url

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type UrlService struct {
	q db.Querier
}

func NewUrlService(q db.Querier) *UrlService {
	return &UrlService{q: q}
}

var ErrInternal = errors.New("Sorry! Something went wrong :(")

func (u *UrlService) GenerateUrl(c context.Context, username string, endpoint string) (string, error) {
	// Check if the user is guest
	if username == "" {
		return u.GenerateRandomUrlAndInsertIntoDb(c)
	}

	user, err := u.q.GetUserFromUsername(c, username)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		slog.Error("No user found", "username", username)
		return "", fmt.Errorf("User %v not found", username)
	}

	switch user.Plan {
	case db.PlanFree:
		{
			return u.GenerateRandomUrlAndInsertIntoDb(c)
		}
	case db.PlanNoBrainer, db.PlanPro:
		{

			// TODO: Check number of endpoints limit

			// Check if the requested endpoint exists
			exists, err := u.q.CheckEndpointExists(c, endpoint)
			if err != nil {
				slog.Error("Unable to check if endpoint already exists", "err", err)
				return "", ErrInternal
			}
			if exists {
				slog.Info("Endpoint already exists", "endpoint", endpoint)
				return "", fmt.Errorf("URL %v already exists", endpoint)
			}

			// endpoint is available
			url := fmt.Sprintf("https://%v.checkpost.io", endpoint)
			_, err = u.q.CreateNewEndpoint(c, db.CreateNewEndpointParams{
				Endpoint: endpoint,
				UserID:   pgtype.Int8{Int64: user.ID, Valid: true},
				Plan:     user.Plan,
			})
			if err != nil {
				slog.Error("Unable to insert new url into db", "endpoint", endpoint, "user", user.ID, "err", err)
				return "", err
			}

			slog.Info("Endpoint created and inserted into db", "endpoint", endpoint, "user", user.ID)

			return url, nil
		}
	}
	return "", fmt.Errorf("invalid user plan")
}

func (s *UrlService) StoreRequestDetails(c *fiber.Ctx) error {
	endpoint := c.Params("endpoint", "")
	if endpoint == "" {
		return fiber.ErrNotFound
	}

	endpointRecord, err := s.q.GetEndpoint(c.Context(), endpoint)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Endpoint not found", "endpoint", endpoint)
		return fiber.ErrNotFound
	}
	slog.Info("Endpoint record found", "endpoint", endpoint)

	userId := endpointRecord.UserID

	var req any
	_ = c.BodyParser(&req)

	strBytes, _ := json.Marshal(req)
	body := string(strBytes)
	headers := c.GetReqHeaders()
	ip := c.IP()
	path := c.Path()
	path, found := strings.CutPrefix(path, "/url/hook")
	if !found {
		return fiber.ErrBadRequest
	}

	method := c.Method()
	query := c.Queries()

	queryBytes, err := json.Marshal(query)
	if err != nil {
		slog.Error("Unable to marshal query params", "err", err)
		return fiber.ErrBadRequest
	}

	headerBytes, err := json.Marshal(headers)
	if err != nil {
		slog.Error("Unable to marshal headers", "err", err)
		return fiber.ErrBadRequest
	}

	if _, err := s.q.CreateNewRequest(c.Context(), db.CreateNewRequestParams{
		UserID:       userId,
		EndpointID:   endpointRecord.ID,
		Method:       db.HttpMethod(strings.ToLower(method)),
		Content:      pgtype.Text{String: body, Valid: true},
		Path:         path,
		ResponseCode: pgtype.Int4{Int32: fiber.StatusOK, Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     ip,
		ContentSize:  int32(len(body)),
	}); err != nil {
		slog.Error("Unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return fiber.ErrInternalServerError
	}

	slog.Info("Endpoint record created", "endpoint", endpoint, "userId", userId)

	return c.SendStatus(fiber.StatusOK)
}

func (s *UrlService) GenerateRandomUrlAndInsertIntoDb(c context.Context) (string, error) {
	// TODO: length from config ?
	randomEndpoint, err := gonanoid.Generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 10)
	if err != nil {
		slog.Error("Unable to generate nano id", "err", err)
		return "", err
	}

	// TODO: Fetch base url from config file
	randomUrl := fmt.Sprintf("https://%v.checkpost.local", randomEndpoint)

	// Inserting into db assuming that no endpoint with that random url existed. We can add
	// a check later on if needed.
	if _, err := s.q.CreateNewGuestEndpoint(c, randomEndpoint); err != nil {
		slog.Error("Unable to insert endpoint", "endpoint", randomUrl, "err", err)
		return "", err
	}
	return randomUrl, nil
}
