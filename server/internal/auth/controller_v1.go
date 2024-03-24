package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthController struct {
	config         *oauth2.Config
	pasetoVerifier *PasetoVerifier
}

func NewGithubAuthController() (*AuthController, error) {
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

	return &AuthController{
		config:         config,
		pasetoVerifier: pasetoVerifier,
	}, err
}

type GithubUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
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

func (ac *AuthController) RegisterRoutes(router fiber.Router) {
	router.Get("/auth/github", ac.LoginHandler)
	router.Get("/auth/github/callback", ac.CallbackHandler)
}

func (ac *AuthController) LoginHandler(c *fiber.Ctx) error {
	fmt.Println("Received login request")
	url := ac.config.AuthCodeURL("not-implemented-yet")
	return c.Redirect(url)
}

func (ac *AuthController) CallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return fmt.Errorf("No auth code received")
	}
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

	var githubUser GithubUser
	err = json.Unmarshal(body, &githubUser)
	if err != nil {
		return err
	}

	// Create a token with 1 month expiry time.
	pasetoToken, err := ac.pasetoVerifier.CreateToken(githubUser.Username, time.Hour*24*30)
	if err != nil {
		return err
	}

	response := AuthResponse{Token: pasetoToken}

	return c.JSON(response)
}
