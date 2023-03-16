package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/trenddapp/backend/pkg/paginator"
	"github.com/trenddapp/backend/service/nft/pkg/model"
)

const (
	URL = "https://api-rinkeby.etherscan.io/api"
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

func (c *client) ListAccountNFTs(
	ctx context.Context,
	address string,
	pageSize int,
	pageToken string,
) ([]model.NFT, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL+"", http.NoBody)
	if err != nil {
		return nil, "", err
	}

	page := 1
	if pageToken != "" {
		id, _, err := paginator.ParsePageToken(pageToken)
		if err != nil {
			return nil, "", err
		}

		page, err = strconv.Atoi(id)
		if err != nil {
			return nil, "", err
		}
	}

	query := request.URL.Query()
	query.Add("action", "tokennfttx")
	query.Add("address", address)
	query.Add("apikey", c.apiKey)
	query.Add("module", "account")
	query.Add("offset", strconv.Itoa(pageSize))
	query.Add("page", strconv.Itoa(page))
	request.URL.RawQuery = query.Encode()

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
		Transactions []struct {
			ContractAddress string `json:"contractAddress"`
			To              string `json:"to"`
			TokenID         string `json:"tokenID"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		var result struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, "", err
		}

		return nil, "", errors.New(result.Message)
	}

	nfts := make([]model.NFT, 0, len(result.Transactions))

	for _, transaction := range result.Transactions {
		if !strings.EqualFold(address, transaction.To) {
			continue
		}

		nfts = append(nfts, model.NFT{
			ContractAddress: transaction.ContractAddress,
			TokenID:         transaction.TokenID,
		})
	}

	nextPageToken := ""
	if len(nfts) == pageSize {
		nextPageToken = paginator.GeneratePageToken(strconv.Itoa(page+1), "UNSPECIFIED")
	}

	return nfts, nextPageToken, nil
}
