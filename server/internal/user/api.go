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
	urlGroup.Get("/another", authmw, uc.GetUserDetailsHandler)
}

type UserDetailsResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Plan      string `json:"plan"`
}

func (uc *UserController) GetUserDetailsHandler(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	slog.Info("Requesting user details", "userId", userId)

	user, err := uc.store.GetUserFromUserId(c.Context(), userId)
	if err != nil {
		return fiber.ErrNotFound
	}

	res := UserDetailsResponse{
		Id:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Plan:      string(user.Plan),
		AvatarUrl: user.AvatarUrl,
	}

	return c.JSON(res)
}
