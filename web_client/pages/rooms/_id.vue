<template>
  <v-app>
    <v-app-bar app flat>
      <v-btn @click="if_show_dialog=true" icon><v-icon>mdi-close</v-icon></v-btn>

            <v-dialog v-model="if_show_dialog" width=500>
              <v-card :loading="exiting">

                <v-card-title>部屋を出ますか？</v-card-title>
                <v-card-actions>
                  <v-row justify="end">
                    <v-btn :disabled="exiting" @click="exitRoom" text color="primary">退室する</v-btn>
                    <v-btn :disabled="exiting" @click="if_show_dialog=false" text>キャンセル</v-btn>
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

        <v-divider inset></v-divider>

        <v-subheader>同じ部屋の他のユーザー</v-subheader>

        <v-list-item>
          <v-container>
            <v-row>
              <v-col
                cols="12" sm="3" md="2" lg="2" xl="2"
                v-for="(user_info, index) in other_users_info"
                :key="index"
                dense
              >
                  <v-card class="ma-2 pa-4" outlined>
                    <v-layout justify-center>
                      <v-icon x-large color="green">mdi-account-circle-outline</v-icon>
                    </v-layout>
                    <v-layout justify-center>
                      <v-card-title>
                        {{ user_info.user_name }}
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
import common from "~/plugins/common"

export default {
  name: "room",
  data: () => ({
    room_name: null,
    entered_time: null,
    room_status: null,
    if_show_dialog: false,
    exiting: false,
    other_users_info: [],
    timeout: null,
    room_timeout: null,
    user_timeout: null
  }),
  async created() {
    common.onAuthStateChanged(this)

    if (this.$store.state.isSignedIn) {
      // 入室時刻を取得
      await this.updateUserData()

      await this.updateRoomInfo()

      // users読み込み
      await this.getOtherUsersData()

      // todo staying awake
    } else {
      await this.$router.push('/')
    }
  },
  computed: {
  },
  methods: {
    async updateRoomInfo() {
      if (this.$store.state.isSignedIn) {
        // 存在する部屋のroom_idでなければならない
        const vm = this;
        const room_id = vm.$store.state.room_id
        let url = new URL("https://us-central1-online-study-room-f1f30.cloudfunctions.net/RoomStatus")
        url.search = new URLSearchParams({room_id}).toString()
        const resp = await fetch(url.toString(), {method: "GET"}).then(response =>
          response.json()
        )
        if (resp.result === 'ok') {
          this.room_name = resp.room_status.room_body.name
          const users = resp.room_status.room_body.users
          if (!users.includes(vm.$store.state.user.user_id)) {
            await this.$router.push('/')
            return
          }
        }
        this.room_status = resp.room_status
        this.room_timeout = setTimeout(() => {
          this.updateRoomInfo()
        }, 10000)
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

      this.user_timeout = setTimeout(() => {
        this.updateUserData()
      }, 15000)
    },
    async getOtherUsersData() {
      const url = new URL('https://us-central1-online-study-room-f1f30.cloudfunctions.net/UserStatus')
      const vm = this
      let info = []
      for (const user of this.room_status.room_body.users) {
        if (user !== vm.$store.state.user.user_id) {
          url.search = new URLSearchParams({user_id: user}).toString()
          const resp = await fetch(url.toString(), {method: 'GET'}).then(r => r.json())
          const data = resp['user_status']

          const study_seconds = new Date().getTime() - new Date(data['user_body'].last_entered).getTime()
          info.push({
            user_name: data['user_body'].name.substr(0, 3),
            time_study: Math.floor(study_seconds / (1000 * 60)).toString() + '分'
          })
        }
      }
      this.other_users_info = info
      this.timeout = setTimeout(() => {
        this.getOtherUsersData()
      }, 10000)
    },
    async exitRoom() {
      this.exiting = true
      const vm = this

      const url = "https://us-central1-online-study-room-f1f30.cloudfunctions.net/ExitRoom"
      const params = new URLSearchParams({
        user_id: vm.$store.state.user.user_id,
        room_id: vm.$store.state.room_id,
        id_token: vm.$store.state.user.id_token,
      })
      const resp = await fetch(url, {
        method: 'POST',
        body: params
      }).then(response => response.json())

      if (resp.result === 'ok') {
        this.$store.commit('setRoomId', null)
        await this.$router.push('/')
      } else {
        console.log('Failed to exit room.')
        console.log(resp)
      }
      this.exiting = false
      this.if_show_dialog = false
    }
  }
}
</script>

<style scoped></style>
