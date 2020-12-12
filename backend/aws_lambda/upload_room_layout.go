package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
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

func CheckRoomLayoutData(roomLayoutData RoomLayoutStruct, client *firestore.Client, ctx context.Context) CustomError {
	var idList []int
	var partitionShapeTypeList []string
	
	if roomLayoutData.RoomId == "" {
		return InvalidRoomLayout.New("please specify a valid room id")
	} else if isExistRoom , _ := IsExistRoom(roomLayoutData.RoomId, client, ctx); ! isExistRoom {
		return InvalidRoomLayout.New("any room of that room id doesn't exist")
	} else if currentVersion, _ := CurrentRoomLayoutVersion(roomLayoutData.RoomId, client, ctx); roomLayoutData.Version != 1 + currentVersion {
		return InvalidRoomLayout.New("please specify a incremented version. latest version is " + strconv.Itoa(currentVersion))
	} else if roomLayoutData.FontSizeRatio == 0.0 {
		return InvalidRoomLayout.New("please specify a valid font size ratio")
	} else if roomLayoutData.RoomShape.Height == 0 || roomLayoutData.RoomShape.Width == 0 {
		return InvalidRoomLayout.New("please specify the room-shape correctly")
	}
	// 横長のレイアウトのみ可
	if roomLayoutData.RoomShape.Width < roomLayoutData.RoomShape.Height {
		return InvalidRoomLayout.New("please make room width larger than room height")
	}

	if len(roomLayoutData.PartitionShapes) > 0 {
		// PartitionのShapeTypeの重複がないか
		for _, p := range roomLayoutData.PartitionShapes {
			if p.Name == "" || p.Width == 0 || p.Height == 0 {
				return InvalidRoomLayout.New("please specify partition shapes correctly")
			}	// ここから正常にifを抜けることがある
			for _, other := range partitionShapeTypeList {
				if other == p.Name {
					return InvalidRoomLayout.New("some partition shape types are duplicated")
				}
			}
			partitionShapeTypeList = append(partitionShapeTypeList, p.Name)
		}
	}
	if len(roomLayoutData.Seats) == 0 {
		return InvalidRoomLayout.New("please specify at least one seat")
	}
	// SeatのIdの重複がないか
	for _, s := range roomLayoutData.Seats {
		for _, other := range idList {
			if other == s.Id {
				return InvalidRoomLayout.New("some seat ids are duplicated")
			}
		}
		idList = append(idList, s.Id)
	}
	// 仕切り
	for _, p := range roomLayoutData.Partitions {
		if p.ShapeType == "" {
			return InvalidRoomLayout.New("please specify valid shape-type to all shapes")
		}
		// 仕切りのShapeTypeに有効なものが指定されているか
		isContained := false
		for _, other := range partitionShapeTypeList {
			if other == p.ShapeType {
				isContained = true
			}
		}
		if ! isContained {
			return InvalidRoomLayout.New("please specify valid shape type, at partition id = " + strconv.Itoa(p.Id))
		}
	}
	return CustomError{Body: nil}
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
			// todo 前後で座席に変更があった場合、現在そのルームにいる人を強制的に退室させる
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
