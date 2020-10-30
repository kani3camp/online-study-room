package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"log"
	"net/http"
	"time"
)

const ROOMS = "rooms"
const USERS = "users"
const HISTORY = "history"
const CONFIG = "config"
const API = "api"
const NEWS = "news"

const ProjectId = "online-study-space"

//var ProjectId = os.Getenv("GOOGLE_CLOUD_PROJECT")	// なんか代入されない

const TimeLimit = 1800 // 秒

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

type RoomStruct struct {
	RoomId string         `json:"room_id"`
	Body   RoomBodyStruct `json:"room_body"`
}

type RoomBodyStruct struct {
	Created time.Time `firestore:"created" json:"created"`
	Name    string    `firestore:"name" json:"name"`
	Users   []string  `firestore:"users" json:"users"`
}

type UserStruct struct {
	UserId      string         `json:"user_id"`
	DisplayName string         `json:"display_name"`
	Body        UserBodyStruct `json:"user_body"`
}

type UserBodyStruct struct {
	In          string    `firestore:"in" json:"in"`
	LastAccess  time.Time `firestore:"last-access" json:"last_access"`
	LastEntered time.Time `firestore:"last-entered" json:"last_entered"`
	LastExited  time.Time `firestore:"last-exited" json:"last_exited"`
	LastStudied time.Time `firestore:"last-studied" json:"last_studied"`
	//Name        string    `firestore:"name" json:"name"` todo firestoreの方でも消す
	Online           bool      `firestore:"online" json:"online"`
	Status           string    `firestore:"status" json:"status"`
	RegistrationDate time.Time `firestore:"registration-date" json:"registration_date"`
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

func InitializeHttpFunc(w *http.ResponseWriter) (context.Context, *firestore.Client) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	return InitializeEventFunc()
}

func InitializeEventFunc() (context.Context, *firestore.Client) {
	ctx := context.Background()
	client, _ := InitializeFirestoreClient(ctx)
	return ctx, client
}

func InitializeFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, ProjectId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}

func InitializeFirebaseApp(ctx context.Context) (*firebase.App, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Println("failed to initialize firebase.App.")
		log.Println(err)
		return nil, err
	}
	return app, err
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
				user, _ := authClient.GetUser(ctx, doc.Ref.ID)
				userList = append(userList, UserStruct{
					UserId:      doc.Ref.ID,
					DisplayName: user.DisplayName,
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

func _CreateNewRoom(roomId string, roomName string, roomType string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
		"name":    roomName,
		"type":    roomType,
		"users":   []string{},
		"created": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	return err
}