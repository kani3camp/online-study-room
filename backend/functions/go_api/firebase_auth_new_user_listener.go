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

func CreateNewUser(userId string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
		"registration-date": time.Now(),
		"last-access":       time.Now(),
		"online":            false,
		"status":            "",
	})
	if err != nil {
		log.Println("failed to create new user")
		log.Println(err)
	}
	return err
}

func FirebaseAuthNewUserListener(ctx context.Context, e AuthEvent)  error {
	_, client := InitializeEventFunc()
	defer client.Close()
	
	userId := e.UID
	_ = CreateNewUser(userId, client, ctx)
	return nil
}
