package bitflyer

type ProductCode int

/// 通貨ペア定義
const (
	BTCJpy ProductCode = iota
	EthJpy
	FxBtcJpy
	EthBtc
	BchBtc
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
