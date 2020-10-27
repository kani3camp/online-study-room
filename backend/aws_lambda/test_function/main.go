package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type TestFunctionParams struct {
	Name string `json:"name"`
	Age string `json:"age"`
}

type Response struct {
	RequestMethod  string `json:"RequestMethod"`
	RequestBody    string `json:"RequestBody"`
	PathParameter  string `json:"PathParameter"`
	QueryParameter string `json:"QueryParameter"`
}


func HandleLambdaEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// httpリクエストの情報を取得
	method := request.HTTPMethod
	body := request.Body
	pathParam := request.PathParameters["pathparam"]
	queryParam := request.QueryStringParameters["queryparam"]
	
	params := TestFunctionParams{}
	_ = json.Unmarshal([]byte(body), &params)
	
	// レスポンスとして返すjson文字列を作る
	res := Response{
		RequestMethod:  method,
		RequestBody:    body,
		PathParameter:  pathParam,
		QueryParameter: queryParam,
	}
	jsonBytes, _ := json.Marshal(res)
	
	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		//Headers: map[string]string{
		//	"Content-Type": "application/json",
		//},
	}, nil
}

func main()  {
	lambda.Start(HandleLambdaEvent)
}
