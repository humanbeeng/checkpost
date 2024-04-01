package url

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

func (u *UrlService) GenerateUrl(c context.Context, email string, endpoint string) (string, error) {
	user, err := u.q.GetUserFromEmail(c, email)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		if email == "" {
			return u.generateRandomUrlAndInsertIntoDb(c)
		} else {
			slog.Warn("No user found", "email", email)
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

func (s *UrlService) generateRandomUrlAndInsertIntoDb(c context.Context) (string, error) {
	randomEndpoint, err := gonanoid.Generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 12)
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
		slog.Error("Unable to insert endpoint", "endpoint", randomUrl, "err", err)
		return "", err
	}
	return randomUrl, nil
}
