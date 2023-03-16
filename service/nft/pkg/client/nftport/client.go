package nftport

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/trenddapp/backend/pkg/paginator"
	"github.com/trenddapp/backend/service/nft/pkg/model"
)

const (
	Chain    = "rinkeby"
	StatusOK = "OK"
	URL      = "https://api.nftport.xyz"
)

var (
	ErrUnknown = errors.New("unknown error")
)

type client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(config *Config) Client {
	return &client{
		apiKey:     config.APIKey,
		httpClient: http.DefaultClient,
	}
}

func (c *client) GetNFT(ctx context.Context, contractAddress string, tokenID string) (model.NFT, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL+"/v0/nfts/"+contractAddress+"/"+tokenID, http.NoBody)
	if err != nil {
		return model.NFT{}, err
	}

	query := request.URL.Query()
	query.Add("chain", Chain)
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Authorization", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return model.NFT{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return model.NFT{}, ErrUnknown
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return model.NFT{}, err
	}

	var result struct {
		NFT struct {
			ContractAddress string `json:"contract_address"`
			TokenID         string `json:"token_id"`
		} `json:"nft"`
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return model.NFT{}, err
	}

	if result.Response != StatusOK {
		return model.NFT{}, ErrUnknown
	}

	return model.NFT{
		ContractAddress: result.NFT.ContractAddress,
		TokenID:         result.NFT.TokenID,
	}, nil
}

func (c *client) ListAccountNFTs(
	ctx context.Context,
	address string,
	pageSize int,
	pageToken string,
) ([]model.NFT, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL+"/v0/accounts/"+address, http.NoBody)
	if err != nil {
		return nil, "", err
	}

	continuation := ""
	if pageToken != "" {
		id, _, err := paginator.ParsePageToken(pageToken)
		if err != nil {
			return nil, "", err
		}

		continuation = id
	}

	query := request.URL.Query()
	query.Add("chain", Chain)
	query.Add("continuation", continuation)
	query.Add("include", "metadata")
	query.Add("page_number", "1")
	query.Add("page_size", "20")
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Authorization", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, "", err
	}

	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", ErrUnknown
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	var result struct {
		Continuation string `json:"continuation"`
		NFTs         []struct {
			ContractAddress string `json:"contract_address"`
			TokenID         string `json:"token_id"`
		} `json:"nfts"`
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}

	if result.Response != StatusOK {
		return nil, "", ErrUnknown
	}

	nfts := make([]model.NFT, 0, len(result.NFTs))

	for _, nft := range result.NFTs {
		nfts = append(nfts, model.NFT{
			ContractAddress: nft.ContractAddress,
			TokenID:         nft.TokenID,
		})
	}

	nextPageToken := ""
	if len(nfts) == pageSize {
		nextPageToken = paginator.GeneratePageToken(result.Continuation, "UNSPECIFIED")
	}

	return nfts, nextPageToken, nil
}

func (c *client) ListContractNFTs(
	ctx context.Context,
	address string,
	pageSize int,
	pageToken string,
) ([]model.NFT, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL+"/v0/nfts/"+address, http.NoBody)
	if err != nil {
		return nil, "", err
	}

	query := request.URL.Query()
	query.Add("chain", Chain)
	query.Add("include", "all")
	query.Add("page_number", "1")
	query.Add("page_size", strconv.Itoa(pageSize))
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Authorization", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", ErrUnknown
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	var result struct {
		NFTs []struct {
			ContractAddress string `json:"contract_address"`
			TokenID         string `json:"token_id"`
		} `json:"nfts"`
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}

	if result.Response != StatusOK {
		return nil, "", ErrUnknown
	}

	nfts := make([]model.NFT, 0, len(result.NFTs))

	for _, nft := range result.NFTs {
		nfts = append(nfts, model.NFT{
			ContractAddress: nft.ContractAddress,
			TokenID:         nft.TokenID,
		})
	}

	return nfts, "", nil
}
