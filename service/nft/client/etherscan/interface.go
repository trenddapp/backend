package etherscan

import (
	"context"

	"github.com/dapp-z/backend/service/nft/model"
)

type Client interface {
	GetAccountNFTs(ctx context.Context, address string) ([]model.NFT, error)
}
