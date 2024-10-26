<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <!-- <v-container>
        <v-flex>
          <h2 style="display: inline-block">
            <v-icon>mdi-youtube</v-icon>
            YouTubeライブ
          </h2>
        </v-flex>
        <v-flex>
          <v-col>
            <v-row justify="center">
              <a
                target="_blank"
                :href="youtubeLink"
              ><h3>ライブ配信を見に行く</h3></a>
            </v-row>
          </v-col>
        </v-flex>
      </v-container> -->

      <v-container>
        <v-flex>
          <v-row>
            <v-flex>
              <v-col>
                <h2 style="display: inline-block">
                  <v-icon>mdi-door-open</v-icon>
                  ルーム
                </h2>
                に入室して作業開始！
              </v-col>
            </v-flex>
            <v-flex>
              <v-col class="d-flex flex-row-reverse">
                <v-btn :disabled="loading" outlined @click="loadRooms">
                  <v-icon>mdi-reload</v-icon>
                </v-btn>
              </v-col>
            </v-flex>
          </v-row>
        </v-flex>

        <v-container v-show="loading" class="fill-height" fluid>
          <v-row align="center" justify="center">
            <v-col class="text-center">
              <div class="big-char">ロード中...</div>
            </v-col>
          </v-row>
        </v-container>

        <v-container v-show="!loading">
          <v-row>
            <v-col
              v-for="(room, index) in rooms"
              :key="room.room_id"
              cols="12"
              sm="4"
              md="3"
              lg="3"
              xl="3"
              dense
            >
              <v-hover v-slot="{ hover }">
                <v-card class="ma-2 pa-3" :elevation="hover ? 10 : 2" @click="enterRoom(index)">
                  <v-card-title>
                    {{ room['room_body'].name }}
                  </v-card-title>
                  <v-card-subtitle> {{ room['room_body']['users'].length }}人 </v-card-subtitle>
                </v-card>
              </v-hover>
            </v-col>
          </v-row>
        </v-container>
      </v-container>

      <Dialog
        :if-show-dialog="if_show_dialog"
        :card-title="dialog_message"
        :accept-needed="false"
        cancel-option-string="OK"
        @cancel="if_show_dialog = false"
      />

      <Dialog
        :if-show-dialog="if_show_dialog_2"
        :card-title="dialog_message"
        :accept-needed="true"
        accept-option-string="サインイン"
        cancel-option-string="閉じる"
        @accept="$router.push('/sign_in')"
        @cancel="if_show_dialog_2 = false"
      />
    </v-main>
    <Footer />
  </v-app>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import common from '../plugins/common'
import NavigationDrawer from '#components'
import ToolBar from '#components'
import RoomLayout from '#components'
import Dialog from '#components'
import { useMainStore } from '#imports'

const router = useRouter()
const store = useMainStore()

type Room = {
  room_body: {
    created: Date
    name: string
    theme_color_hex: string
    type: string
    users: []
  }
  room_id: string
}

// Reactive state
const rooms = ref<Room[]>([])
const if_show_dialog = ref(false)
const if_show_dialog_2 = ref(false)
const dialog_message = ref('')
const loading = ref(false)
const youtubeLink = common.key.youtubeLink

const drawer = computed({
  get: () => store.drawer,
  set: (value: boolean) => store.setDrawer(value),
})

onMounted(async () => {
  common.onAuthStateChanged()

  await loadRooms()
})

const enterRoom = async (roomIndex: number) => {
  const selected_room_id = rooms.value[roomIndex].room_id
  const selected_room_name = rooms.value[roomIndex].room_body.name

  dialog_message.value = `${selected_room_name}の部屋 に入室しますか？`

  if (store.isSignedIn) {
    const vm = this
    store.setRoomId(selected_room_id)
    store.setRoomName(selected_room_name)
    await router.push(`/enter/${selected_room_id}`)
  } else {
    dialog_message.value = 'サインインしてください。'
    if_show_dialog_2.value = true
  }
}

const loadRooms = async () => {
  loading.value = true
  const url = common.apiLink.rooms
  const resp = await common.httpGet(url, {})
  if (resp.result === 'ok') {
    rooms.value = resp.rooms
  } else {
    console.log(resp.message)
  }
  loading.value = false
}
</script>

<style>
h2 {
  color: #36479f;
}

iframe {
  border-width: 0;
}

.big-char {
  font-size: 2rem;
  color: #7f828b;
}
</style>
