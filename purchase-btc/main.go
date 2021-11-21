package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	/// APIKEYを取得
	//apiKey, err := getParameter("purchase_btc-api-key")
	//if err != nil {
	//	return getErrorResponse(err.Error()), err
	//}

	/// API SECRETを取得
	apiSecret, err := getParameter("purchase_btc-api-secret")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}
	//
	//t, err := bitflyer.GetTicker(bitflyer.BTCJpy)
	//if err != nil {
	//	return events.APIGatewayProxyResponse{
	//		Body:       "Bad Request",
	//		StatusCode: 400,
	//	}, nil
	//}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("API Secret >>> %+v", apiSecret),
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
