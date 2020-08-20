package go_api

import (
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
func FirebaseAuthNewUserListener(ctx context.Context, e AuthEvent)  error {
	_, client := InitializeEventFunc()
	defer client.Close()
	
	userId := e.UID
	_, err := createNewUser(userId, "", client, ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
