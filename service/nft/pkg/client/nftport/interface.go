package nftport

import (
	"context"

	"github.com/trenddapp/backend/service/nft/pkg/model"
)

type Client interface {
	GetNFT(ctx context.Context, contractAddress string, tokenID string) (model.NFT, error)
	ListAccountNFTs(
		ctx context.Context,
		address string,
		pageSize int,
		pageToken string,
	) ([]model.NFT, string, error)
	ListContractNFTs(
		ctx context.Context,
		address string,
		pageSize int,
		pageToken string,
	) ([]model.NFT, string, error)
}
