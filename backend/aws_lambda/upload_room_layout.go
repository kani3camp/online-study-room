package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

type UploadRoomLayoutParams struct {
	RoomLayoutData RoomLayoutStruct `json:"room_layout_data"`
	Password      string `json:"password"`
}

type UploadRoomLayoutResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}


func UploadRoomLayout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("UploadRoomLayout()")
	ctx, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)
	
	var apiResp UploadRoomLayoutResponseStruct
	body := request.Body
	params := UploadRoomLayoutParams{}
	_ = json.Unmarshal([]byte(body), &params)

	roomLayoutData := params.RoomLayoutData
	password := params.Password
	
	if password != os.Getenv("password") {
		apiResp.Result = ERROR
		apiResp.Message = "Invalid password."
	} else {
		customErr := CheckRoomLayoutData(roomLayoutData, client, ctx)
		if customErr.Body != nil {
			apiResp.Result = ERROR
			apiResp.Message = customErr.Body.Error()
		} else {
			err := SaveRoomLayout(roomLayoutData, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
				apiResp.Message = "failed. please review the log in CloudWatch."
			} else {
				apiResp.Result = OK
			}
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(UploadRoomLayout)
}
