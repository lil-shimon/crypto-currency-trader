package main

import (
	"bytes"
	csv "encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const ApiUrl = "https://coincheck.com"

type CurrencyType int

type Res struct {
	Rate   string `json:"rate"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

const (
	Btc CurrencyType = iota
	Eth
	Mona
	Plt
)

func (ct CurrencyType) String() string {
	switch ct {
	case Btc:
		return "btc_jpy"
	case Eth:
		return "eth_jpy"
	case Mona:
		return "mona_jpy"
	case Plt:
		return "plt_jpy"
	default:
		return "btc_jpy"
	}
}

/// save chart data to csv
func fetchChartByCt(ct string) {
	url := ApiUrl + "/api/exchange/orders/rate"

	res, err := CreateHttpRequest("GET", url, nil, map[string]string{"pair": ct, "order_type": "buy", "amount": "1"}, nil)

	if err != nil {
		panic(err)
	}

	var r Res
	err = json.Unmarshal(res, &r)
	if err != nil {
		panic(err)
	}

	var now = time.Now()
	record := readFromCsv()
	record = append(record, []string{r.Amount, r.Price, r.Rate, timeToStr(now)})
	writeToCsv(record)
	print(r.Amount, r.Price, r.Rate, timeToStr(now), "\n")
}

func timeToStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func writeToCsv(data [][]string) {
	file, err := os.Create("chart.csv")
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(file)
	w.WriteAll(data)

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func readFromCsv() [][]string {
	file, err := os.ReadFile("chart.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(file))
	record, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return record
}

func CreateHttpRequest(method, url string, header, query map[string]string, data []byte) ([]byte, error) {

	/// GET POSTしか受け付けない
	if method != "GET" && method != "POST" {
		return nil, errors.New("method was neither GET or POST")
	}

	/// HTTP requestを作成
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	/// Query Parameterを作成
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}

	/// エンコード
	req.URL.RawQuery = q.Encode()

	/// Headerの処理
	for key, value := range header {
		req.Header.Add(key, value)
	}

	/// HTTP Client
	httpC := &http.Client{}

	/// requestを投げる
	res, err := httpC.Do(req)
	if err != nil {
		return nil, err
	}

	/// 開放処理
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return // ignore error
		}
	}(res.Body)

	/// resの中身を取得
	bd, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bd, nil
}

func main() {
	for range time.Tick(time.Second * 1) {
		fetchChartByCt(Btc.String())
	}
}
