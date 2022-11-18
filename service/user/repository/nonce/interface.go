package nonce

import (
	"context"

	"github.com/trenddapp/backend/service/user/model"
)

type Repository interface {
	GetNonceByUserID(ctx context.Context, userID string) (model.Nonce, error)
	CreateNonce(ctx context.Context, in model.Nonce) (model.Nonce, error)
	DeleteNonce(ctx context.Context, id string) (model.Nonce, error)
}
