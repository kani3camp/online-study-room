package main

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
	CreateTime time.Time           `json:"createTime"`
	Fields     EventRoomBodyStruct `json:"fields"`
	Name       string              `json:"name"`
	UpdateTime time.Time           `json:"updateTime"`
}

type EventRoomBodyStruct struct {
	Created TimestampValue `firestore:"created"`
	Name    StringValue    `firestore:"name"`
	Users   ArrayValue     `firestore:"users"`
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
	Values []StringValue `json:"values"`
}

// todo 要修正
// ユーザーの入退室がトリガー
func UpdateUserDoc(ctx context.Context, e FirestoreEvent) error {
	now := time.Now()
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	_, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	previousUsers := e.OldValue.Fields.Users.ArrayValue.Values
	newUsers := e.Value.Fields.Users.ArrayValue.Values

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

	if len(enteredUser) > 1 {
		log.Fatalln("more than 1 people entered : ", enteredUser)
	}
	if len(leftUser) > 1 {
		log.Fatalln("more than 1 people left : ", leftUser)
	}
	if len(enteredUser) > 1 || len(leftUser) > 1 {
		log.Println("previousUsers : ", previousUsers)
		log.Println("newUsers : ", newUsers)
	}

	usersCollectionRef := client.Collection(USERS)
	fullPath := strings.Split(e.Value.Name, "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	doc := strings.Join(pathParts[1:], "/")

	if len(enteredUser) > 0 {
		log.Printf("Entered! entered_user is %s\n", enteredUser[0])
		userId := enteredUser[0]
		seatId := 0 // todo
		_, err = usersCollectionRef.Doc(userId).Set(ctx, map[string]interface{}{
			"online": true,
			"in":     doc,
			"seat-id": seatId,
		}, firestore.MergeAll)
		if err != nil {
			log.Println("failed to update user info of " + userId + ".")
		}
		_ = RecordLastAccess(userId, client, ctx)
		_ = RecordEnteredTime(userId, client, ctx)
		_ = RecordHistory(map[string]interface{}{
			"activity": EnterActivity,
			"room":     doc,
			"user-id":  userId,
			"date":     now,
		}, client, ctx)
		roomBody, _ := RetrieveRoomInfo(doc, client, ctx)
		authClient, _ := InitializeFirebaseAuthClient(ctx)
		user, _ := authClient.GetUser(ctx, userId)
		defer SendLiveChatMessage(user.DisplayName+"さんが"+roomBody.Name+"ルームに入りました。", client, ctx)
	} else if len(leftUser) > 0 {
		log.Printf("Left! left_user is %s\n", leftUser[0])
		userId := leftUser[0]
		_, err = usersCollectionRef.Doc(userId).Set(ctx, map[string]interface{}{
			"online":       false,
			"in":           "",
			"seat-id": 0,
			"last-studied": now,
		}, firestore.MergeAll)
		if err != nil {
			log.Fatalln("Failed to update user info of " + userId)
		}
		_ = RecordLastAccess(userId, client, ctx)
		_ = RecordExitedTime(userId, client, ctx)
		_ = RecordHistory(map[string]interface{}{
			"activity": LeaveActivity,
			"room":     doc,
			"user-id":  userId,
			"date":     now,
		}, client, ctx)
		defer UpdateTotalTime(userId, doc, now, client, ctx)

		roomBody, _ := RetrieveRoomInfo(doc, client, ctx)
		authClient, _ := InitializeFirebaseAuthClient(ctx)
		user, _ := authClient.GetUser(ctx, userId)
		defer SendLiveChatMessage(user.DisplayName+"さんが"+roomBody.Name+"ルームを出ました。", client, ctx)
	} else {
		log.Println("no changes?")
	}
	return nil
}
