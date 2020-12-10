package main

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
	newUserStatus := UserBodyStruct{
		RegistrationDate: time.Now(),
		LastAccess: time.Now(),
		Online: false,
		TotalStudyTime: 0,
		TotalBreakTime: 0,
	}
	_, err := client.Collection(USERS).Doc(userId).Set(ctx, newUserStatus)
	if err != nil {
		log.Println("failed to create new user")
		log.Println(err)
	}
	return err
}

func FirebaseAuthNewUserListener(ctx context.Context, e AuthEvent) error {
	log.Println("FirebaseAuthNewUserListener()")
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	userId := e.UID
	_ = CreateNewUser(userId, client, ctx)
	return nil
}
