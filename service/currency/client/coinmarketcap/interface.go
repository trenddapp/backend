package coinmarketcap

import "context"

type Client interface {
	GetConversionRate(ctx context.Context, symbol string) (float64, error)
}
