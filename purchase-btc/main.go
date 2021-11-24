package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"math"
	"purchase-btc/order"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	/// Tickerのchannel定義 [order.Ticker]
	tickerChan := make(chan *order.Ticker)

	/// errorのchannel
	errChan := make(chan error)

	/// channel閉じる
	defer close(tickerChan)
	defer close(errChan)

	/// Tickerを取得
	/// 非同期 go routine
	go order.GetTicker(tickerChan, errChan, order.BTCJpy)

	/// 各channelから値を受信
	t := <-tickerChan
	err := <-errChan

	if err != nil {
		return getErrorResponse(err.Error()), nil
	}

	/// APIKEYをSystem Managerから取得
	apiKey, err := getParameter("purchase_btc-api-key")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	/// API SECRETをSystem Managerから取得
	apiSecret, err := getParameter("purchase_btc-api-secret")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	/// 現在価格の95％を買価格定義
	bp := RoundDecimal(t.Ltp * 0.95)

	/// 注文条件を指定
	o := order.Order{
		ProductCode:    order.BTCJpy.String(),
		ChildOrderType: order.Limit.String(),
		Side:           order.Buy.String(),
		Price:          bp,
		Size:           0.001,
		MinuteToExpire: 4320, /// 3days
		TimeInForce:    order.Gtc.String(),
	}

	client := order.NewAPIClient(apiKey, apiSecret)

	/// 買い注文を入れる
	oRes, err := client.PlaceOrder(&o)
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("res >>> %+v", oRes),
		StatusCode: 200,
	}, nil
}

/// System Managerからパラメータを取得する関数
func getParameter(key string) (string, error) {

	/// sessionオブジェクトを作成
	s := session.Must(session.NewSessionWithOptions(session.Options{
		/// ShardConfigEnable -> ~/.aws/config
		/// local動作用
		SharedConfigState: session.SharedConfigEnable,
	}))

	/// System Managerのクライアントを作成
	svc := ssm.New(s, aws.NewConfig().WithRegion("ap-northeast-1"))

	/// Input Structを生成
	params := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	/// System Managerに登録されている値を取得
	res, err := svc.GetParameter(params)
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}

func getErrorResponse(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: 400,
	}
}

func main() {
	lambda.Start(handler)
}
