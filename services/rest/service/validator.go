package service

import (
	"fmt"

	"github.com/qiaopengjun5162/web3-market-services/services/rest/model"
)

type Validator struct {
}

func (v *Validator) validateSupportAssetRequest(req *model.SupportAssetRequest) error {
	if req.AssetName == "" {
		return fmt.Errorf("asset_name is required")
	}
	if req.ConsumerToken == "" {
		return fmt.Errorf("consumer_token is required")
	}
	return nil
}

func (v *Validator) validateMarketPriceRequest(req *model.MarketPriceRequest) error {
	if req.AssetName == "" {
		return fmt.Errorf("asset_name is required")
	}
	if req.ConsumerToken == "" {
		return fmt.Errorf("consumer_token is required")
	}
	return nil
}
