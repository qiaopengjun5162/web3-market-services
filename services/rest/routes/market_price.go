package routes

import (
	"net/http"

	"github.com/ethereum/go-ethereum/log"

	"github.com/qiaopengjun5162/web3-market-services/services/rest/model"
)

func (r *Routes) GetSupportAsset(w http.ResponseWriter, req *http.Request) {
	assetName := req.URL.Query().Get("asset_name")
	consumerToken := req.URL.Query().Get("consumer_token")
	assetRequest := &model.SupportAssetRequest{
		AssetName:     assetName,
		ConsumerToken: consumerToken,
	}
	resp, err := r.srv.GetSupportAsset(assetRequest)
	if err != nil {
		log.Error("GetSupportAsset error: %v", err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, resp); err != nil {
		log.Error("GetSupportAsset jsonResponse error: %v", err)
		return
	}

}

func (r *Routes) GetMarketPrice(w http.ResponseWriter, req *http.Request) {
	assetName := req.URL.Query().Get("asset_name")
	consumerToken := req.URL.Query().Get("consumer_token")
	assetRequest := &model.MarketPriceRequest{
		AssetName:     assetName,
		ConsumerToken: consumerToken,
	}
	resp, err := r.srv.GetMarketPrice(assetRequest)
	if err != nil {
		log.Error("GetMarketPrice error: %v", err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, resp); err != nil {
		log.Error("GetMarketPrice jsonResponse error: %v", err)
		return
	}
}
