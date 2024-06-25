package auth

import (
	"context"
	"encoding/json"
	"errors"
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

type MailResponseItem struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func (a *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	slog.Info("Received login request")
	// TODO: Add state to oauth request
	a.oauthConfig.Scopes = append(a.oauthConfig.Scopes, "user:email")
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

	// Fetch basic user information
	userReq, err := http.NewRequestWithContext(c.Context(), http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	userReq.Header.Add("Authorization", "Bearer "+token.AccessToken)

	userRes, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return nil, err
	}

	// Fetch user emails
	emailUrl, err := url.Parse("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}

	emailReq, err := http.NewRequestWithContext(c.Context(), http.MethodGet, emailUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	emailReq.Header.Add("Authorization", "Bearer "+token.AccessToken)

	emailRes, err := http.DefaultClient.Do(emailReq)
	if err != nil {
		return nil, err
	}

	defer userRes.Body.Close()
	defer emailRes.Body.Close()

	userBody, err := io.ReadAll(userRes.Body)
	if err != nil {
		return nil, err
	}

	emailBody, err := io.ReadAll(emailRes.Body)
	if err != nil {
		return nil, err
	}

	var githubUser GithubUser
	err = json.Unmarshal(userBody, &githubUser)
	if err != nil {
		return nil, err
	}

	if githubUser.Name == "" {
		githubUser.Name = githubUser.Username
	}

	var userEmails []MailResponseItem
	err = json.Unmarshal(emailBody, &userEmails)
	if err != nil {
		return nil, err
	}

	for _, email := range userEmails {
		if email.Primary {
			githubUser.Email = email.Email
			break
		}
	}

	return &githubUser, nil
}
