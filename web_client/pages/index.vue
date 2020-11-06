<template>
  <v-app>
    <NavigationDrawer></NavigationDrawer>

    <ToolBar></ToolBar>

    <v-main>
      <v-container>
        <v-flex>
          <h2 style="display: inline-block">ルーム</h2>
          <span>に入室して作業開始！</span>
        </v-flex>

        <v-container v-show="loading"
                     class="fill-height"
                     fluid
        >
          <v-row
            align="center"
            justify="center"
          >
            <v-col class="text-center">
              <div class="big-char">Loading...</div>
            </v-col>
          </v-row>
        </v-container>

        <v-container v-show="! loading">
          <v-row>
            <v-col
              cols="12" sm="4" md="3" lg="3" xl="3"
              v-for="(room, index) in rooms"
              :key="room.room_id"
              @click="confirmEntering(index)"
              dense
            >
              <v-hover v-slot:default="{ hover }">
                <v-card class="ma-2 pa-3" :elevation="hover ? 10 : 2">
                  <v-layout justify-center>
                    <v-card-title>
                      {{ room.room_body.name }}
                    </v-card-title>
                  </v-layout>
                </v-card>
              </v-hover>
            </v-col>
          </v-row>
        </v-container>

      </v-container>

      <v-container>
        <v-flex>
          <h2 style="display: inline-block">YouTubeライブ</h2>
        </v-flex>

<!--        todo-->

      </v-container>





      <v-dialog v-model="if_show_dialog" width=500>
        <v-card class="mx-auto" outlined :loading="entering">
          <v-card-title>{{ selected_room_name }}の部屋 に入室しますか？</v-card-title>

          <v-card-actions box-sizing>
            <v-spacer></v-spacer>
            <v-btn :disabled="entering" @click="enterRoom" text color="primary">入室する</v-btn>
            <v-btn :disabled="entering" @click="if_show_dialog=false" pr-0 text>キャンセル</v-btn>
          </v-card-actions>

        </v-card>
      </v-dialog>


      <v-dialog v-model="if_show_dialog_2" width=500>
        <v-card class="mx-auto" outlined>
          <v-card-title>{{ dialog_message }}</v-card-title>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn @click="if_show_dialog_2=false" text pr-0>閉じる</v-btn>
          </v-card-actions>

        </v-card>
      </v-dialog>

    </v-main>
    <Footer></Footer>
  </v-app>
</template>

<script>
  import common from "~/plugins/common"
  import firebase from '../plugins/firebase'
  import Logo from "@/components/Logo"
  import NavigationDrawer from "@/components/NavigationDrawer"
  import ToolBar from "@/components/ToolBar"

  export default {
    components: {
      Logo,
      NavigationDrawer,
      ToolBar
    },
    data: () => ({
      rooms: null,
      if_show_dialog: false,
      if_show_dialog_2: false,
      dialog_message: null,
      selected_index: null,
      selected_room_name: null,
      entering: false,
      loading: false,
    }),
    computed: {
      drawer: {
        get() {
          return this.$store.state.drawer
        },
        set(value) {
          this.$store.commit('setDrawer', value)
        }
      }
    },
    async created() {
      common.onAuthStateChanged(this)

      this.loading = true
      const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms'
      const resp = await common.httpGet(url, {})
      if (resp.result === 'ok') {
        this.rooms = resp.rooms
      } else {
        console.log(resp.message)
      }
      this.loading = false
    },
    methods: {
      confirmEntering(index) {
        this.selected_index = index
        this.selected_room_name = this.rooms[this.selected_index]['room_body'].name
        this.if_show_dialog = true
      },
      async enterRoom() {
        if (this.$store.state.isSignedIn) {
          const vm = this

          this.entering = true

          const selected_room_id = this.rooms[this.selected_index].room_id
          const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/enter_room'
          const params = {
            user_id: this.$store.state.user.user_id,
            room_id: selected_room_id,
            id_token: this.$store.state.user.id_token,
          }
          const res = await common.httpPost(url, params).catch((e) => {
            console.log(e)
            vm.if_show_dialog = false
            vm.dialog_message = '通信に失敗しました。もう一度試してください。'
            vm.if_show_dialog_2 = true
          })

          if (res.result === 'ok') {
            this.if_show_dialog = false
            this.$store.commit('setRoomId', selected_room_id)
            await this.$router.push('/rooms/' + selected_room_id)
          } else {
            console.log(res)
            this.if_show_dialog = false
            this.dialog_message = 'エラーが発生しました。'
            this.if_show_dialog_2 = true
          }

          this.entering = false
        } else {
          this.if_show_dialog = false
          this.dialog_message = 'サインインしてください。'
          this.if_show_dialog_2 = true
        }
      }
    }
  }
</script>

<style>


.big-char {
  font-size: 3rem;
  color: #7F828B;
}
</style>
