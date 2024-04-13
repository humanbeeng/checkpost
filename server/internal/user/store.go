package user

import (
	"context"

	db "github.com/humanbeeng/checkpost/server/db/sqlc"
)

type UserQuerier interface {
	GetUserFromUsername(ctx context.Context, username string) (db.User, error)
}

type UserStore struct {
	q db.Querier
}

func NewUserStore(q db.Querier) *UserStore {
	return &UserStore{
		q: q,
	}
}

func (us UserStore) GetUserFromUsername(ctx context.Context, username string) (db.User, error) {
	return us.q.GetUserFromUsername(ctx, username)
}
