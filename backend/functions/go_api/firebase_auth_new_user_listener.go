package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"time"
)

type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

func createNewUser(userId string, client *firestore.Client, ctx context.Context) (*firestore.WriteResult, error) {
	return client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"registration-date": firestore.ServerTimestamp, // todo この書き方でいいんだっけ？
		"last-access":       firestore.ServerTimestamp,
		"online":            false,
		"status":            "",
	})
}

func FirebaseAuthNewUserListener(ctx context.Context, e AuthEvent)  error {
	_, client := InitializeEventFunc()
	defer client.Close()
	
	userId := e.UID
	_, err := createNewUser(userId, client, ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
