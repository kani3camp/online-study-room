package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"time"
)

// todo これもcommon.goでよくない？
// 全オンラインユーザーの最終アクセス時間を調べ、タイムアウトを判断して処理
func UpdateDatabase(client *firestore.Client, ctx context.Context)  {
	fmt.Println("updating database...")
	
	users, _ := RetrieveOnlineUsers(client, ctx)
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
