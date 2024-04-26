package user

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	store *UserStore
}

func NewUserController(store *UserStore) *UserController {
	return &UserController{
		store: store,
	}
}

func (uc *UserController) RegisterRoutes(app *fiber.App, authmw fiber.Handler) {
	urlGroup := app.Group("/user")

	urlGroup.Get("/", authmw, uc.GetUserDetailsHandler)
}

type UserDetailsResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Plan     string `json:"plan"`
}

func (uc *UserController) GetUserDetailsHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	slog.Info("Requesting user details", "userId", userId)

	user, err := uc.store.GetUserFromUserId(c.Context(), userId)
	if err != nil {
		return fiber.ErrNotFound
	}

	res := UserDetailsResponse{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Plan:     string(user.Plan),
	}

	return c.JSON(res)
}
