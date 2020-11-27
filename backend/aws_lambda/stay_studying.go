package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

func StayStudying(_ context.Context, request events.APIGatewayWebsocketProxyRequest) events.APIGatewayProxyResponse {
	//connectionID := request.RequestContext.ConnectionID
	
	
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

