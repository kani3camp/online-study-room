package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

type CreateNewNewsParams struct {
	Password      string `json:"password"`
	NewsTitle      string `json:"news_title"`
	TextBody string    `json:"text_body"`
}

type CreateNewNewsResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func CreateNewNews(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("CreateNewNews()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp CreateNewNewsResponseStruct
	body := request.Body
	params := CreateNewNewsParams{}
	_ = json.Unmarshal([]byte(body), &params)

	if params.Password == "" || params.NewsTitle == "" || params.TextBody == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if params.Password != os.Getenv("password") {
		apiResp.Result = ERROR
		apiResp.Message = "invalid password"
	} else {
		err := _CreateNewNews(params.NewsTitle, params.TextBody, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed to create the news."
		} else {
			apiResp.Result = OK
			apiResp.Message = "successfully created news titled " + params.NewsTitle
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(CreateNewNews)
}
