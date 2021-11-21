package bitflyer

import (
	"encoding/json"
	"purchase-btc/utils"
)

const URL = "https://api.bitflyer.com"
const productCodeKey = "product_code"

// Ticker / Tickerを格納する構造体
type Ticker struct {
	ProductCode     string  `json:"product-code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func GetTicker(code ProductCode) (*Ticker, error) {

	/// URL定義
	url := URL + "/v1/ticker"

	/// request
	res, err := utils.CreateHttpRequest("GET", url, nil, map[string]string{productCodeKey: code.String()}, nil)
	if err != nil {
		return nil, err
	}

	/// Ticker
	var t Ticker

	/// レスポンスをTickerに変換する
	err = json.Unmarshal(res, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
