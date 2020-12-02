<template>
  <v-app>
    <v-app-bar
      app
      flat
    >
      <v-btn
        icon
        @click="if_show_dialog=true"
      >
        <v-icon>mdi-close</v-icon>
      </v-btn>

      <v-dialog
        v-model="if_show_dialog"
        width="500"
      >
        <v-card :loading="exiting">
          <v-card-title>部屋を出ますか？</v-card-title>
          <v-card-actions>
            <v-spacer />
            <v-btn
              :disabled="exiting"
              text
              color="primary"
              @click="exitRoom"
            >
              退室する
            </v-btn>
            <v-btn
              :disabled="exiting"
              text
              @click="if_show_dialog=false"
            >
              キャンセル
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>


      <v-layout justify-center>
        <v-toolbar-title>{{ room_name }} の部屋</v-toolbar-title>
      </v-layout>
    </v-app-bar>

    <v-main>
      <v-list subheader>
        <v-list-item>
          <v-layout justify-center>
            入室時刻：{{ entered_time }}
          </v-layout>
        </v-list-item>

        <v-divider />

        <v-subheader>同じ部屋の他のユーザー</v-subheader>

        <v-list-item>
          <v-container>
            <v-row>
              <v-col
                v-for="(user_info, index) in other_users_info"
                :key="index"
                cols="12"
                sm="3"
                md="2"
                lg="2"
                xl="2"
                dense
              >
                <v-card
                  class="ma-2 pa-4"
                  outlined
                >
                  <v-layout justify-center>
                    <v-icon
                      x-large
                      color="green"
                    >
                      mdi-account-circle-outline
                    </v-icon>
                  </v-layout>
                  <v-layout justify-center>
                    <v-card-title>
                      {{ user_info.display_name }}
                    </v-card-title>
                  </v-layout>
                  <v-layout justify-center>
                    {{ user_info.time_study }}
                  </v-layout>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </v-list-item>
      </v-list>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import firebase from 'firebase'
import { UserStore } from '~/utils/store-accessor'

export default Vue.extend({
  name: 'Room',
  data: () => ({
    room_name: null,
    entered_time: null,
    room_status: null,
    if_show_dialog: false,
    exiting: false,
    other_users_info: [],
    stay_awake_timeout: null,
    user_timeout: null,
    socket: null,
    is_socket_open: false,
  }),
  async created() {
    const vm = this
    this.$onAuthStateChanged(vm)

    if (vm.$store.state.isSignedIn) {
      // 入室時刻を取得
      vm.user_timeout = setTimeout(() => {
        vm.updateUserData()
      }, 5000)

      await vm.fetchRoomData()
      await vm.startStudying()
      await vm.stayStudying()
    } else {
      await vm.$router.push('/')
    }
  },
  destroyed() {
    clearTimeout(this.stay_awake_timeout)
    clearTimeout(this.user_timeout)
    this.socket.close()
  },
  methods: {
    startStudying() {
      const vm = this
      vm.socket = new WebSocket('wss://0ieer51ju9.execute-api.ap-northeast-1.amazonaws.com/production')
      vm.socket.onopen = async () => {
        vm.is_socket_open = true
        const params = {
          action: 'connect',
          user_id: firebase.auth().currentUser.uid,
          id_token: await firebase.auth().currentUser.getIdToken(false),
          room_id: vm.$store.state.room_id,
          device_type: '',
        }
        vm.socket.send(JSON.stringify(params))
      }
      vm.socket.onmessage = async (event) => {
        const resp = JSON.parse(event.data)
        if (resp.is_ok) {
          const info = []
          let amIin = false
          for (const user of resp.users) {
            if (user.user_id !== firebase.auth().currentUser.uid) {
              const studySeconds = new Date().getTime() - new Date(user.user_body.last_entered).getTime()
              info.push({
                display_name: user.display_name.substr(0, 3),
                time_study: Math.floor(studySeconds / (1000 * 60)).toString() + '分',
              })
            } else {
              amIin = true
            }
          }
          if (!amIin) {
            console.log('部屋に自分がいないので退室処理')
            await this.$router.push('/')
          }
          this.other_users_info = info
        } else {
          console.error(resp.message)
          await vm.$router.push('/')
        }
      }
    },
    async stayStudying() {
      if (UserStore.info.isSignedIn) {
        const vm = this
        if (vm.is_socket_open) {
          const params = {
            user_id: firebase.auth().currentUser.uid,
            id_token: await firebase.auth().currentUser.getIdToken(false),
            room_id: vm.$store.state.room_id,
            device_type: '',
          }
          vm.socket.send(JSON.stringify(params))
        }

        vm.stay_awake_timeout = setTimeout(() => {
          vm.stayStudying()
        }, 10000)
      } else {
        await this.$router.push('/')
      }
    },
    async fetchRoomData() {
      if (UserStore.info.isSignedIn) {
        // const vm = this
        const roomId = UserStore.info.roomId
        const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_status'
        const params = { room_id: roomId }
        const resp = await this.$httpGet(url, params)

        if (resp.result === 'ok') {
          this.room_name = resp.room_status.room_body.name
          this.room_status = resp.room_status
        }
      } else {
        await this.$router.push('/')
      }
    },
    async updateUserData() {
      await this.$getUserData(this)
      const dateTime = UserStore.info.lastEntered
      if (dateTime) {
        this.entered_time = dateTime.getHours() + '時' + dateTime.getMinutes() + '分'
      }
    },
    async exitRoom() {
      const vm = this

      vm.exiting = true

      const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/exit_room'
      const params = {
        user_id: firebase.auth().currentUser.uid,
        room_id: vm.$store.state.room_id,
        id_token: await firebase.auth().currentUser.getIdToken(false),
      }
      const resp = await this.$httpPost(url, params)

      if (resp.result === 'ok') {
        UserStore.setRoomId('')
        await vm.$router.push('/')
      } else {
        console.log('Failed to exit room.')
        console.log(resp)
      }
      vm.exiting = false
      vm.if_show_dialog = false
    },
  },
})
</script>

<style scoped></style>

