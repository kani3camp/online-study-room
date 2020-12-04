package main

import (
	"cloud.google.com/go/firestore"
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
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"time"
)

const ROOMS = "rooms"
const USERS = "users"
const HISTORY = "history"
const CONFIG = "config"
const API = "api"
const NEWS = "news"

const ProjectId = "online-study-space"
const PathToServiceAccount = "path/to/serviceAccount.json" // todo

//var ProjectId = os.Getenv("GOOGLE_CLOUD_PROJECT")	// なんか代入されない

const TimeLimit = 180 // 秒

const user_id = "user_id"
const room_id = "room_id"
const id_token = "id_token"

const RoomDoesNotExist = "room does not exist."
const InvalidParams = "invalid parameters."
const InvalidUser = "invalid user."
const InvalidValue = "invalid value."
const OK = "ok"
const ERROR = "error"
const UserAuthFailed = "user authentication failed."
const UserDoesNotExist = "user does not exist."
const Failed = "failed"

const EnterActivity = "entering"
const LeaveActivity = "leaving"

type RoomStruct struct {
	RoomId string         `json:"room_id"`
	Body   RoomBodyStruct `json:"room_body"`
}

type RoomBodyStruct struct {
	Created       time.Time `firestore:"created" json:"created"`
	Name          string    `firestore:"name" json:"name"`
	Users         []string  `firestore:"users" json:"users"`
	Type          string    `firestore:"type" json:"type"`
	ThemeColorHex string    `firestore:"theme-color-hex" json:"theme_color_hex"`
}

type UserStruct struct {
	UserId      string         `json:"user_id"`
	DisplayName string         `json:"display_name"`
	Body        UserBodyStruct `json:"user_body"`
}

type UserBodyStruct struct {
	In               string    `firestore:"in" json:"in"`
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

type EnteringAndLeavingHistoryStruct struct {
	Activity string    `firestore:"activity"`
	Room     string    `firestore:"room"`
	Date     time.Time `firestore:"date"`
	UserId   string    `firestore:"user-id"`
}

func InitializeHttpFunc() (context.Context, *firestore.Client) {
	return InitializeEventFunc()
}

func InitializeEventFunc() (context.Context, *firestore.Client) {
	ctx := context.Background()
	client, _ := InitializeFirestoreClient(ctx)
	return ctx, client
}

func InitializeFirestoreClient(ctx context.Context) (*firestore.Client, error) {
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

func InitializeFirebaseApp(ctx context.Context) (*firebase.App, error) {
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
	app, _ := InitializeFirebaseApp(ctx)
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return authClient, err
}

func CloseFirestoreClient(client *firestore.Client) {
	err := client.Close()
	if err != nil {
		log.Println("failed to close firestore client.")
	}
}

func RetrieveFirebaseCredentialInBytes() ([]byte, error) {
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

func IsUserVerified(userId string, idToken string, client *firestore.Client, ctx context.Context) (bool, error) {
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
	log.Println("IsExistRoom() is running. roomId : " + roomId + ".")
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return roomDoc.Exists(), nil
}

func IsInRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsInRoom() is running. roomId : " + roomId + ". userId : " + userId + ".")
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, err
	}
	var room RoomBodyStruct
	err = roomDoc.DataTo(&room)
	if err != nil {
		log.Println(err)
		return false, err
	}
	users := room.Users
	for _, u := range users {
		if u == userId {
			return true, nil
		}
	}
	return false, nil
}

func IsInUsers(userId string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsInUsers() is running. userId : " + userId + ".")
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		return false, err
	}
	return userDoc.Exists(), nil
}

func IsOnline(userId string, client *firestore.Client, ctx context.Context) (bool, error) {
	log.Println("IsOnline() is running. userId : " + userId + ".")
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false, err // エラーの場合もfalseを返すので注意
	} else {
		return userDoc.Data()["online"].(bool), nil
	}
}

func LeaveRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) error {
	log.Println("LeaveRoom() is running. roomId : " + roomId + ". userId : " + userId + ".")
	var err error
	if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
		log.Println(RoomDoesNotExist)
	} else if isInRoom, _ := IsInRoom(roomId, userId, client, ctx); !isInRoom {
		log.Println("you are not in the room.")
	} else {
		// 退室処理
		_, err = client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
			"users": firestore.ArrayRemove(userId),
		}, firestore.MergeAll)
		if err != nil {
			log.Println("failed to remove user from room.")
		} else {
		}
	}
	return err
}

func RetrieveRooms(client *firestore.Client, ctx context.Context) ([]RoomStruct, error) {
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
			room.Body.Users = []string{}
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func RetrieveOnlineUsers(client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
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

func RetrieveRoomUsers(roomId string, client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
	var err error
	var users []UserStruct

	authClient, _ := InitializeFirebaseAuthClient(ctx)

	roomInfo, err := RetrieveRoomInfo(roomId, client, ctx)
	if err != nil {
	} else {
		for _, userId := range roomInfo.Users {
			userBody, err := RetrieveUserInfo(userId, client, ctx)
			if err != nil {
			} else {
				user, _ := authClient.GetUser(ctx, userId)
				users = append(users, UserStruct{
					UserId:      userId,
					DisplayName: user.DisplayName,
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
	var roomBodyStruct RoomBodyStruct

	room, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return RoomBodyStruct{}, err
	} else {
		_ = room.DataTo(&roomBodyStruct)
		if roomBodyStruct.Users == nil {
			roomBodyStruct.Users = []string{} // jsonにした時、中身がない場合にnullではなく[]にする
		}
		return roomBodyStruct, nil
	}
}

func RetrieveNews(numNews int, client *firestore.Client, ctx context.Context) ([]NewsStruct, error) {
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
	_, _, err := client.Collection(HISTORY).Add(ctx,
		details,
	)
	if err != nil {
		log.Println("failed to make a record.")
	}
	return err
}

func RecordLastAccess(userId string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-access": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func RecordEnteredTime(userId string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-entered": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func RecordExitedTime(userId string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-exited": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateStatusMessage(userId string, statusMessage string, client *firestore.Client, ctx context.Context) error {
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
	println("InWhichRoom() running.")
	rooms, err := RetrieveRooms(client, ctx)
	if err != nil {
		log.Println(err)
	} else {
		for _, room := range rooms {
			users := room.Body.Users
			for _, user := range users {
				if user == userId {
					return room.RoomId, nil
				}
			}
		}
	}
	return "", err
}

func _CreateNewRoom(roomId string, roomName string, roomType string, themeColorHex string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
		"name":            roomName,
		"type":            roomType,
		"users":           []string{},
		"created":         time.Now(),
		"theme-color-hex": themeColorHex,
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}

// 全オンラインユーザーの最終アクセス時間を調べ、タイムアウトを判断して処理
func _UpdateDatabase(client *firestore.Client, ctx context.Context) {
	log.Println("updating database...")

	users, err := RetrieveOnlineUsers(client, ctx)
	if err != nil {
		log.Println("RetrieveOnlineUsers() failed.")
	}
	if len(users) > 0 {
		for _, u := range users {
			lastAccess := u.Body.LastAccess
			timeElapsed := time.Now().Sub(lastAccess)
			if timeElapsed.Seconds() > TimeLimit {
				log.Printf("%s is put over time.\n", u.UserId)
				currentRoom := u.Body.In
				_ = LeaveRoom(currentRoom, u.UserId, client, ctx)
			}
		}
	}
}

func Response(jsonBytes []byte) (events.APIGatewayProxyResponse, error) {
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

func _EnterRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
		"users": firestore.ArrayUnion(userId),
	}, firestore.MergeAll)
	if err != nil {
		log.Println("failed _EnterRoom().")
		log.Println(err)
	}
	return err
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
	var historyData EnteringAndLeavingHistoryStruct

	docs, err := client.Collection(HISTORY).Where("user-id", "==", userId).Where("room", "==", roomId).Where("activity", "==", EnterActivity).OrderBy("date", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalln("could not fetch entering history: " + err.Error())
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
			log.Fatalln("Failed to update user info of " + userId)
		}
	} else if roomType == "break" {
		totalBreakTime = totalBreakTime + duration
		_, err = client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
			"total-break-time": int(totalBreakTime.Seconds()),
		}, firestore.MergeAll)
		if err != nil {
			log.Fatalln("Failed to update user info of " + userId)
		}
	}
}
