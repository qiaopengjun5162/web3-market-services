package grpc

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"

	"github.com/qiaopengjun5162/web3-market-services/database"
	"github.com/qiaopengjun5162/web3-market-services/protobuf/market"
)

func (ms *MarketRpcServices) GetSupportAsset(ctx context.Context, req *market.SupportAssetRequest) (*market.SupportAssetResponse, error) {
	return &market.SupportAssetResponse{
		ReturnCode: 100,
		Message:    "support this asset",
		IsSupport:  true,
	}, nil
}

func (ms *MarketRpcServices) GetMarketPrice(ctx context.Context, req *market.MarketPriceRequest) (*market.MarketPriceResponse, error) {
	var marketPriceWrite []*database.MarketPrice
	var coinRateWrite []*database.OfficialCoinRate

	marketPriceBTC := &database.MarketPrice{
		GUID:      uuid.New(),
		AssetName: "BTC",
		PriceUsdt: "80000",
		Volume:    "800000000",
		Rate:      "10",
		Timestamp: uint64(time.Now().Unix()),
	}
	marketPriceWrite = append(marketPriceWrite, marketPriceBTC)

	coinRateItem := &database.OfficialCoinRate{
		GUID:      uuid.New(),
		AssetName: "Cny",
		BaseAsset: "USD",
		Price:     "7.3",
		Timestamp: uint64(time.Now().Unix()),
	}

	coinRateWrite = append(coinRateWrite, coinRateItem)

	err := ms.db.MarketPrice.StoreMarketPrice(marketPriceWrite)
	if err != nil {
		log.Error("store market price fail", "err", err)
		return nil, err
	}

	err = ms.db.OfficialCoinRate.StoreOfficialCoinRate(coinRateWrite)
	if err != nil {
		log.Error("store coin rate fail", "err", err)
		return nil, err
	}

	var marketPriceReturns []*market.MarketPrice
	var coinRateReturns []*market.OfficialCoinRate

	assetPriceList, err := ms.db.MarketPrice.QueryMarketPriceByAsset("all")
	if err != nil {
		return nil, err
	}
	for _, value := range assetPriceList {
		mItem := &market.MarketPrice{
			AssetName:   value.AssetName,
			AssetPrice:  value.PriceUsdt,
			AssetVolume: value.Volume,
			AssetRate:   value.Rate,
		}
		marketPriceReturns = append(marketPriceReturns, mItem)
	}
	coinRateList, err := ms.db.OfficialCoinRate.QueryOfficialCoinRateByAsset("all")
	if err != nil {
		return nil, err
	}
	for _, value := range coinRateList {
		mItem := &market.OfficialCoinRate{
			Name: value.AssetName,
			Rate: value.Price,
		}
		coinRateReturns = append(coinRateReturns, mItem)
	}

	return &market.MarketPriceResponse{
		ReturnCode:       100,
		Message:          "get asset market success",
		MarketPrice:      marketPriceReturns,
		OfficialCoinRate: coinRateReturns,
	}, nil
}
