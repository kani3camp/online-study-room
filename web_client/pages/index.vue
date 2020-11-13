<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <v-row>
            <v-flex>
              <v-col>
                <h2 style="display: inline-block">
                  ルーム
                </h2>
                に入室して作業開始！
              </v-col>
            </v-flex>
            <v-spacer/>
            <v-flex>
              <v-col class="d-flex flex-row-reverse">
                <v-btn
                  outlined
                  @click="loadRooms"
                >
                  <v-icon>mdi-reload</v-icon>
                </v-btn>
              </v-col>
            </v-flex>
          </v-row>

        </v-flex>

        <v-container
          v-show="loading"
          class="fill-height"
          fluid
        >
          <v-row
            align="center"
            justify="center"
          >
            <v-col class="text-center">
              <div class="big-char">
                Loading...
              </div>
            </v-col>
          </v-row>
        </v-container>

        <v-container v-show="! loading">
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
              @click="confirmEntering(index)"
            >
              <v-hover v-slot="{ hover }">
                <v-card
                  class="ma-2 pa-3"
                  :elevation="hover ? 10 : 2"
                >
                  <v-layout justify-center>
                    <v-card-title>
                      {{ room['room_body'].name }}
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
          <h2 style="display: inline-block">
            YouTubeライブ
          </h2>
        </v-flex>

        <!--        todo-->
      </v-container>


      <v-dialog
        v-model="if_show_dialog"
        width="500"
      >
        <v-card
          class="mx-auto"
          outlined
          :loading="entering"
        >
          <v-card-title>{{ selected_room_name }}の部屋 に入室しますか？</v-card-title>

          <v-card-actions box-sizing>
            <v-spacer />
            <v-btn
              :disabled="entering"
              text
              color="primary"
              @click="enterRoom"
            >
              入室する
            </v-btn>
            <v-btn
              :disabled="entering"
              pr-0
              text
              @click="if_show_dialog=false"
            >
              キャンセル
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>


      <v-dialog
        v-model="if_show_dialog_2"
        width="500"
      >
        <v-card
          class="mx-auto"
          outlined
        >
          <v-card-title>{{ dialog_message }}</v-card-title>

          <v-card-actions>
            <v-spacer />
            <v-btn
              text
              pr-0
              @click="if_show_dialog_2=false"
            >
              閉じる
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-main>
    <Footer />
  </v-app>
</template>

<script>
import common from '~/plugins/common'
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'

export default {
  components: {
    NavigationDrawer,
    ToolBar,
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
      },
    },
  },
  async created() {
    common.onAuthStateChanged(this)

    await this.loadRooms()
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
    },
    async loadRooms() {
      this.loading = true
      const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms'
      const resp = await common.httpGet(url, {})
      if (resp.result === 'ok') {
        this.rooms = resp.rooms
      } else {
        console.log(resp.message)
      }
      this.loading = false
    }
  },
}
</script>

<style>
main {
  background-color: #99b5b1;
}

.big-char {
  font-size: 3rem;
  color: #7f828b;
}
</style>
