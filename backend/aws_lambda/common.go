package main

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)


const (
	ROOMS = "rooms"
	USERS = "users"
	HISTORY = "history"
	CONFIG = "config"
	API = "api"
	ApiGateway = "api-gateway"
	RoomLayoutsInfo = "room-layouts-info"
	RoomLayouts = "room-layouts"
	NEWS = "news"

	ProjectId = "online-study-space"
	TimeLimit = 180 // 秒

	UserId = "user_id"
	RoomId = "room_id"
	IdToken = "id_token"

	RoomDoesNotExist = "room does not exist"
	InvalidParams = "invalid parameters"
	InvalidUser = "invalid user"
	InvalidValue = "invalid value"
	OK = "ok"
	ERROR = "error"
	UserAuthFailed = "user authentication failed"
	UserDoesNotExist = "user does not exist"
	Failed = "failed"

	EnterActivity = "entering"
	LeaveActivity = "leaving"
	NewRoomLayoutActivity = "new-room-layout"

)


type ErrorType uint
const (
	Unknown ErrorType = iota
	SeatNotAvailable
	UserNotInTheRoom
	UserNotInAnyRoom
	NoSuchUserExists
	RoomNotExist
)

// ErrorTypeを取得すること(Type()関数が使えること)に特化した型であり、ここではCustomError型の構造体をこの型として代入（もしくは型アサーション）できる
// ただ、今回はこのTypeGetterを使う型アサーションは、CustomErrorにすればよさげ？？わざわざこれを定義しなくてもよさそうだが...
//type TypeGetter interface {
//	Type() ErrorType
//}
type CustomError struct {
	ErrorType ErrorType
	Body      error
}
func (et ErrorType) New(message string) CustomError {
	return CustomError{ErrorType: et, Body: errors.New(message)}
}
func (et ErrorType) Wrap(err error, message string) CustomError {
	return CustomError{ErrorType: et, Body: errors.Wrap(err, message)}
}
//func (e CustomError) Error() string {
//	return e.Body.Error()
//}
//func (e CustomError) Type() ErrorType {
//	return e.ErrorType
//}
//func Wrap(err error, message string) CustomError {
//	we := errors.Wrap(err, message)
//	if customError, ok := err.(TypeGetter); ok {
//		return CustomError{ErrorType: customError.Type(), Body: we}
//	}
//	return CustomError{ErrorType: Unknown, Body: we}
//}
//func Cause(err error) error {
//	return errors.Cause(err)
//}
//func GetErrorType(err error) ErrorType {
//	for {
//		if e, ok := err.(TypeGetter); ok {
//			return e.Type()
//		}
//		break
//	}
//	return Unknown
//}


type RoomStruct struct {
	RoomId string         `json:"room_id"`
	Body   RoomBodyStruct `json:"room_body"`
}

type RoomBodyStruct struct {
	Created       time.Time `firestore:"created" json:"created"`
	Name          string    `firestore:"name" json:"name"`
	Users         []UserSeatSetStruct `firestore:"users" json:"users"`
	Type          string    `firestore:"type" json:"type"`
	ThemeColorHex string    `firestore:"theme-color-hex" json:"theme_color_hex"`
}

type UserSeatSetStruct struct {
	SeatId int `firestore:"seat-id" json:"seat_id"`
	UserId string `firestore:"user-id" json:"user_id"`
}

type UserStruct struct {
	UserId      string         `json:"user_id"`
	DisplayName string         `json:"display_name"`
	Body        UserBodyStruct `json:"user_body"`
}

type UserBodyStruct struct {
	In               string    `firestore:"in" json:"in"`
	SeatId int `firestore:"seat-id" json:"seat_id"`
	ConnectionId string `firestore:"connection-id" json:"connection_id"`
	LastAccess       time.Time `firestore:"last-access" json:"last_access"`
	LastEntered      time.Time `firestore:"last-entered" json:"last_entered"`
	LastExited       time.Time `firestore:"last-exited" json:"last_exited"`
	LastStudied      time.Time `firestore:"last-studied" json:"last_studied"`
	Online           bool      `firestore:"online" json:"online"`
	Status           string    `firestore:"status" json:"status"`
	RegistrationDate time.Time `firestore:"registration-date" json:"registration_date"`
	TotalStudyTime   int64     `firestore:"total-study-time" json:"total_study_time"`
	TotalBreakTime   int64     `firestore:"total-break-time" json:"total_break_time"`
}

type NewsStruct struct {
	NewsId   string         `json:"news_id"`
	NewsBody NewsBodyStruct `json:"news_body"`
}

type NewsBodyStruct struct {
	Created  time.Time `firestore:"created" json:"created"`
	Updated  time.Time `firestore:"updated" json:"updated"`
	Title    string    `firestore:"title" json:"title"`
	TextBody string    `firestore:"text-body" json:"text_body"`
}

//type EnteringAndLeavingHistoryStruct struct {
//	Activity string    `firestore:"activity"`
//	Room     string    `firestore:"room"`
//	Date     time.Time `firestore:"date"`
//	UserId   string    `firestore:"user-id"`
//}

type RoomLayoutsInfoConfigStruct struct {
	BucketName  string            `firestore:"bucket-name"`
	ObjectPaths map[string]string `firestore:"object-paths"`
}

type RoomLayoutStruct struct {
	RoomId string `json:"room_id" firestore:"room-id"`
	Version int `json:"version" firestore:"version"`
	FontSizeRatio float32 `json:"font_size_ratio" firestore:"font-size-ratio"`
	RoomShape struct {
		Height int `json:"height" firestore:"height"`
		Width int `json:"width" firestore:"width"`
	} `json:"room_shape" firestore:"room-shape"`
	SeatShape struct {
		Height int `json:"height" firestore:"height"`
		Width int `json:"width" firestore:"width"`
	} `json:"seat_shape" firestore:"seat-shape"`
	PartitionShapes []struct{
		Name string `json:"name" firestore:"name"`
		Width int `json:"width" firestore:"width"`
		Height int `json:"height" firestore:"height"`
	} `json:"partition_shapes" firestore:"partition-shapes"`
	Seats []struct {
		Id int `json:"id" firestore:"id"`
		X int `json:"x" firestore:"x"`
		Y int `json:"y" firestore:"y"`
		IsVacant bool `json:"is_vacant"`	// これはfirestoreには保存したくない
	} `json:"seats" firestore:"seats"`
	Partitions []struct {
		Id int `json:"id" firestore:"id"`
		X int `json:"x" firestore:"x"`
		Y int `json:"y" firestore:"y"`
		ShapeType string `json:"shape_type" firestore:"shape-type"`
	} `json:"partitions" firestore:"partitions"`
}

type ApiGatewayConfigStruct struct {
	// https://docs.aws.amazon.com/ja_jp/apigateway/latest/developerguide/apigateway-how-to-call-websocket-api-connections.html
	WebsocketConnectionBaseUrl string `json:"websocket_connection_base_url" firestore:"websocket-connection-base-url"`
}

type EnterLeaveHistory struct {
	Activity string `firestore:"activity"`
	RoomId string `firestore:"room-id"`
	UserId string `firestore:"user-id"`
	Date time.Time `firestore:"date"`
}


func InitializeHttpFuncWithFirestore() (context.Context, *firestore.Client) {
	log.Println("InitializeHttpFuncWithFirestore()")
	return InitializeEventFuncWithFirestore()
}

func InitializeEventFuncWithFirestore() (context.Context, *firestore.Client) {
	log.Println("InitializeEventFuncWithFirestore()")
	ctx := context.Background()
	client, _ := InitializeFirestoreClient(ctx)
	return ctx, client
}

func InitializeHttpFuncWithCloudStorage() (context.Context, *storage.Client) {
	log.Println("InitializeHttpFuncWithCloudStorage()")
	return InitializeEventFuncWithCloudStorage()
}

func InitializeEventFuncWithCloudStorage() (context.Context, *storage.Client) {
	log.Println("InitializeEventFuncWithCloudStorage()")
	ctx := context.Background()
	client, _ := InitializeCloudStorageClient(ctx)
	return ctx, client
}

func InitializeFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	log.Println("InitializeFirestoreClient()")
	var client *firestore.Client
	var err1, err2 error
	awsCredential, err1 := RetrieveFirebaseCredentialInBytes()
	if err1 != nil {
		client, err2 = firestore.NewClient(ctx, ProjectId)
	} else {
		sa := option.WithCredentialsJSON(awsCredential)
		client, err2 = firestore.NewClient(ctx, ProjectId, sa)
	}
	if err2 != nil {
		log.Println(err2)
		return nil, err2
	}
	return client, nil
}

func InitializeCloudStorageClient(ctx context.Context) (*storage.Client, error) {
	log.Println("InitializeCloudStorageClient()")
	var client *storage.Client
	var err1, err2 error
	awsCredential, err1 := RetrieveCloudStorageCredentialInBytes()
	if err1 != nil {
		client, err2 = storage.NewClient(ctx)
	} else {
		sa := option.WithCredentialsJSON(awsCredential)
		client, err2 = storage.NewClient(ctx, sa)
	}
	if err2 != nil {
		log.Println(err2)
		return nil, err2
	}
	return client, nil
}

func InitializeFirebaseApp(ctx context.Context) (*firebase.App, error) {
	log.Println("InitializeFirebaseApp()")
	var app *firebase.App
	var err1, err2 error
	awsCredential, err1 := RetrieveFirebaseCredentialInBytes()
	if err1 != nil {
		app, err2 = firebase.NewApp(ctx, nil)
	} else {
		sa := option.WithCredentialsJSON(awsCredential)
		app, err2 = firebase.NewApp(ctx, nil, sa)
	}
	if err2 != nil {
		log.Println("failed to initialize firebase.App.")
		log.Println(err2)
		return nil, err2
	}
	return app, nil
}

func InitializeFirebaseAuthClient(ctx context.Context) (*auth.Client, error) {
	log.Println("InitializeFirebaseAuthClient()")
	app, _ := InitializeFirebaseApp(ctx)
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return authClient, err
}

func CloseFirestoreClient(client *firestore.Client) {
	log.Println("CloseFirestoreClient()")
	err := client.Close()
	if err != nil {
		log.Println("failed to close firestore client.")
	} else {
	}
}

func CloseCloudStorageClient(client *storage.Client) {
	log.Println("CloseCloudStorageClient()")
	err := client.Close()
	if err != nil {
		log.Println("failed to close cloud storage client.")
	} else {
	}
}

func RetrieveFirebaseCredentialInBytes() ([]byte, error) {
	log.Println("RetrieveFirebaseCredentialInBytes()")
	secretName := "firestore-service-account"
	region := "ap-northeast-1"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return []byte(secretString), nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		_len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			//return
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:_len])
		return []byte(decodedBinarySecret), nil
	}
}

func RetrieveCloudStorageCredentialInBytes() ([]byte, error) {
	log.Println("RetrieveCloudStorageCredentialInBytes()")
	secretName := "cloudstorage-service-account"
	region := "ap-northeast-1"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return []byte(secretString), nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		_len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return nil, err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:_len])
		return []byte(decodedBinarySecret), nil
	}

	// Your code goes here.
}

func IsUserVerified(userId string, idToken string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsUserVerified()")
	authClient, _ := InitializeFirebaseAuthClient(ctx)
	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return false, err
	} else if userId != token.UID {
		return false, nil
	} else if isInUsers, _ := IsInUsers(userId, client, ctx); !isInUsers {
		return false, nil
	}
	return true, nil
}

func IsExistRoom(roomId string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsExistRoom()")
	log.Println("roomId : " + roomId + ".")
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return roomDoc.Exists(), nil
}

func IsInRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) (bool, int, error) {
	log.Println("IsInRoom()")
	log.Println("IsInRoom() is running. roomId : " + roomId + ". userId : " + userId + ".")
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, 0, err
	}
	var room RoomBodyStruct
	err = roomDoc.DataTo(&room)
	if err != nil {
		log.Println(err)
		return false, 0, err
	}
	users := room.Users
	for _, u := range users {
		if u.UserId == userId {
			return true, u.SeatId, nil
		}
	}
	return false, 0, nil
}

func IsInUsers(userId string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsInUsers()")
	log.Println("\t userId : " + userId + ".")
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		return false, err
	}
	return userDoc.Exists(), nil
}

func IsOnline(userId string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsOnline()")
	log.Println("IsOnline() is running. userId : " + userId + ".")
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, err // エラーの場合もfalseを返すので注意
	} else {
		return userDoc.Data()["online"].(bool), nil
	}
}

func LeaveRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) CustomError {
	log.Println("LeaveRoom()")
	log.Println("roomId : " + roomId + ". user : " + userId + ".")
	if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
		errString := RoomDoesNotExist
		log.Println(errString)
		return RoomNotExist.New(errString)
	} else if isInRoom, seatId, _ := IsInRoom(roomId, userId, client, ctx); !isInRoom {
		log.Println("you are not in the room.")
		return CustomError{Body: nil}
	} else {
		// userのconnection idを削除
		_ = SetConnectionId(userId, "", client, ctx)

		// 退室処理
		userSeatSet := UserSeatSetStruct{
			UserId: userId,
			SeatId: seatId,
		}
		_, err := client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
			"users": firestore.ArrayRemove(userSeatSet),
		}, firestore.MergeAll)
		if err != nil {
			log.Println("failed to remove user from room")
		}

		// user status を更新
		_, err = client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
			"online":       false,
			"in":           "",
			"seat-id": 0,
			"last-studied": time.Now(),
		}, firestore.MergeAll)
		if err != nil {
			log.Println("Failed to update user info of " + userId)
		}
		_ = RecordLastAccess(userId, client, ctx)
		_ = RecordExitedTime(userId, client, ctx)
		_ = RecordHistory(EnterLeaveHistory{
			Activity: LeaveActivity,
			RoomId:   roomId,
			UserId:  userId,
			Date:    time.Now(),
		}, client, ctx)
		defer UpdateTotalTime(userId, roomId, time.Now(), client, ctx)

		roomBody, _ := RetrieveRoomInfo(roomId, client, ctx)
		authClient, _ := InitializeFirebaseAuthClient(ctx)
		user, _ := authClient.GetUser(ctx, userId)
		defer SendLiveChatMessage(user.DisplayName+"さんが"+roomBody.Name+"ルームを出ました。", client, ctx)
	}
	return CustomError{Body: nil}
}

func SetConnectionId(userId string, connectionId string, client *firestore.Client, ctx context.Context) error {
	log.Println("SetConnectionId()")
	// inやonlineはいじらない
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"connection-id": connectionId,
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func RetrieveConnectionId(userId string, client *firestore.Client, ctx context.Context) (string, error) {
	log.Println("RetrieveConnectionId()")
	var userBody UserBodyStruct
	doc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_ = doc.DataTo(&userBody)
	return userBody.ConnectionId, nil
}

func RetrieveCurrentSeatId(userId string, client *firestore.Client, ctx context.Context) (int, CustomError) {
	log.Println("RetrieveCurrentSeatId()")
	userInfo, err := RetrieveUserInfo(userId, client, ctx)
	if err != nil {
		return 0, CustomError{ErrorType: Unknown, Body: err}
	}
	roomId := userInfo.In
	if roomId == "" {
		errString := "the user is not in any room now"
		return 0, UserNotInTheRoom.New(errString)
	}
	roomStatus, err := RetrieveRoomInfo(roomId, client, ctx)
	if err != nil {
		return 0, CustomError{ErrorType: Unknown, Body: err}
	}
	for _, user := range roomStatus.Users {
		if user.UserId == userId {
			return user.SeatId, CustomError{Body: nil}
		}
	}
	errString := "the user is not in " + roomId
	log.Println(errString)
	return 0, UserNotInAnyRoom.New(errString)
}

func RetrieveRooms(client *firestore.Client, ctx context.Context) ([]RoomStruct, error) {
	log.Println("RetrieveRooms()")
	var rooms []RoomStruct

	// roomsのコレクションを取得
	iter := client.Collection(ROOMS).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
			return []RoomStruct{}, err
		}
		var _room RoomBodyStruct
		_ = doc.DataTo(&_room)
		room := RoomStruct{
			RoomId: doc.Ref.ID,
			Body:   _room,
		}
		if room.Body.Users == nil {
			room.Body.Users = []UserSeatSetStruct{}
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// user statusのonlineフラグから判断
func RetrieveOnlineUsersAsStatus(client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
	log.Println("RetrieveOnlineUsersAsStatus()")
	userDocs, err := client.Collection(USERS).Documents(ctx).GetAll()
	if err != nil {
		log.Println(err)
		return []UserStruct{}, err
	}

	app, _ := InitializeFirebaseApp(ctx)
	authClient, _ := app.Auth(ctx)

	if len(userDocs) == 0 {
		log.Println("there is no user.")
		return []UserStruct{}, nil
	} else {
		var userList []UserStruct
		for _, doc := range userDocs {
			var _user UserBodyStruct
			_ = doc.DataTo(&_user)
			if _user.Online {
				var displayName string
				user, err := authClient.GetUser(ctx, doc.Ref.ID)
				if err != nil {
					// これはfirebase authに登録されてないテストユーザーの場合、例外として起こりうる。
					log.Println("failed authClient.GetUser().")
					displayName = ""
				} else {
					displayName = user.DisplayName
				}
				userList = append(userList, UserStruct{
					UserId:      doc.Ref.ID,
					DisplayName: displayName,
					Body:        _user,
				})
			}
		}
		if userList == nil {
			userList = []UserStruct{}
		}
		return userList, nil
	}
}

// 実際のroomsのusersから判断
func RetrieveOnlineUsersInRooms(client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
	log.Println("RetrieveOnlineUsersInRooms()")
	roomDocs, err := client.Collection(ROOMS).Documents(ctx).GetAll()
	if err != nil {
		log.Println(err)
		return []UserStruct{}, err
	}

	app, _ := InitializeFirebaseApp(ctx)
	authClient, _ := app.Auth(ctx)

	if len(roomDocs) <= 0 {
		log.Println("there is no room.")
		return []UserStruct{}, nil
	} else {
		var userList []UserStruct
		for _, roomDoc := range roomDocs {
			var _room RoomBodyStruct
			_ = roomDoc.DataTo(&_room)
			for _, userSeat := range _room.Users {
				userId := userSeat.UserId
				userInfo, _ := RetrieveUserInfo(userId, client, ctx)
				var displayName string
				user, err := authClient.GetUser(ctx, userId)
				if err != nil {
					// これはfirebase authに登録されてないテストユーザーの場合、例外として起こりうる。
					log.Println("failed authClient.GetUser().")
					displayName = ""
				} else {
					displayName = user.DisplayName
				}
				userList = append(userList, UserStruct{
					UserId:      userId,
					DisplayName: displayName,
					Body:        userInfo,
				})
			}
		}
		if userList == nil {
			userList = []UserStruct{}
		}
		return userList, nil
	}
}

func RetrieveRoomUsers(roomId string, client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
	log.Println("RetrieveRoomUsers()")
	var err error
	var users []UserStruct

	authClient, _ := InitializeFirebaseAuthClient(ctx)

	roomInfo, err := RetrieveRoomInfo(roomId, client, ctx)
	if err != nil {
	} else {
		for _, user := range roomInfo.Users {
			userBody, err := RetrieveUserInfo(user.UserId, client, ctx)
			if err != nil {
			} else {
				userInfo, _ := authClient.GetUser(ctx, user.UserId)
				users = append(users, UserStruct{
					UserId:      user.UserId,
					DisplayName: userInfo.DisplayName,
					Body:        userBody,
				})
			}
		}
	}
	if users == nil {
		users = []UserStruct{}
	}
	return users, err
}

func RetrieveRoomInfo(roomId string, client *firestore.Client, ctx context.Context) (RoomBodyStruct, error) {
	log.Println("RetrieveRoomInfo()")
	var roomBodyStruct RoomBodyStruct

	room, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return RoomBodyStruct{}, err
	} else {
		_ = room.DataTo(&roomBodyStruct)
		if roomBodyStruct.Users == nil {
			roomBodyStruct.Users = []UserSeatSetStruct{} // jsonにした時、中身がない場合にnullではなく[]にする
		}
		return roomBodyStruct, nil
	}
}

func RetrieveNews(numNews int, client *firestore.Client, ctx context.Context) ([]NewsStruct, error) {
	log.Println("RetrieveNews()")
	var newsList []NewsStruct
	var err error

	// roomsのコレクションを取得
	iter := client.Collection(NEWS).OrderBy("updated", firestore.Desc).Limit(numNews).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {
			log.Printf("failed to iterate: %v\n", err)
			return []NewsStruct{}, err
		}
		var _news NewsBodyStruct
		_ = doc.DataTo(&_news)
		news := NewsStruct{
			NewsId:   doc.Ref.ID,
			NewsBody: _news,
		}
		newsList = append(newsList, news)
	}
	if newsList == nil {
		newsList = []NewsStruct{}
	}
	return newsList, err
}

func RetrieveUserInfo(userId string, client *firestore.Client, ctx context.Context) (UserBodyStruct, error) {
	log.Println("RetrieveUserInfo()")
	var userBodyStruct UserBodyStruct
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
	} else {
		_ = userDoc.DataTo(&userBodyStruct)
	}
	return userBodyStruct, err
}

func RecordHistory(details interface{}, client *firestore.Client, ctx context.Context) error {
	log.Println("RecordHistory()")
	_, _, err := client.Collection(HISTORY).Add(ctx,
		details,
	)
	if err != nil {
		log.Println("failed to make a record.")
	}
	return err
}

func RecordLastAccess(userId string, client *firestore.Client, ctx context.Context) error {
	log.Println("RecordLastAccess()")
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-access": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func RecordEnteredTime(userId string, client *firestore.Client, ctx context.Context) error {
	log.Println("RecordEnteredTime()")
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-entered": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func RecordExitedTime(userId string, client *firestore.Client, ctx context.Context) error {
	log.Println("RecordExitedTime()")
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-exited": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateStatusMessage(userId string, statusMessage string, client *firestore.Client, ctx context.Context) error {
	log.Println("UpdateStatusMessage()")
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-access": time.Now(),
		"status":      statusMessage,
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func InWhichRoom(userId string, client *firestore.Client, ctx context.Context) (string, error) {
	log.Println("InWhichRoom()")
	rooms, err := RetrieveRooms(client, ctx)
	if err != nil {
		log.Println(err)
	} else {
		for _, room := range rooms {
			users := room.Body.Users
			for _, user := range users {
				if user.UserId == userId {
					return room.RoomId, nil
				}
			}
		}
	}
	return "", err
}

func _CreateNewRoom(roomId string, roomName string, roomType string, themeColorHex string, client *firestore.Client, ctx context.Context) error {
	log.Println("_CreateNewRoom()")
	_, err := client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
		"name":            roomName,
		"type":            roomType,
		"users":           []UserSeatSetStruct{},
		"created":         time.Now(),
		"theme-color-hex": themeColorHex,
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

// 全オンラインユーザーの最終アクセス時間を調べ、タイムアウトを判断して処理
func _UpdateDatabase(client *firestore.Client, ctx context.Context) error {
	log.Println("_UpdateDatabase()")

	// onlineとなっているユーザー
	onlineUsersAsStatus, err := RetrieveOnlineUsersAsStatus(client, ctx)
	if err != nil {
		log.Println("RetrieveOnlineUsersAsStatus() failed.")
		return err
	}
	if len(onlineUsersAsStatus) > 0 {
		for _, u := range onlineUsersAsStatus {
			// User statusがonlineなのにルームにいないことがないか
			if isInRoom, _, _ := IsInRoom(u.Body.In, u.UserId, client, ctx); !isInRoom {
				log.Printf("%s is not in the room though online is true.\n", u.UserId)
				_, _ = client.Collection(USERS).Doc(u.UserId).Set(ctx, map[string]interface{}{
					"online": false,
					"connection-id": "",
				}, firestore.MergeAll)
			}
			
			// LastAccessから時間が経ってるのにルームに残ってないか
			lastAccess := u.Body.LastAccess
			timeElapsed := time.Now().Sub(lastAccess)
			if timeElapsed.Seconds() > TimeLimit {
				log.Printf("%s is put over time.\n", u.UserId)
				currentRoom := u.Body.In
				_ = LeaveRoom(currentRoom, u.UserId, client, ctx)
			}
			
			// connection idが設定されているか
			if u.Body.ConnectionId == "" {
				log.Printf("%s has no connection id though online is true.\n", u.UserId)
			}
			
			// inが設定されているか
			if u.Body.In == "" {
				log.Printf("%s's In is not given though online is true.\n", u.UserId)
			}
		}
	}
	
	// roomsに実際にいるユーザー
	onlineUsersInRooms, _ := RetrieveOnlineUsersInRooms(client, ctx)
	log.Println(onlineUsersInRooms)
	if len(onlineUsersAsStatus) != len(onlineUsersInRooms) {
		log.Println("len(onlineUsersAsStatus) != len(onlineUsersInRooms)")
	}
	for _, userInRoom := range onlineUsersInRooms {
		userInfo, _ := RetrieveUserInfo(userInRoom.UserId, client, ctx)
		if !userInfo.Online {
			log.Printf("%s is in room(%s) though online is false.\n", userInRoom.UserId, userInRoom.Body.In)
			roomId, _ := InWhichRoom(userInRoom.UserId, client, ctx)
			_ = LeaveRoom(roomId, userInRoom.UserId, client, ctx)
		}
		// todo room_layout的に有効な席に座っているかどうか

	}
	
	// todo
	// offlineとなっているユーザー
	// user statusがofflineなのにルームにいないか
	// offlineなのにconnection idが設定されていないか
	return nil
}

func Response(jsonBytes []byte) (events.APIGatewayProxyResponse, error) {
	log.Println("Response()")
	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200, // これないとInternal Server Errorになる
		Headers: map[string]string{
			"Content-Type": "application/json",
			//"Access-Control-Allow-Origin": "*",
			//"Access-Control-Allow-Methods": "GET,POST,HEAD,OPTIONS",
			//"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func CheckIfSeatAvailable(roomId string, seatId int, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("CheckIfSeatAvailable()")
	// そのIDの席が存在するか
	roomLayout, err := RetrieveRoomLayout(roomId, client, ctx)
	if err != nil {
		return false, err
	}
	seatDetected := false
	for _, seat := range roomLayout.Seats {
		if seat.Id == seatId {
			seatDetected = true
		}
	}
	if !seatDetected {
		return false, nil
	}

	// その席が空いているか
	var roomBody RoomBodyStruct
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, err
	}
	_ = roomDoc.DataTo(&roomBody)
	for _, userSeat := range roomBody.Users {
		if userSeat.SeatId == seatId {
			return false, nil
		}
	}
	return true, nil
}

func _EnterRoom(roomId string, userId string, seatId int, client *firestore.Client, ctx context.Context) CustomError {
	log.Println("_EnterRoom()")
	isSeatAvailable, err := CheckIfSeatAvailable(roomId, seatId, client, ctx)
	if err != nil {
		return CustomError{ErrorType: Unknown, Body: err}
	}
	if !isSeatAvailable {
		errString := "that seat is not available now (room id: " + roomId + ", seat id: " + strconv.Itoa(seatId) + ")"
		log.Println(errString)
		return SeatNotAvailable.New(errString)
	}
	userSeatSet := UserSeatSetStruct{
		UserId: userId,
		SeatId: seatId,
	}
	_, err = client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
		"users": firestore.ArrayUnion(userSeatSet),
	}, firestore.MergeAll)
	if err != nil {
		log.Println("failed _EnterRoom().")
		log.Println(err)
		return CustomError{Body: err}
	}
	//user statusを更新
	_, err = client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"online": true,
		"in":     roomId,
		"seat-id": seatId,
	}, firestore.MergeAll)
	if err != nil {
		log.Println("failed to update user info of " + userId + ".")
		log.Println(err)
		return CustomError{Body: err}
	}
	_ = RecordLastAccess(userId, client, ctx)
	_ = RecordEnteredTime(userId, client, ctx)
	_ = RecordHistory(EnterLeaveHistory{
		Activity: EnterActivity,
		RoomId:     roomId,
		UserId:  userId,
		Date:     time.Now(),
	}, client, ctx)
	roomBody, _ := RetrieveRoomInfo(roomId, client, ctx)
	authClient, _ := InitializeFirebaseAuthClient(ctx)
	user, _ := authClient.GetUser(ctx, userId)
	defer SendLiveChatMessage(user.DisplayName+"さんが"+roomBody.Name+"ルームに入りました。", client, ctx)

	return CustomError{ErrorType: Unknown, Body: nil}
}

// todo 定期実行
//func CheckHistoryConsistency(client *firestore.Client, ctx context.Context) error {
//	// 全ユーザーの入退室の整合性がとれているか（入室と退室が必ずペアで記録できているか）
//	iter := client.Collection(HISTORY).Where("user-id", "==", userId).OrderBy("date", firestore.Asc).Documents(ctx)
//	var historyData EnteringAndLeavingHistoryStruct
//	for {
//		doc, err := iter.Next()
//		if err == iterator.Done {
//			return nil
//		}
//		if err != nil {
//			return err
//		}
//		_ = doc.DataTo(&historyData)
//		if historyData.Activity == EnterActivity {
//			enteredDate := historyData.Date
//		} else if historyData.Activity == LeaveActivity {
//			leftDate := historyData.Date
//		}
//
//	}
//}

func UpdateTotalTime(userId string, roomId string, leftDate time.Time, client *firestore.Client, ctx context.Context) {
	log.Println("UpdateTotalTime()")
	var historyData EnterLeaveHistory

	docs, err := client.Collection(HISTORY).Where("user-id", "==", userId).Where("room-id", "==", roomId).Where("activity", "==", EnterActivity).OrderBy("date", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Println("could not fetch entering history: " + err.Error())
	}
	_ = docs[0].DataTo(&historyData)
	enteredDate := historyData.Date
	duration := leftDate.Sub(enteredDate)
	log.Printf("duration: %v", duration)

	roomBody, _ := RetrieveRoomInfo(roomId, client, ctx)
	roomType := roomBody.Type

	userBody, _ := RetrieveUserInfo(userId, client, ctx)
	totalStudyTime := time.Duration(userBody.TotalStudyTime) * time.Second
	totalBreakTime := time.Duration(userBody.TotalBreakTime) * time.Second

	if roomType == "study" {
		totalStudyTime = totalStudyTime + duration
		log.Printf("new totalStudyTime: %v", totalStudyTime)
		_, err = client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
			"total-study-time": int(totalStudyTime.Seconds()),
		}, firestore.MergeAll)
		if err != nil {
			log.Println("Failed to update total-study-time of " + userId)
		}
	} else if roomType == "break" {
		totalBreakTime = totalBreakTime + duration
		_, err = client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
			"total-break-time": int(totalBreakTime.Seconds()),
		}, firestore.MergeAll)
		if err != nil {
			log.Println("Failed to update total-break-time of " + userId)
		}
	}
}

func RetrieveRoomLayout(roomId string, client *firestore.Client, ctx context.Context) (RoomLayoutStruct, error) {
	log.Println("RetrieveRoomLayout()")
	var roomLayout RoomLayoutStruct
	doc, err := client.Collection(CONFIG).Doc(RoomLayoutsInfo).Collection(RoomLayouts).Doc(roomId).Get(ctx)
	if err != nil {
		log.Printf("failed to process client.Collection(CONFIG).Doc(ROOM_LAYOUTS_INFO).Collection(ROOM_LAYOUTS).Doc(roomId).Get(ctx), %v", err)
		return RoomLayoutStruct{}, err
	}
	_ = doc.DataTo(&roomLayout)
	roomLayout.RoomId = roomId
	return roomLayout, nil
}

func CurrentRoomLayoutVersion(roomId string, client *firestore.Client, ctx context.Context) (int, error) {
	log.Println("CurrentRoomLayoutVersion()")
	doc, err := client.Collection(CONFIG).Doc(RoomLayoutsInfo).Collection(RoomLayouts).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var roomLayout RoomLayoutStruct
	_ = doc.DataTo(&roomLayout)
	return roomLayout.Version, nil
}

func SaveRoomLayout(roomLayout RoomLayoutStruct, client *firestore.Client, ctx context.Context) error {
	log.Println("SaveRoomLayout()")

	// 履歴を保存
	oldRoomLayout, err := RetrieveRoomLayout(roomLayout.RoomId, client, ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	_ = RecordHistory(map[string]interface{}{
		"activity": NewRoomLayoutActivity,
		"old-room-layout": oldRoomLayout,
		"new-room-layout": roomLayout,
		"date": time.Now(),
	}, client, ctx)

	// 保存
	_, err = client.Collection(CONFIG).Doc(RoomLayoutsInfo).Collection(RoomLayouts).Doc(roomLayout.RoomId).Set(ctx, roomLayout)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Disconnect(connectionId string, client *firestore.Client, ctx context.Context) error {
	log.Println("Disconnect()")
	var apiGatewayConfig ApiGatewayConfigStruct
	configDoc, _ := client.Collection(CONFIG).Doc(ApiGateway).Get(ctx)
	_ = configDoc.DataTo(&apiGatewayConfig)
	websocketConnectionBaseUrl := apiGatewayConfig.WebsocketConnectionBaseUrl
	url := websocketConnectionBaseUrl + connectionId

	req, _ := http.NewRequest("DELETE", url, nil)
	httpClient := new(http.Client)
	resp, _ := httpClient.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	responseString := string(byteArray)
	log.Println(responseString)
	// todo まずPostmanで試す

	return nil
}

func FindUserWithConnectionId(connectionId string, client *firestore.Client, ctx context.Context) (UserStruct, CustomError) {
	log.Println("FindUserWithConnectionId()")
	onlineUsers, _ := RetrieveOnlineUsersAsStatus(client, ctx)
	for _, user := range onlineUsers {
		if user.Body.ConnectionId == connectionId {
			return user, CustomError{Body: nil}
		}
	}
	return UserStruct{}, NoSuchUserExists.New("no such user exists, searched by connection id:" + connectionId)
}

func (roomLayout RoomLayoutStruct) SetIsVacant(client *firestore.Client, ctx context.Context) RoomLayoutStruct {
	for i, seat := range roomLayout.Seats {
		roomLayout.Seats[i].IsVacant = true
		roomStatus, _ := RetrieveRoomInfo(roomLayout.RoomId, client, ctx)
		for _, usedSeat := range roomStatus.Users {
			if usedSeat.SeatId == seat.Id {
				roomLayout.Seats[i].IsVacant = false
				break
			}
		}
	}
	return roomLayout
}

