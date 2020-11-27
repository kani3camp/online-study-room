package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func StartStudying(_ context.Context, request events.APIGatewayWebsocketProxyRequest) events.APIGatewayProxyResponse {
	//connectionID := request.RequestContext.ConnectionID

	return 
}

func main()  {
	lambda.Start(StartStudying)
}