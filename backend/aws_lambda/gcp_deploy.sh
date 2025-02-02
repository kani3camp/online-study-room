# gcloudプロジェクト切り替えコマンド
# gcloud config set project <project-id>


function_name="UpdateUserDoc"
#project_id="online-study-space"
project_id="test-online-study-space"

#gcloud functions deploy $function_name \
#  --trigger-http \
#  --runtime go113 \
#  --allow-unauthenticated

gcloud functions deploy $function_name \
  --runtime go113 \
  --trigger-event providers/cloud.firestore/eventTypes/document.update \
  --trigger-resource "projects/$project_id/databases/(default)/documents/rooms/{roomId}"

#gcloud functions deploy $function_name \
#  --runtime go113 \
#  --trigger-event providers/firebase.auth/eventTypes/user.create


# Windows向け（コマンドプロンプトにコピペ）
gcloud functions deploy UserStatus --trigger-http --runtime go113 --allow-unauthenticated
gcloud functions deploy UpdateUserDoc --runtime go113 --trigger-event providers/cloud.firestore/eventTypes/document.update --trigger-resource "projects/online-study-space/databases/(default)/documents/rooms/{roomId}"
gcloud functions deploy FirebaseAuthNewUserListener --runtime go113 --trigger-event providers/firebase.auth/eventTypes/user.create

