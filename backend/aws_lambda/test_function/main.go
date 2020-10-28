package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type TestFunctionParams struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

type Response struct {
	RequestMethod  string `json:"RequestMethod"`
	RequestBody    string `json:"RequestBody"`
	PathParameter  string `json:"PathParameter"`
	QueryParameter string `json:"QueryParameter"`
	Name string `json:"name"`
	Age int `json:"age"`
}


func HandleLambdaEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// httpリクエストの情報を取得
	body := request.Body	// POSTの場合
	queryParam := request.QueryStringParameters["name"]	// GETの場合。"name"は適宜keyに変える

	params := TestFunctionParams{}
	_ = json.Unmarshal([]byte(body), &params)

	// レスポンスとして返すjson文字列を作る
	res := Response{
		RequestBody:    body,
		QueryParameter: queryParam,
		Name: params.Name,
		Age: params.Age,
	}
	jsonBytes, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		StatusCode: 200,	// これないとInternal Server Errorになる
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main()  {
	lambda.Start(HandleLambdaEvent)
}
