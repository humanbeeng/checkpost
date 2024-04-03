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

var (
	ErrEndpointAlreadyExists = errors.New("endpoint already exists")
	ErrNoUser                = errors.New("no user found")
)

func (u *UrlService) GenerateUrl(c context.Context, username string, endpoint string) (string, error) {
	slog.Info("Generate url request received", "username", username, "endpoint", endpoint)
	user, err := u.q.GetUserFromUsername(c, username)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		if username == "" {
			return u.generateRandomUrlAndInsertIntoDb(c)
		} else {
			slog.Warn("No user found", "username", username)
			return "", ErrNoUser
		}
	}

	switch user.Plan {
	case db.PlanFree:
		{
			return u.generateRandomUrlAndInsertIntoDb(c)
		}
	case db.PlanNoBrainer, db.PlanPro:
		{
			// Check if endpoint exists
			_, err := u.q.GetEndpoint(c, endpoint)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					// endpoint is available
					url := fmt.Sprintf("https://%v.checkpost.local", endpoint)
					_, err := u.q.CreateNewEndpoint(c, db.CreateNewEndpointParams{
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
				} else {
					slog.Error("Unable to fetch existing url", "endpoint", endpoint)
					return "", err
				}
			} else {
				slog.Error("endpoint already exists", "endpoint", endpoint)
				return "", ErrEndpointAlreadyExists
			}
		}
	}

	return "", nil
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
		ResponseCode: pgtype.Int4{Int32: fiber.StatusOK, Valid: true},
		QueryParams:  queryBytes,
		Headers:      headerBytes,
		SourceIp:     ip,
		ContentSize:  int32(len(body)),
	}); err != nil {
		slog.Error("Unable to create new request record", "endpoint", endpoint, "userId", userId, "err", err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *UrlService) generateRandomUrlAndInsertIntoDb(c context.Context) (string, error) {
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
	if _, err := s.q.CreateNewEndpoint(c, db.CreateNewEndpointParams{
		Endpoint: randomEndpoint,
	}); err != nil {
		// TODO: Add a check to see if the random url also exists (improbable, but i think we can add it// without any cost)
		slog.Error("Unable to insert endpoint", "endpoint", randomUrl, "err", err)
		return "", err
	}
	return randomUrl, nil
}
