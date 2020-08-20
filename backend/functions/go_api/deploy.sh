#!/bin/zsh
function_name="CreateNewRoom"
project_id="online-study-room-f1f30"

gcloud functions deploy $function_name \
  --trigger-http \
  --runtime go113 \
  --allow-unauthenticated

#gcloud functions deploy $function_name \
#  --runtime go113 \
#  --trigger-event providers/cloud.firestore/eventTypes/document.update \
#  --trigger-resource "projects/$project_id/databases/(default)/documents/rooms/{roomId}"

#gcloud functions deploy $function_name \
#  --runtime go113 \
#  --trigger-event providers/firebase.auth/eventTypes/user.create
