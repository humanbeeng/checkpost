package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthController struct {
	config *oauth2.Config
}

func NewGithubAuthController() *AuthController {
	key := os.Getenv("GITHUB_KEY")
	secret := os.Getenv("GITHUB_SECRET")
	redirectUrl := os.Getenv("GITHUB_REDIRECT")

	config := &oauth2.Config{
		ClientID:     key,
		ClientSecret: secret,
		RedirectURL:  redirectUrl,
		Endpoint:     github.Endpoint,
		Scopes:       []string{},
	}

	return &AuthController{
		config: config,
	}
}

type AuthenticatedUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}

func (ac *AuthController) RegisterRoutes(router fiber.Router) {
	router.Get("/auth/github", ac.LoginHandler)
	router.Get("/auth/github/callback", ac.CallbackHandler)
	router.Get("/auth/logout", ac.LogoutHandler)
}

func (ac *AuthController) LoginHandler(c *fiber.Ctx) error {
	url := ac.config.AuthCodeURL("not-implemented-yet")
	return c.Redirect(url)
}

func (ac *AuthController) CallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	// exchange the auth code that retrieved from github via
	// URL query parameter into an access token.
	token, err := ac.config.Exchange(c.Context(), code)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	baseUrl, err := url.Parse("https://api.github.com/user")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var authenticatedUser AuthenticatedUser

	err = json.Unmarshal(body, &authenticatedUser)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", authenticatedUser)

	return nil
}

func (ac *AuthController) LogoutHandler(c *fiber.Ctx) error {
	return nil
}
