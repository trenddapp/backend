package wordle

import (
	"context"

	"github.com/trenddapp/backend/service/wordle/pkg/model"
)

type Repository interface {
	GetWordle(ctx context.Context, id string, opts ...Option) (model.Wordle, error)
	CreateWordle(ctx context.Context, in model.Wordle) (model.Wordle, error)
	UpdateWordle(ctx context.Context, in model.Wordle, opts ...Option) (model.Wordle, error)
	DeleteWordle(ctx context.Context, id string, opts ...Option) (model.Wordle, error)
	ListWordles(ctx context.Context, opts ...Option) ([]model.Wordle, string, error)
}
