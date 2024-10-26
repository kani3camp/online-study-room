<template>
  <v-app>
    <v-app-bar app flat>
      <v-btn icon @click="if_show_dialog = true">
        <v-icon>mdi-close</v-icon>
      </v-btn>

      <Dialog
        :if-show-dialog="if_show_dialog"
        :accept-needed="true"
        :loading="exiting"
        accept-option-string="出る"
        cancel-option-string="キャンセル"
        card-title="ルームを出ますか？"
        @accept="closeSocket"
        @cancel="if_show_dialog = false"
      />

      <Dialog
        :if-show-dialog="if_show_result"
        :accept-needed="false"
        cancel-option-string="閉じる"
        card-title="切断されました！ トップページに戻ります。"
        @cancel="exitRoom"
      />

      <Dialog
        :if-show-dialog="if_show_error_dialog"
        :accept-needed="false"
        cancel-option-string="閉じる"
        card-title="エラーがおきました(泣) トップページに戻ります。"
        @cancel="exitRoom"
      />

      <v-layout justify-center>
        <v-toolbar-title>{{ room_name }} の部屋</v-toolbar-title>
      </v-layout>
    </v-app-bar>

    <v-main v-show="is_entered">
      <v-list subheader>
        <v-list-item>
          <v-layout justify-center> 入室時刻：{{ entered_time }} </v-layout>
        </v-list-item>

        <v-divider />

        <v-list-item> 入室時間の上限は2時間となってます。今後伸びる予定です。 </v-list-item>
        <v-list-item>
          注意：スマホの方は、できるだけこの画面を閉じないようにしてください。通信が切断されて退室してしまいます。
        </v-list-item>
        <v-list-item>
          この問題については、今後のアップデートで解決する予定です。アイデアも募集中です。
        </v-list-item>

        <v-divider />

        <v-subheader>同じ部屋の他のユーザー</v-subheader>
      </v-list>

      <v-container style="max-width: 800px">
        <RoomLayout :room-id="room_id" :layout="room_layout" />
      </v-container>
    </v-main>

    <v-main v-show="!is_entered">
      <v-container v-show="!is_entered" class="fill-height" fluid>
        <v-row align="center" justify="center">
          <v-col class="text-center">
            <div class="big-char">ルームに入室中...</div>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import common from '~/plugins/common'
import firebase from '~/plugins/firebase'
import Dialog from '#components'
import RoomLayout from '#components'
import { useMainStore } from '~/stores'

const store = useMainStore()
const router = useRouter()

onBeforeRouteLeave((to, from, next) => {
  window.onbeforeunload = null
  next()
})

const room_id = ref('')
const room_name = ref('')
const room_layout = ref(null) // room_layoutの初期値はnull
const entered_time = ref(new Date().getHours() + '時' + new Date().getMinutes() + '分')
const room_status = ref(null)
const if_show_dialog = ref(false)
const if_show_result = ref(false)
const if_show_error_dialog = ref(false)
const exiting = ref(false)
const stay_awake_timeout = ref<NodeJS.Timeout | null>(null)
const socket = ref<WebSocket | null>(null)
const is_socket_open = ref(false)
const is_entered = ref(false)

onMounted(async () => {
  // これ意味ある？↓
  common.onAuthStateChanged()

  room_id.value = store.room_id
  room_name.value = store.room_name

  if (store.isSignedIn) {
    await startStudying()
  } else {
    await router.push('/')
  }
})

onMounted(() => {
  window.onbeforeunload = (e) => showAlert(e)
})

onBeforeUnmount(() => {
  clearTimeout(stay_awake_timeout)
  socket.value?.close()
  if (store.room_id != null) {
    exitRoom()
  }
})

const playSound = () => {
  const audioElement = document.getElementById('audioElement') as HTMLAudioElement
  audioElement.play()
}

const showAlert = (e: BeforeUnloadEvent) => {
  // e.returnValue = '退室する場合は退室ボタンを押してください。'
  e.preventDefault()
  // 一部のブラウザでは、ユーザーにカスタムメッセージを表示しないため、空文字を設定します
  e.returnValue = ''
}

const startStudying = async () => {
  console.log('startStudying')
  // websocket
  socket.value = new WebSocket(common.apiLink.websocket)
  socket.value.onopen = async () => {
    console.log('socket opened.')
    const params = {
      action: 'connect',
      user_id: firebase.auth().currentUser?.uid,
      id_token: await firebase.auth().currentUser?.getIdToken(false),
      room_id: store.room_id,
      seat_id: store.seat_id,
    }
    socket.value?.send(JSON.stringify(params))
  }
  socket.value.onmessage = async (event) => {
    const resp = JSON.parse(event.data)
    if (resp['is_ok']) {
      if (!is_entered) {
        console.log('入室成功！！')
        is_entered.value = true
        is_socket_open.value = true
        await stayStudying()
      }
      room_layout.value = resp['room_layout']
      let amIin = false
      for (const user of resp['users']) {
        if (user.user_id !== firebase.auth().currentUser?.uid) {
        } else {
          amIin = true
        }
      }
      if (!amIin) {
        console.log('部屋に自分がいないので退室')
        if_show_error_dialog.value = true
      }
    } else {
      console.error(resp.message)
      if_show_error_dialog.value = true
    }
  }
  socket.value.onclose = async () => {
    console.log('socket closed.')
    exiting.value = false
    store.setRoomId('')
    if_show_result.value = true
    is_socket_open.value = false
  }
  socket.value.onerror = async () => {
    console.error('socket error.')
    if_show_error_dialog.value = true
  }
}

const stayStudying = async () => {
  if (store.isSignedIn) {
    if (is_socket_open) {
      const params = {
        action: 'stay',
        user_id: firebase.auth().currentUser?.uid,
        id_token: await firebase.auth().currentUser?.getIdToken(false),
        room_id: store.room_id,
      }
      socket.value?.send(JSON.stringify(params))
    }

    stay_awake_timeout.value = setTimeout(() => {
      stayStudying()
    }, 10000)
  } else {
    await router.push('/')
  }
}

const closeSocket = () => {
  if (is_socket_open) {
    exiting.value = true
    socket.value?.close()
  } else {
    exitRoom()
  }
}

const exitRoom = async () => {
  await router.push('/')
}
</script>
<style scoped></style>
