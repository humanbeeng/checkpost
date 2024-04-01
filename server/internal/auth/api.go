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
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/humanbeeng/checkpost/server/db/sqlc"
	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	config         *oauth2.Config
	pasetoVerifier *PasetoVerifier
	q              db.Querier
}

func NewGithubAuthHandler(querier db.Querier) (*AuthHandler, error) {
	key := os.Getenv("GITHUB_KEY")
	secret := os.Getenv("GITHUB_SECRET")
	symmetricKey := os.Getenv("PASETO_KEY")

	pasetoVerifier, err := NewPasetoVerifier(symmetricKey)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     key,
		ClientSecret: secret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{},
	}

	return &AuthHandler{
		config:         config,
		pasetoVerifier: pasetoVerifier,
		q:              querier,
	}, err
}

func (ac *AuthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/auth/github", ac.LoginHandler)
	app.Get("/auth/github/callback", ac.CallbackHandler)
	// TODO: Remove this
	app.Get("/auth/token", ac.GenerateToken)
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
	// TODO: Add state to oauth request
	url := a.config.AuthCodeURL("none")
	return c.Redirect(url)
}

func (a *AuthHandler) CallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		slog.Warn("Code not found in callback url")
	}

	githubUser, err := a.exchangeCodeForUser(c, code)
	if err != nil {
		slog.Error("Unable to exchange code for github user", "err", err)
		return fiber.ErrInternalServerError
	}

	user, err := a.q.GetUserFromEmail(context.Background(), githubUser.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = a.q.CreateUser(c.Context(), db.CreateUserParams{
				Name:     githubUser.Name,
				Username: githubUser.Username,
				Plan:     db.PlanFree,
				Email:    githubUser.Email,
			})
			if err != nil {
				slog.Error("Unable to create new user", "err", err)
				return fiber.ErrInternalServerError
			}
		} else {
			slog.Error("Unable to fetch existing user", "err", err)
		}
	}

	// Create token and encrypt it
	pasetoToken, err := a.pasetoVerifier.CreateToken(user.Username, user.ID, time.Hour*24*30)
	if err != nil {
		fmt.Println("err", err)
		return err
	}

	res := AuthResponse{Token: pasetoToken}
	return c.JSON(res)
}

// exchange the auth code that retrieved from github via URL query parameter into an access token.
func (a *AuthHandler) exchangeCodeForUser(c *fiber.Ctx, code string) (*GithubUser, error) {
	token, err := a.config.Exchange(c.Context(), code)
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

func (ac *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	token, err := ac.pasetoVerifier.CreateToken("humanbeeng", 1, time.Hour*12)
	if err != nil {
		slog.Error("err", "err", err)
		return fiber.ErrInternalServerError
	}

	return c.SendString(token)
}
