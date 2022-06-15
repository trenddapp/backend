package nftport

import (
	"context"

	"github.com/dapp-z/backend/service/nft/model"
)

type Client interface {
	GetAccountNFTs(ctx context.Context, address string) ([]model.NFT, error)
	GetContractNFTs(ctx context.Context, address string) ([]model.NFT, error)
	GetNFT(ctx context.Context, contractAddress string, tokenID string) (*model.NFT, error)
}
