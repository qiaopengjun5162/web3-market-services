// Package service provides the service layer for the application.
// It contains business logic and interacts with other layers such as data access and presentation.
package service

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/qiaopengjun5162/web3-market-services/database"
	"github.com/qiaopengjun5162/web3-market-services/services/rest/model"
)

type RestService interface {
	GetSupportAsset(*model.SupportAssetRequest) (*model.SupportAssetResponse, error)
	GetMarketPrice(*model.MarketPriceRequest) (*model.MarketPriceResponse, error)
}

type HandleSvc struct {
	v                    *Validator
	marketPriceView      database.MarketPriceViewer
	officialCoinRateView database.OfficialCoinRateViewer
}

func NewHandleSvc(v *Validator, marketPriceView database.MarketPriceViewer, officialCoinRateView database.OfficialCoinRateViewer) *HandleSvc {
	return &HandleSvc{
		v:                    v,
		marketPriceView:      marketPriceView,
		officialCoinRateView: officialCoinRateView,
	}
}

func (h HandleSvc) GetSupportAsset(request *model.SupportAssetRequest) (*model.SupportAssetResponse, error) {
	if err := h.v.validateSupportAssetRequest(request); err != nil {
		log.Error("validateSupportAssetRequest error: %v", err)
		return nil, err
	}

	return &model.SupportAssetResponse{
		ReturnCode: 100,
		Message:    "get asset support success",
		IsSupport:  true,
	}, nil
}

func (h HandleSvc) GetMarketPrice(request *model.MarketPriceRequest) (*model.MarketPriceResponse, error) {
	if err := h.v.validateMarketPriceRequest(request); err != nil {
		log.Error("validateMarketPriceRequest error: %v", err)
		return nil, err
	}
	assetPriceList, err := h.marketPriceView.QueryMarketPriceByAsset(request.AssetName)
	if err != nil {
		log.Error("QueryMarketPriceByAsset error: %v", err)
		return nil, err
	}
	var marketPriceList []*model.MarketPrice
	for _, assetPrice := range assetPriceList {
		mItem := &model.MarketPrice{
			AssetName:   assetPrice.AssetName,
			AssetPrice:  assetPrice.PriceUsdt,
			AssetVolume: assetPrice.Volume,
			AssetRate:   assetPrice.Rate,
		}
		marketPriceList = append(marketPriceList, mItem)
	}

	ocrList, err := h.officialCoinRateView.QueryOfficialCoinRateByAsset(request.AssetName)
	if err != nil {
		log.Error("QueryOfficialCoinRateByAsset error: %v", err)
		return nil, err
	}
	var officialCoinRateList []*model.OfficialCoinRate
	for _, ocr := range ocrList {
		mItem := &model.OfficialCoinRate{
			Name: ocr.AssetName,
			Rate: ocr.Price,
		}
		officialCoinRateList = append(officialCoinRateList, mItem)
	}
	return &model.MarketPriceResponse{
		ReturnCode:           100,
		Message:              "get asset market success",
		MarketPriceList:      marketPriceList,
		OfficialCoinRateList: officialCoinRateList,
	}, nil
}
