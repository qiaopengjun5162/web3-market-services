package client

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/qiaopengjun5162/web3-market-services/services/rest/model"
)

var errMarketPriceHttpError = errors.New("call market price http error")

type MarketPriceClient interface {
	GetSupportAsset(assetNName string) (bool, error)
	GetMarketPrice(assetName string) (*model.MarketPriceResponse, error)
}

type Client struct {
	client *resty.Client
}

func NewMarketPriceClient(url string) *Client {
	client := resty.New()
	client.SetBaseURL(url)
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		statusCode := resp.StatusCode()
		if statusCode >= 400 {
			method := resp.Request.Method
			baseUrl := resp.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, baseUrl, errMarketPriceHttpError)
		}
		return nil
	})
	return &Client{client: client}
}

func (c *Client) GetSupportAsset(assetName string, consumer_token string) (bool, error) {
	res, err := c.client.R().
		SetQueryParam("asset_name", assetName).
		SetQueryParam("consumer_token", consumer_token).
		//SetQueryParams(map[string]string{
		//	"asset_name": assetName,
		//}).
		SetResult(model.SupportAssetResponse{}).
		Get("/api/v1/get_support_asset")
	if err != nil {
		return false, errors.New("call market price http error")
	}
	ret, ok := res.Result().(*model.SupportAssetResponse)
	if !ok {
		return false, errors.New("call market price http error")
	}
	return ret.IsSupport, nil
}

func (c *Client) GetMarketPrice(assetName string, consumer_token string) (*model.MarketPriceResponse, error) {
	res, err := c.client.R().
		//SetQueryParam("asset_name", assetName).
		SetQueryParams(map[string]string{
			"asset_name": assetName, "consumer_token": consumer_token}).
		SetResult(&model.MarketPriceResponse{}).
		Get("/api/v1/get_market_price")
	if err != nil {
		return nil, errors.New("call market price http error")
	}
	ret, ok := res.Result().(*model.MarketPriceResponse)
	if !ok {
		return nil, errors.New("call market price http error")
	}
	return ret, nil
}
