package fastforex

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romutchio/crypto-calculator/internal/entity"
	"github.com/valyala/fasthttp"
)

const FetchOneURL = "fetch-one?from=%s&to=%s&api_key=%s"
const FetchMultiURL = "fetch-multi?from=%s&to=%s&api_key=%s"
const CryptoFetchPricesURL = "/crypto/fetch-prices?pairs=%s&api_key=%s"

func (c *Client) FetchOne(from string, to string) (*entity.FetchOneResponse, error) {
	url := fmt.Sprintf(FetchOneURL, from, to, c.config.Token)
	resp, err := c.request(url, fasthttp.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("request fail: %w", err)
	}
	result := &entity.FetchOneResponse{}
	if err = json.Unmarshal(resp, result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return result, nil
}

func (c *Client) FetchMulti(from string, tos ...string) (*entity.FetchMultiResponse, error) {
	to := strings.Join(tos, ",")
	url := fmt.Sprintf(FetchMultiURL, from, to, c.config.Token)
	resp, err := c.request(url, fasthttp.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("request fail: %w", err)
	}

	result := &entity.FetchMultiResponse{}
	if err = json.Unmarshal(resp, result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return result, nil
}

func (c *Client) CryptoFetchPrices(pairs ...string) (*entity.CryptoFetchPrices, error) {
	pair := strings.Join(pairs, ",")
	url := fmt.Sprintf(CryptoFetchPricesURL, pair, c.config.Token)
	resp, err := c.request(url, fasthttp.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("request fail: %w", err)
	}

	result := &entity.CryptoFetchPrices{}
	if err = json.Unmarshal(resp, result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return result, nil
}
