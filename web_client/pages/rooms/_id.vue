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
            <v-row justify="end">
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
            </v-row>
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

<script>
import common from '~/plugins/common'

export default {
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
  }),
  async created() {
    common.onAuthStateChanged(this)

    if (this.$store.state.isSignedIn) {
      // 入室時刻を取得
      this.user_timeout = setTimeout(() => {
        this.updateUserData()
      }, 5000)

      await this.fetchRoomData()

      await this.stayAwake()
    } else {
      await this.$router.push('/')
    }
  },
  destroyed() {
    clearTimeout(this.stay_awake_timeout)
    clearTimeout(this.user_timeout)
  },
  methods: {
    async stayAwake() {
      if (this.$store.state.isSignedIn) {
        // 存在する部屋のroom_idでなければならない
        const vm = this
        const room_id = vm.$store.state.room_id
        let url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/staying_awake'
        let params = {
          user_id: vm.$store.state.user.user_id,
          id_token: vm.$store.state.user.id_token,
        }
        const resp = await common.httpPost(url, params)

        if (resp.result === 'ok') {
          let info = []
          let amIin = false
          for (const user of resp['users']) {
            if (user.user_id !== vm.$store.state.user.user_id) {
              const study_seconds = new Date().getTime() - new Date(user['user_body'].last_entered).getTime()
              info.push({
                display_name: user.display_name.substr(0, 3),
                time_study: Math.floor(study_seconds / (1000 * 60)).toString() + '分',
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
          // todo
        }
        this.room_status = resp.room_status
        this.stay_awake_timeout = setTimeout(() => {
          this.stayAwake()
        }, 10000)
      } else {
        await this.$router.push('/')
      }
    },
    async fetchRoomData() {
      if (this.$store.state.isSignedIn) {
        const vm = this
        const room_id = vm.$store.state.room_id
        let url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_status'
        let params = { room_id }
        const resp = await common.httpGet(url, params)

        if (resp.result === 'ok') {
          this.room_name = resp.room_status['room_body'].name
          this.room_status = resp.room_status
        }
      } else {
        await this.$router.push('/')
      }
    },
    async updateUserData() {
      await common.getUserData(this)
      const date_time = this.$store.state.user.last_entered
      if (date_time) {
        this.entered_time = date_time.getHours() + '時' + date_time.getMinutes() + '分'
      }
    },
    async exitRoom() {
      this.exiting = true
      const vm = this

      const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/exit_room'
      const params = {
        user_id: vm.$store.state.user.user_id,
        room_id: vm.$store.state.room_id,
        id_token: vm.$store.state.user.id_token,
      }
      const resp = await common.httpPost(url, params)

      if (resp.result === 'ok') {
        this.$store.commit('setRoomId', null)
        await this.$router.push('/')
      } else {
        console.log('Failed to exit room.')
        console.log(resp)
      }
      this.exiting = false
      this.if_show_dialog = false
    },
  },
}
</script>

<style scoped></style>
