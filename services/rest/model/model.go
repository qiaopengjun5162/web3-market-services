package model

type SupportAssetRequest struct {
	AssetName     string `json:"asset_name"`
	ConsumerToken string `json:"consumer_token"`
}

type SupportAssetResponse struct {
	IsSupport  bool   `json:"is_support"`
	Message    string `json:"message"`
	ReturnCode uint64 `json:"return_code"`
}

type OfficialCoinRate struct {
	Name string `json:"name"`
	Rate string `json:"rate"`
}

type MarketPrice struct {
	AssetName   string `json:"asset_name"`
	AssetPrice  string `json:"asset_price"`
	AssetVolume string `json:"asset_volume"`
	AssetRate   string `json:"asset_rate"`
}

type MarketPriceRequest struct {
	AssetName     string `json:"asset_name"`
	ConsumerToken string `json:"consumer_token"`
}

type MarketPriceResponse struct {
	MarketPriceList      []*MarketPrice      `json:"market_price_list"`
	Message              string              `json:"message"`
	OfficialCoinRateList []*OfficialCoinRate `json:"official_coin_rate_list"`
	ReturnCode           uint64              `json:"return_code"`
}
