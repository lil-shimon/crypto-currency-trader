package order

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"purchase-btc/utils"
	"strconv"
	"time"
)

const URL = "https://api.bitflyer.com"
const productCodeKey = "product_code"
const bitflyerBtcMinAmount = 0.001

/// 小数点第4位以下を切り捨てるため
const btcPlace = 4.0

// APIClient / APIClientの構造体
type APIClient struct {
	apiKey    string
	apiSecret string
}

// NewAPIClient / APIClientのコンストラクター
/// 引数が多いと可読性が下がるので回避する
func NewAPIClient(apiKey, apiSecret string) *APIClient {
	return &APIClient{apiKey, apiSecret}
}

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

// Order / 新規注文の構造体
type Order struct {
	ProductCode    string  `json:"product_code"`
	ChildOrderType string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Price          float64 `json:"price"`
	Size           float64 `json:"size"`
	MinuteToExpire int     `json:"minute_to_expire"`
	TimeInForce    string  `json:"time_in_force"`
}

// Res OrderRes / 新規注文のレスポンス構造体
type Res struct {
	ChildOrderAcceptanceId string `json:"child_order_acceptance_id"`
}

func GetTicker(ch chan *Ticker, errCh chan error, code ProductCode) {

	/// URL定義
	url := URL + "/v1/ticker"

	/// request
	res, err := utils.CreateHttpRequest("GET", url, nil, map[string]string{productCodeKey: code.String()}, nil)
	if err != nil {
		ch <- nil
		errCh <- err
		return
	}

	/// Ticker
	var t Ticker

	/// レスポンスをTickerに変換する
	err = json.Unmarshal(res, &t)
	if err != nil {
		ch <- nil
		errCh <- err
		return
	}

	/// 返す順番間違えるとブロックされるので注意....
	ch <- &t
	errCh <- nil
}

func (client *APIClient) getHeader(method, path string, body []byte) map[string]string {

	/// UNIXのタイプスタンプを取得
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	/// Signatureを作成
	/// Signature = ACCESS-TIMESTAMP, HTTP method, request path, request bodyを文字列で連結したもの
	/// をAPI SECRETでHMAC=SHA256署名
	text := ts + method + path + string(body)
	mac := hmac.New(sha256.New, []byte(client.apiSecret))
	mac.Write([]byte(text))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       client.apiKey,
		"ACCESS-TIMESTAMP": ts,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

// PlaceOrder / 新規注文を出す関数
func (client *APIClient) PlaceOrder(order *Order) (*Res, error) {

	method := "POST"

	/// API URLを定義
	path := "/v1/me/sendchildorder"
	url := URL + path

	/// POSTするデータを定義
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	/// HEADERを取得
	header := client.getHeader(method, path, data)

	/// HTTP REQUEST
	res, err := utils.CreateHttpRequest(method, url, header, map[string]string{}, data)
	if err != nil {
		return nil, err
	}

	/// responseをorderResに格納
	var orderRes Res
	err = json.Unmarshal(res, &orderRes)
	if err != nil {
		return nil, err
	}

	/// orderResの中身を正しく取得できているかを確認
	if len(orderRes.ChildOrderAcceptanceId) == 0 {
		return nil, errors.New(string(res))
	}

	return &orderRes, nil
}

// GetByLogic / 関数をreturn
/// return func and float64 (高階関数)
func GetByLogic(strategy int) func(float64, *Ticker) (float64, float64) {
	var logic func(float64, *Ticker) (float64, float64)

	/// TODO: ロジックを追加
	switch strategy {
	case 1:
		/// LTPの98.5%の価格で注文を行う
		logic = func(budget float64, t *Ticker) (float64, float64) {
			var buyPrice, buySize float64

			/// TODO: 他の通貨ペア追加時には判定を行う
			/// 最終取引価格98.5%を取得
			buyPrice = utils.RoundDecimal(t.Ltp * 0.985)

			/// 注文数量を取得
			buySize = utils.CalcAmount(buyPrice, budget, bitflyerBtcMinAmount, btcPlace)

			return buyPrice, buySize
		}
		break
	default:

		/// Best Askを注文価格とする
		logic = func(budget float64, t *Ticker) (float64, float64) {
			var buyPrice, buySize float64

			/// BestAskを取得
			buyPrice = utils.RoundDecimal(t.BestAsk)

			/// 注文数量を取得
			buySize = utils.CalcAmount(buyPrice, budget, bitflyerBtcMinAmount, btcPlace)

			return buyPrice, buySize
		}
		break
	}
	return logic
}

func PlaceOrderWithParams(client *APIClient, price, size float64) (*Res, error) {

	/// 注文条件を指定
	o := Order{
		ProductCode:    BTCJpy.String(),
		ChildOrderType: Limit.String(),
		Side:           Buy.String(),
		Price:          price,
		Size:           size,
		MinuteToExpire: 4320, /// 3days
		TimeInForce:    Gtc.String(),
	}

	/// 買い注文を入れる
	orderRes, err := client.PlaceOrder(&o)
	if err != nil {
		return nil, err
	}
	return orderRes, nil
}
