package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"log"
	"net/http"
	"time"
)

const ROOMS = "rooms"
const USERS = "users"
const HISTORY = "history"
const CONFIG = "config"
const API = "api"

const ProjectId = "online-study-room-f1f30"
//var ProjectId = os.Getenv("GOOGLE_CLOUD_PROJECT")	// なんか代入されない

const TimeLimit = 1800 // 秒

const user_id = "user_id"
const room_id = "room_id"
const id_token = "id_token"

const RoomDoesNotExist = "Room does not exist."
const InvalidParams = "Invalid parameters."
const NoRoom = "There is no room."
const InvalidUser = "Invalid user."
const OK = "ok"
const ERROR = "error"
const UserAuthFailed = "User authentication failed."
const UserDoesNotExist = "User does not exist."

type RoomStruct struct {
	RoomId string `json:"room_id"`
	Body   RoomBodyStruct `json:"room_body"`
}

type RoomBodyStruct struct {
	Created time.Time `firestore:"created" json:"created"`
	Name    string    `firestore:"name" json:"name"`
	Users   []string  `firestore:"users" json:"users"`
}

type UserStruct struct {
	UserId string `json:"user_id"`
	Body   UserBodyStruct `json:"user_body"`
}

type UserBodyStruct struct {
	In          string    `firestore:"in" json:"in"`
	LastAccess  time.Time `firestore:"last-access" json:"last_access"`
	LastEntered time.Time `firestore:"last-entered" json:"last_entered"`
	LastExited  time.Time `firestore:"last-exited" json:"last_exited"`
	LastStudied time.Time `firestore:"last-studied" json:"last_studied"`
	Name        string    `firestore:"name" json:"name"`
	Online      bool      `firestore:"online" json:"online"`
	Status      string    `firestore:"status" json:"status"`
	RegistrationDate time.Time `firestore:"registration-date" json:"registration_date"`
}


func InitializeHttpFunc(w *http.ResponseWriter) (context.Context, *firestore.Client) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	return InitializeEventFunc()
}

func InitializeEventFunc() (context.Context, *firestore.Client) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, ProjectId)
	if err != nil {
		log.Fatalln(err)
	}
	return ctx, client
}

func InitializeFirebaseApp(ctx context.Context) (*firebase.App, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Println("Failed to initialize firebase.App.")
		log.Fatalln(err)
	}
	return app, err
}

func IsUserVerified(userId string, idToken string, ctx context.Context) bool {
	app, _ := InitializeFirebaseApp(ctx)
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}
	
	return userId == token.UID
}


func IsExistRoom(roomId string, client *firestore.Client, ctx context.Context) bool {
	log.Println("IsExistRoom() is running. roomId : " + roomId + ".")
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {log.Println(err)}
	return roomDoc.Exists()
}

func IsInRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) bool {
	log.Println("IsInRoom() is running. roomId : " + roomId + ". userId : " + userId + ".")
	roomDoc, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {log.Println(err)}
	var room RoomBodyStruct
	err = roomDoc.DataTo(&room)
	if err != nil {log.Println(err)}
	users := room.Users
	for _, u := range users {
		if u == userId {
			return true
		}
	}
	return false
}

func IsInUsers(userId string, client *firestore.Client, ctx context.Context) bool {
	log.Println("IsInUsers() is running. userId : " + userId + ".")
	userDoc, _ := client.Collection(USERS).Doc(userId).Get(ctx)
	return userDoc.Exists()
}

func IsOnline(userId string, client *firestore.Client, ctx context.Context) bool {
	log.Println("IsOnline() is running. userId : " + userId + ".")
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
		return false // エラーの場合もfalseを返すので注意
	} else {
		return userDoc.Data()["online"].(bool)
	}
}

func LeaveRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) (string, error) {
	log.Println("LeaveRoom() is running. roomId : " + roomId + ". userId : " + userId)
	var statement string
	var err error
	if IsExistRoom(roomId, client, ctx) {
		if IsInRoom(roomId, userId, client, ctx) {
			// 退室処理
			_, err = client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
				"users": firestore.ArrayRemove(userId),
			}, firestore.MergeAll)
			if err != nil {
				statement = err.Error()
			} else {
				statement = "Successfully exited " + roomId
			}
		} else {
			statement = "You are not in the room."
		}
	} else {
		err = errors.New(RoomDoesNotExist)
		statement = RoomDoesNotExist
	}
	log.Println(statement)
	return statement, err
}

func _OnlineUsers(client *firestore.Client, ctx context.Context) []UserStruct {
	userDocs, err := client.Collection(USERS).Documents(ctx).GetAll()
	if err != nil {log.Println(err)}
	
	if len(userDocs) == 0 {
		log.Println("There is no user.")
		return []UserStruct{}
	} else {
		var userList []UserStruct
		for _, doc := range userDocs {
			var _user UserBodyStruct
			_ = doc.DataTo(&_user)
			if _user.Online {
				userList = append(userList, UserStruct{
					UserId: doc.Ref.ID,
					Body:   _user,
				})
			}
		}
		return userList
	}
}


func Record(details interface{}, client *firestore.Client, ctx context.Context) {
	_, _, err := client.Collection(HISTORY).Add(ctx,
		details,
	)
	if err != nil {
		log.Println("Failed to make a record.")
	}
}

func RecordLastAccess(userId string, client *firestore.Client, ctx context.Context)  {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-access": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {log.Println(err)}
}

func RecordEnteredTime(userId string, client *firestore.Client, ctx context.Context)  {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-entered": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {log.Println(err)}
}

func RecordExitedTime(userId string, client *firestore.Client, ctx context.Context)  {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"last-exited": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {log.Println(err)}
}

func InWhichRoom(userId string, client *firestore.Client, ctx context.Context) (string, error) {
	rooms, err := retrieveRooms(client, ctx)
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