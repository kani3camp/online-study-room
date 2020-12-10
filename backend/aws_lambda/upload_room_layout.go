package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"strconv"
)

type UploadRoomLayoutParams struct {
	RoomLayoutData RoomLayoutStruct `json:"room_layout_data"`
	Password      string `json:"password"`
}

type UploadRoomLayoutResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

func CheckRoomLayoutData(roomLayoutData RoomLayoutStruct, client *firestore.Client, ctx context.Context) error {
	var idList []int
	var partitionShapeTypeList []string
	
	if roomLayoutData.RoomId == "" {
		return errors.New("please specify a valid room id")
	} else if isExistRoom , _ := IsExistRoom(roomLayoutData.RoomId, client, ctx); ! isExistRoom {
		return errors.New("any room of that room id doesn't exist")
	} else if currentVersion, _ := CurrentRoomLayoutVersion(roomLayoutData.RoomId, client, ctx); roomLayoutData.Version != 1 + currentVersion {
		return errors.New("please specify a valid version. current version is " + strconv.Itoa(currentVersion))
	} else if roomLayoutData.FontSizeRatio == 0.0 {
		return errors.New("please specify a valid font size ratio")
	} else if roomLayoutData.RoomShape.Height == 0 || roomLayoutData.RoomShape.Width == 0 {
		return errors.New("please specify the room-shape correctly")
	} else {
		if len(roomLayoutData.PartitionShapes) != 0 {
			// PartitionのShapeTypeの重複がないか
			for _, p := range roomLayoutData.PartitionShapes {
				if p.Name == "" || p.Width == 0 || p.Height == 0 {
					return errors.New("please specify partition shapes correctly")
				}	// ここから正常にifを抜けることがある
				for _, other := range partitionShapeTypeList {
					if other == p.Name {
						return errors.New("some partition shape types are duplicated")
					}
				}
				partitionShapeTypeList = append(partitionShapeTypeList, p.Name)
			}
		} else if len(roomLayoutData.Partitions) == 0 {
			return errors.New("please specify partition shapes")
		}
		if len(roomLayoutData.Seats) == 0 {
			return errors.New("please specify at least one seat")
		} else {
			// SeatのIdの重複がないか
			for _, s := range roomLayoutData.Seats {
				for _, other := range idList {
					if other == s.Id {
						return errors.New("some seat ids are duplicated")
					}
				}
				idList = append(idList, s.Id)
			}
			
			if len(roomLayoutData.Partitions) != 0 {
				for _, p := range roomLayoutData.Partitions {
					if p.ShapeType == "" {
						return errors.New("please specify valid shape-type to all shapes")
					}
					isContained := false
					for _, other := range partitionShapeTypeList {
						if other == p.ShapeType {
							isContained = true
						}
					}
					if ! isContained {
						return errors.New("please specify valid shape type")
					}
				}
			}
			return nil
		}
	}
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
		err := CheckRoomLayoutData(roomLayoutData, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = err.Error()
		} else {
			err = SaveRoomLayout(roomLayoutData, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
				apiResp.Message = "failed"
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
