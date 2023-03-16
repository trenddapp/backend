package user

import (
	"context"

	"github.com/trenddapp/backend/service/user/pkg/model"
)

type Repository interface {
	GetUser(ctx context.Context, id string) (model.User, error)
	GetUserByAddress(ctx context.Context, address string) (model.User, error)
	CreateUser(ctx context.Context, in model.User) (model.User, error)
	UpdateUser(ctx context.Context, in model.User) (model.User, error)
	DeleteUser(ctx context.Context, id string) (model.User, error)
}
