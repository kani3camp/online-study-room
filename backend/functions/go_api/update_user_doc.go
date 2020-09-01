package go_api

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)


type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	Fields     EventRoomBodyStruct `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

type EventRoomBodyStruct struct {
	Created TimestampValue `firestore:"created"`
	Name    StringValue    `firestore:"name"`
	Users   ArrayValue  `firestore:"users"`
}

type TimestampValue struct {
	TimestampValue time.Time `json:"timestampValue"`
}
type IntegerValue struct {
	IntegerValue string `json:"integerValue"`
}
type StringValue struct {
	StringValue string `json:"stringValue"`
}
type ArrayValue struct {
	ArrayValue Values `json:"arrayValue"`
}
type Values struct {
	Values  []StringValue  `json:"values"`
}



// ユーザーの入退室がトリガー
func UpdateUserDoc(ctx context.Context, e FirestoreEvent) error {
	_, client := InitializeEventFunc()
	defer client.Close()
	
	_, err := metadata.FromContext(ctx)
	if err != nil {return fmt.Errorf("metadata.FromContext: %v", err)}
	
	previousUsers := e.OldValue.Fields.Users.ArrayValue.Values
	newUsers := e.Value.Fields.Users.ArrayValue.Values
	
	log.Println(previousUsers)
	log.Println(newUsers)
	
	var enteredUser, leftUser []string
	
	for _, newUser := range newUsers {
		if len(previousUsers) > 0 {
			for j, previousUser := range previousUsers {
				if newUser == previousUser {
					break
				} else if j+1 == len(previousUsers) {
					enteredUser = append(enteredUser, newUser.StringValue)
				}
			}
		} else if len(previousUsers) == 0 {
			enteredUser = append(enteredUser, newUser.StringValue)
		}
	}
	for _, previousUser := range previousUsers {
		if len(newUsers) > 0 {
			for j, newUser := range newUsers {
				if previousUser == newUser {
					break
				} else if j+1 == len(newUsers) {
					leftUser = append(leftUser, previousUser.StringValue)
				}
			}
		} else if len(newUsers) == 0 {
			leftUser = append(leftUser, previousUser.StringValue)
		}
	}
	
	if len(enteredUser) > 1 {log.Fatalln("More than 1 people entered : ",  enteredUser)}
	if len(leftUser) > 1 {log.Fatalln("More than 1 people left : ",  leftUser)}
	
	usersCollectionRef := client.Collection(USERS)
	fullPath := strings.Split(e.Value.Name, "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	doc := strings.Join(pathParts[1:], "/")
	
	if len(enteredUser) > 0 {
		log.Printf("Entered! entered_user is %s while left_user.length is %d\n", enteredUser[0], len(leftUser))
		userId := enteredUser[0]
		_, err = usersCollectionRef.Doc(userId).Set(ctx, map[string]interface{}{
			"online": true,
			"in":     doc,
			"last-access": firestore.ServerTimestamp,
		}, firestore.MergeAll)
		if err != nil {
			log.Fatalln("Failed to update user info of " + userId)
		}
		RecordLastAccess(userId, client, ctx)
		RecordEnteredTime(userId, client, ctx)
		Record(map[string]interface{}{
			"activity": "entering",
			"room":     doc,
			"user-id":  userId,
			"time":     firestore.ServerTimestamp,
		}, client, ctx)
		userBody, err1 := GetUserInfo(userId, client, ctx)
		roomBody, err2 := GetRoomInfo(doc, client, ctx)
		if err1 != nil {
			log.Fatalln(err1)
		} else if err2 != nil {
			log.Fatalln(err2)
		} else {
			defer SendLiveChatMessage(userBody.Name + "さんが" + roomBody.Name + "に入室しました。", client, ctx)
		}
	} else if len(leftUser) > 0 {
		log.Printf("Left! left_user is %s while entered_user.length is %d\n", leftUser[0], len(enteredUser))
		userId := leftUser[0]
		_, err = usersCollectionRef.Doc(userId).Set(ctx, map[string]interface{}{
			"online":       false,
			"in":           "",
			"last-studied": firestore.ServerTimestamp,
			"last-access": firestore.ServerTimestamp,
		}, firestore.MergeAll)
		if err != nil {
			log.Fatalln("Failed to update user info of " + userId)
		}
		RecordLastAccess(userId, client, ctx)
		RecordExitedTime(userId, client, ctx)
		Record(map[string]interface{}{
			"activity": "leaving",
			"room":     doc,
			"user-id":  userId,
			"time":     firestore.ServerTimestamp,
		}, client, ctx)
		userBody, err1 := GetUserInfo(userId, client, ctx)
		roomBody, err2 := GetRoomInfo(doc, client, ctx)
		if err1 != nil {
			log.Fatalln(err1)
		} else if err2 != nil {
			log.Fatalln(err2)
		} else {
			defer SendLiveChatMessage(userBody.Name + "さんが" + roomBody.Name + "を退室しました。", client, ctx)
		}
	} else {
		log.Println("No changes?")
	}
	return nil
}
