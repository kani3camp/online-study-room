package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RoomLayoutResponseStruct struct {
	Result         string `json:"result"`
	Message        string `json:"message"`
	RoomLayoutData RoomLayoutStruct `json:"room_layout_data"`
}

func RoomLayout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp RoomLayoutResponseStruct

	roomId := request.QueryStringParameters[room_id]

	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		roomLayout, err := RetrieveRoomLayout(roomId, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed."
		} else {
			apiResp.Result = OK
			apiResp.RoomLayoutData = roomLayout
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(RoomLayout)
}

//func RoomLayout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	ctx, cloudStorageClient := InitializeHttpFuncWithCloudStorage()
//	_, firestoreClient := InitializeEventFuncWithFirestore()
//	defer CloseCloudStorageClient(cloudStorageClient)
//
//	var apiResp RoomLayoutResponseStruct
//
//	roomId := request.QueryStringParameters[room_id]
//
//	if roomId == "" {
//		apiResp.Result = ERROR
//		apiResp.Message = InvalidParams
//	} else {
//		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
//		defer cancel()
//
//		doc, err := firestoreClient.Collection(CONFIG).Doc(ROOM_LAYOUTS_INFO).Get(ctx)
//		if err != nil {
//			log.Printf("failed to process firestoreClient.Collection(CONFIG).Doc(ROOM_LAYOUTS_INFO).Get(ctx), %v", err)
//			apiResp.Result = ERROR
//			apiResp.Message = "failed."
//		} else {
//			var roomLayoutsInfoConfig RoomLayoutsInfoConfigStruct
//			_ = doc.DataTo(&roomLayoutsInfoConfig)
//			bucketName := roomLayoutsInfoConfig.BucketName
//			objectName := roomLayoutsInfoConfig.ObjectPaths[roomId]
//			if objectName == "" {
//				apiResp.Result = ERROR
//				apiResp.Message = "invalid room id."
//			} else {
//				rc, err := cloudStorageClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
//				if err != nil {
//					log.Printf("Object(%q).NewReader: %v", objectName, err)
//					apiResp.Result = ERROR
//					apiResp.Message = "failed to read data file."
//				} else {
//					defer rc.Close()
//
//					data, err := ioutil.ReadAll(rc)
//					if err != nil {
//						log.Printf("ioutil.ReadAll: %v", err)
//						apiResp.Result = ERROR
//						apiResp.Message = "failed to read data file."
//					} else {
//						apiResp.RoomLayoutData = string(data)
//						apiResp.Result = OK
//					}
//				}
//			}
//		}
//	}
//
//	bytes, _ := json.Marshal(apiResp)
//	return Response(bytes)
//}
