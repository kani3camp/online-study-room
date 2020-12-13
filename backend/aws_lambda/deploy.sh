# Windows

# 新しく関数をデプロイしたら、タイムアウトが3秒のため長めにしておく
# change_user_info, create_new_room, end_studying, news, online_users, room_layout, room_status, rooms,
# send_contact_form, start_studying, stay_studying, update_database,
# upload_room_layout, user_status

set GOOS=linux
go build -o main common.go send_live_chat_message.go user_status.go
C:\Users\momom\go\bin\build-lambda-zip.exe -output main.zip main
aws lambda create-function --function-name end_studying --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::652333062396:role/service-role/my-first-golang-lambda-function-role-cb8uw4th
aws lambda update-function-code --function-name user_status --zip-file fileb://main.zip


# Mac OS

GOOS=linux go build -o main common.go news.go
zip main.zip main

aws lambda create-function --function-name change_user_info --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::652333062396:role/service-role/my-first-golang-lambda-function-role-cb8uw4th

aws lambda update-function-code --function-name news --zip-file fileb://main.zip
