package coinmarketcap

import (
	"context"

	"github.com/trenddapp/backend/service/currency/pkg/model"
)

type Client interface {
	GetRate(ctx context.Context, symbol string) (model.Rate, error)
}
