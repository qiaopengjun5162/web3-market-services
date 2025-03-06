package client

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSupportAsset(t *testing.T) {
	client := NewMarketPriceClient("http://127.0.0.1:8080")
	result, err := client.GetSupportAsset("all", "consumer_token")
	if err != nil {
		t.Error("GetSupportAsset error:", err)
		return
	}
	t.Log("GetSupportAsset result:", result)
	fmt.Printf("GetSupportAsset result: %v\n", result)
}

func TestMarketPrice(t *testing.T) {
	client := NewMarketPriceClient("http://127.0.0.1:8080")
	result, err := client.GetMarketPrice("all", "consumer_token")
	if err != nil {
		t.Error("GetMarketPrice error:", err)
		return
	}
	res, _ := json.Marshal(result)
	t.Log("GetMarketPrice result:", string(res))
	fmt.Printf("GetMarketPrice result: %v\n", string(res))
}
