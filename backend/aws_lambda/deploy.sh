# Windows

# 新しく関数をデプロイしたら、タイムアウトが3秒のため長めにしておく
# change_user_info, create_new_news, create_new_room, end_studying, news, online_users, room_layout, room_status, rooms,
# send_contact_form, start_studying, stay_studying, update_database,
# upload_room_layout, user_status
#
# test_change_user_info, test_create_new_news, test_create_new_room, test_end_studying, test_news, test_online_users,
# test_room_layout, test_room_status, test_rooms,
# test_send_contact_form, test_start_studying, test_stay_studying, test_update_database,
# test_upload_room_layout, test_user_status

set GOOS=linux
go build -o main common.go send_live_chat_message.go     create_new_news.go
C:\Users\momom\go\bin\build-lambda-zip.exe -output main.zip main
aws lambda create-function --function-name     create_new_news     --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::652333062396:role/service-role/my-first-golang-lambda-function-role-cb8uw4th
aws lambda update-function-code --function-name     create_new_news     --zip-file fileb://main.zip


# Mac OS

GOOS=linux go build -o main common.go news.go
zip main.zip main

aws lambda create-function --function-name change_user_info --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::652333062396:role/service-role/my-first-golang-lambda-function-role-cb8uw4th

aws lambda update-function-code --function-name news --zip-file fileb://main.zip
