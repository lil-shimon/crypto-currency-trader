package bitflyer

// ProductCode / 通貨ペア定義
type ProductCode int

const (
	BTCJpy ProductCode = iota
	EthJpy
	FxBtcJpy
	EthBtc
	BchBtc
)

// OrderType / OrderType
/// Limit = 指値
/// Market = 成行
type OrderType int

const (
	Limit OrderType = iota
	Market
)

// Side / 買いOR売り
type Side int

const (
	Buy Side = iota
	Sell
)

// TimeInForce / 執行数量条件
/// Gtc = Good Til Canceled 注文が約定するかキャンセルされるまで有効
/// Ioc = Immediate or Cancel 指定した価格かそれよりも有効な価格で即時に一部あるいは全部を約定させ、約定しなかった注文数量をキャンセル
/// Fok = Fill or Kill 発注の全数量が即座に約定しない場合当該注文をキャンセル
type TimeInForce int

const (
	Gtc TimeInForce = iota
	Ioc
	Fok
)

/// プロダクトコードをStringに型変換
func (code ProductCode) String() string {
	switch code {
	case BTCJpy:
		return "BTC_JPY"
	case EthJpy:
		return "ETH_JPY"
	case FxBtcJpy:
		return "FX_BTC_JPY"
	case EthBtc:
		return "ETH_BTC"
	case BchBtc:
		return "BCH_BTC"
	default:
		return "BTC_JPY"
	}
}

/// オーダータイプをStringに変換
func (ot OrderType) String() string {
	switch ot {
	case Limit:
		return "LIMIT"
	case Market:
		return "MARKET"
	default:
		return "LIMIT"
	}
}

/// 売買を文字列に変換
func (s Side) String() string {
	switch s {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return "BUY"
	}
}

/// タイムインフォースを文字列に変換
func (tif TimeInForce) String() string {
	switch tif {
	case Gtc:
		return "GTC"
	case Ioc:
		return "IOC"
	case Fok:
		return "FOK"
	default:
		return "GTC"
	}
}
