package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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

	/// クライアントを生成
	client := order.NewAPIClient(apiKey, apiSecret)

	/// budget定義
	budget := 10000.0

	/// PriceとSizeを取得 (1 = LTP * 0.985の関数を取得)
	/// カリー化
	p, s := order.GetByLogic(1)(budget, t)

	/// 注文
	oRes, err := order.PlaceOrderWithParams(client, p, s)

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
