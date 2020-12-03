<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>ルーム一覧</h2>
        </v-flex>
      </v-container>

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
            @click="confirmEntering(index)"
          >
            <v-hover v-slot="{ hover }">
              <v-card
                class="ma-2 pa-3"
                :elevation="hover ? 10 : 2"
              >
                <v-card-title>
                  {{ room['room_body'].name }}
                </v-card-title>
                <v-card-subtitle>
                  {{ room['room_body']['users'].length }}人
                </v-card-subtitle>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
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
          <v-card-title>
            {{ selected_room_name }}の部屋 に入室しますか？
          </v-card-title>

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
              @click="if_show_dialog = false"
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
              @click="if_show_dialog_2 = false"
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

<script lang="ts">
import NavigationDrawer from '@/components/NavigationDrawer.vue'
import ToolBar from '@/components/ToolBar.vue'
import firebase from 'firebase'
import Vue from 'vue'
import { UserStore } from '@/store'

export default Vue.extend({
  name: 'AllRooms',
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
  async created() {
    this.$onAuthStateChanged(this)

    this.loading = true
    const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms'
    const resp = await this.$httpGet(url, {})
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
      this.selected_room_name = this.rooms[this.selected_index].room_body.name
      this.if_show_dialog = true
    },
    async enterRoom() {
      if (UserStore.info.isSignedIn) {
        const vm = this

        this.entering = true

        const selectedRoomId = this.rooms[this.selected_index].room_id
        const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/enter_room'
        const params = {
          user_id: firebase.auth().currentUser.uid,
          room_id: selectedRoomId,
          id_token: await firebase.auth().currentUser.getIdToken(false),
        }
        const res = await this.$httpPost(url, params).catch((e) => {
          console.log(e)
          vm.if_show_dialog = false
          vm.dialog_message = '通信に失敗しました。もう一度試してください。'
          vm.if_show_dialog_2 = true
        })

        if (res.result === 'ok') {
          this.if_show_dialog = false
          UserStore.setRoomId(selectedRoomId)
          await this.$router.push('/rooms/' + selectedRoomId)
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
  },
})
</script>


<style scoped></style>
