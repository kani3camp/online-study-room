package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

func StartStudying(_ context.Context, request events.APIGatewayWebsocketProxyRequest) events.APIGatewayProxyResponse {
	//connectionID := request.RequestContext.ConnectionID
	

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

func main()  {
	lambda.Start(StartStudying)
}