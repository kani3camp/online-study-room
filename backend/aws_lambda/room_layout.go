package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"time"
)

type RoomLayoutResponseStruct struct {
	Result         string `json:"result"`
	Message        string `json:"message"`
	RoomLayoutData string `json:"room_layout_data"`
}

func RoomLayout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cloudStorageclient := InitializeHttpFuncWithCloudStorage()
	_, firestoreClient := InitializeEventFuncWithFirestore()
	defer CloseCloudStorageClient(cloudStorageclient)

	var apiResp RoomLayoutResponseStruct

	roomId := request.QueryStringParameters[room_id]

	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		doc, err := firestoreClient.Collection(CONFIG).Doc(ROOM_LAYOUTS_INFO).Get(ctx)
		if err != nil {
			log.Printf("failed to process firestoreClient.Collection(CONFIG).Doc(ROOM_LAYOUTS_INFO).Get(ctx), %v", err)
			apiResp.Result = ERROR
			apiResp.Message = "failed."
		} else {
			var roomLayoutsInfoConfig RoomLayoutsInfoConfigStruct
			_ = doc.DataTo(&roomLayoutsInfoConfig)
			bucketName := roomLayoutsInfoConfig.BucketName
			objectName := roomLayoutsInfoConfig.ObjectPaths[roomId]
			if objectName == "" {
				apiResp.Result = ERROR
				apiResp.Message = "invalid room id."
			} else {
				rc, err := cloudStorageclient.Bucket(bucketName).Object(objectName).NewReader(ctx)
				if err != nil {
					log.Printf("Object(%q).NewReader: %v", objectName, err)
					apiResp.Result = ERROR
					apiResp.Message = "failed to read data file."
				} else {
					defer rc.Close()

					data, err := ioutil.ReadAll(rc)
					if err != nil {
						log.Printf("ioutil.ReadAll: %v", err)
						apiResp.Result = ERROR
						apiResp.Message = "failed to read data file."
					} else {
						apiResp.RoomLayoutData = string(data)
						apiResp.Result = OK
					}
				}
			}
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(RoomLayout)
}
