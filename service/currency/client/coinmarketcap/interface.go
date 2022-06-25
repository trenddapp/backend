package coinmarketcap

import (
	"context"

	"github.com/trenddapp/backend/service/currency/model"
)

type Client interface {
	GetRate(ctx context.Context, symbol string) (model.Rate, error)
}
