package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/humanbeeng/checkpost/server/config"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/humanbeeng/checkpost/server/internal/core"
	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	config         *config.AppConfig
	oauthConfig    *oauth2.Config
	pasetoVerifier *core.PasetoVerifier
	q              db.Querier
}

func NewGithubAuthHandler(config *config.AppConfig, querier db.Querier) (*AuthHandler, error) {
	key := config.Github.Key
	secret := config.Github.Secret
	symmetricKey := config.Paseto.Key

	pasetoVerifier, err := core.NewPasetoVerifier(symmetricKey)
	if err != nil {
		return nil, err
	}

	oauthConfig := &oauth2.Config{
		ClientID:     key,
		ClientSecret: secret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{},
	}

	return &AuthHandler{
		oauthConfig:    oauthConfig,
		config:         config,
		pasetoVerifier: pasetoVerifier,
		q:              querier,
	}, err
}

func (ac *AuthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/auth/github", ac.LoginHandler)
	app.Get("/auth/github/callback", ac.CallbackHandler)
}

type GithubUser struct {
	Name      string `json:"name"`
	Username  string `json:"login"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

func (a *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	slog.Info("Received login request")
	// TODO: Add state to oauth request
	a.oauthConfig.Scopes = append(a.oauthConfig.Scopes, "email")
	url := a.oauthConfig.AuthCodeURL("none")
	return c.Redirect(url)
}

func (a *AuthHandler) CallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		slog.Warn("Code not found in callback url")
	}

	githubUser, err := a.exchangeCodeForUser(c, code)
	if err != nil {
		slog.Error("unable to exchange code for github user", "err", err)
		return fiber.ErrInternalServerError
	}

	fmt.Println("Fetched github user", githubUser)

	user, err := a.q.GetUserFromEmail(context.Background(), githubUser.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = a.q.CreateUser(c.Context(), db.CreateUserParams{
				Name:      githubUser.Name,
				AvatarUrl: githubUser.AvatarUrl,
				Username:  githubUser.Username,
				Plan:      db.PlanFree,
				Email:     githubUser.Email,
			})
			if err != nil {
				slog.Error("unable to create new user", "err", err)
				return fiber.ErrInternalServerError
			}
		} else {
			slog.Error("unable to fetch existing user", "err", err)
		}
	}

	// Create token and encrypt it
	args := core.CreateTokenArgs{
		Username: user.Username,
		UserId:   user.ID,
		Plan:     db.PlanFree,
	}
	pasetoToken, err := a.pasetoVerifier.CreateToken(args, time.Hour*24*30)
	if err != nil {
		return err
	}

	res := AuthResponse{Token: pasetoToken}
	return c.JSON(res)
}

// exchange the auth code that retrieved from github via URL query parameter into an access token.
func (a *AuthHandler) exchangeCodeForUser(c *fiber.Ctx, code string) (*GithubUser, error) {
	token, err := a.oauthConfig.Exchange(c.Context(), code)
	if err != nil {
		return nil, err
	}

	baseUrl, err := url.Parse("https://api.github.com/user")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(c.Context(), http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	githubRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer githubRes.Body.Close()

	body, err := io.ReadAll(githubRes.Body)
	if err != nil {
		return nil, err
	}

	var githubUser GithubUser
	err = json.Unmarshal(body, &githubUser)
	if err != nil {
		return nil, err
	}

	return &githubUser, nil
}
