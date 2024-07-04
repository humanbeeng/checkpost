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
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	config            *config.AppConfig
	githubOAuthConfig *oauth2.Config
	googleOAuthConfig *oauth2.Config
	pasetoVerifier    *core.PasetoVerifier
	q                 db.Querier
}

func NewAuthHandler(config *config.AppConfig, querier db.Querier) (*AuthHandler, error) {
	symmetricKey := config.Paseto.Key

	pasetoVerifier, err := core.NewPasetoVerifier(symmetricKey)
	if err != nil {
		return nil, err
	}

	githubOauthConfig := &oauth2.Config{
		ClientID:     config.Github.ClientId,
		ClientSecret: config.Github.Secret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}

	googleOAuthConfig := &oauth2.Config{
		ClientID:     config.Google.ClientId,
		ClientSecret: config.Google.Secret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"email"},
		RedirectURL:  config.Google.RedirectUrl,
	}

	return &AuthHandler{
		githubOAuthConfig: githubOauthConfig,
		googleOAuthConfig: googleOAuthConfig,
		config:            config,
		pasetoVerifier:    pasetoVerifier,
		q:                 querier,
	}, err
}

func (ac *AuthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/auth/github", ac.GithubLoginHandler)
	app.Get("/auth/google", ac.GoogleLoginHandler)
	app.Get("/auth/github/callback", ac.GithubCallbackHandler)
	app.Get("/auth/google/callback", ac.GoogleCallbackHandler)
}

type OAuthUser struct {
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

func (a *AuthHandler) GithubLoginHandler(c *fiber.Ctx) error {
	slog.Info("Received Github login request")
	// TODO: Add state to oauth request
	url := a.githubOAuthConfig.AuthCodeURL("none")
	return c.Redirect(url)
}

func (a *AuthHandler) GoogleLoginHandler(c *fiber.Ctx) error {
	slog.Info("Received Google login request")
	// TODO: Add state to oauth request
	url := a.googleOAuthConfig.AuthCodeURL("none")
	return c.Redirect(url)
}

func (a *AuthHandler) GoogleCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")

	if code == "" {
		slog.Warn("Code not found in callback url")
		return fiber.ErrBadRequest
	}

	slog.Info("Received Google callback")

	googleUser, err := a.exchangeGoogleCodeForUser(c, code)
	if err != nil {
		slog.Error("unable to exchange code for google user", "err", err)
		return fiber.ErrInternalServerError
	}

	user, err := a.q.GetUserFromEmail(context.Background(), googleUser.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info("Creating new user", "username", googleUser.Name, "email", googleUser.Email)
			user, err = a.q.CreateUser(c.Context(), db.CreateUserParams{
				Name:      googleUser.Name,
				AvatarUrl: googleUser.AvatarUrl,
				Username:  googleUser.Username,
				Plan:      db.PlanFree,
				Email:     googleUser.Email,
			})
			if err != nil {
				slog.Error("unable to create new user", "err", err)
				return fiber.ErrInternalServerError
			}
		} else {
			slog.Error("unable to fetch existing user", "err", err)
		}
	} else {
		slog.Info("Logging in existing user", "username", user.Username, "email", user.Email)
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

func (a *AuthHandler) GithubCallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		slog.Warn("Code not found in callback url")
		return fiber.ErrBadRequest
	}

	githubUser, err := a.exchangeGithubCodeForUser(c, code)
	if err != nil {
		slog.Error("unable to exchange code for github user", "err", err)
		return fiber.ErrInternalServerError
	}

	user, err := a.q.GetUserFromEmail(context.Background(), githubUser.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info("Creating new user", "username", githubUser.Username, "email", githubUser.Email)
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
	} else {
		slog.Info("Logging in existing user", "username", user.Username, "email", user.Email)
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

func (a *AuthHandler) exchangeGoogleCodeForUser(c *fiber.Ctx, code string) (*OAuthUser, error) {
	slog.Info("Exchanging Google code for user")

	token, err := a.googleOAuthConfig.Exchange(c.Context(), code)
	if err != nil {
		slog.Error("unable to exchange google callback code for token", "err", err)
		return nil, err
	}

	baseUrl, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
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

	defer userRes.Body.Close()

	userBody, _ := io.ReadAll(userRes.Body)

	var googleUser struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	err = json.Unmarshal(userBody, &googleUser)
	if err != nil {
		return nil, err
	}

	oauthUser := OAuthUser{
		Name:      googleUser.Name,
		Email:     googleUser.Email,
		AvatarUrl: googleUser.Picture,
		Username:  googleUser.Email,
	}

	return &oauthUser, nil
}

// exchange the auth code that retrieved from github via URL query parameter into an access token.
func (a *AuthHandler) exchangeGithubCodeForUser(c *fiber.Ctx, code string) (*OAuthUser, error) {
	slog.Info("Exchanging Github code for user")
	token, err := a.githubOAuthConfig.Exchange(c.Context(), code)
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

	var githubUser OAuthUser
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
