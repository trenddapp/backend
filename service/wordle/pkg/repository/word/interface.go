package word

import (
	"context"

	"github.com/trenddapp/backend/service/wordle/pkg/model"
)

type Repository interface {
	GetRandomWord(ctx context.Context, locale string) (model.Word, error)
	IsValidWord(ctx context.Context, in model.Word) (bool, error)
}
