package utils

import "math"

// RoundDecimal / 四捨五入関数
func RoundDecimal(num float64) float64 {
	return math.Round(num)
}

/// 少数を切り上げる関数
func roundUp(num, places float64) float64 {

	/// 10のplaces乗を行う
	sft := math.Pow(10, places)
	return RoundDecimal(num*sft) / sft
}

// CalcAmount / 予算をもとに購入数量を決める関数
func CalcAmount(price, budget, minAmount, places float64) float64 {

	/// 購入数量を予算と価格から算出
	a := roundUp(budget/price, places)

	/// 購入数量が最小取引価格を下回る場合は最小取引価格を返す
	if a < minAmount {
		return minAmount
	} else {
		return a
	}
}
