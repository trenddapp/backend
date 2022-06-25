package etherscan

import (
	"context"

	"github.com/trenddapp/backend/service/nft/model"
)

type Client interface {
	ListAccountNFTs(
		ctx context.Context,
		address string,
		pageSize int,
		pageToken string,
	) ([]model.NFT, string, error)
}
