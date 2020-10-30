# Windows

set GOOS=linux

aws lambda create-function --function-name test_function --runtime go1.x --zip-file fileb://test_function/main.zip --handler main --role arn:aws:iam::652333062396:role/service-role/my-first-golang-lambda-function-role-cb8uw4th

aws lambda update-function-code --function-name test_function --zip-file fileb://test_function/main.zip


# Mac OS

GOOS=linux go build -o main common.go change_user_info.go
zip main.zip main

aws lambda create-function --function-name change_user_info --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::652333062396:role/service-role/my-first-golang-lambda-function-role-cb8uw4th

aws lambda update-function-code --function-name change_user_info --zip-file fileb://main.zip
