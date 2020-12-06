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
          >
            <v-hover v-slot="{ hover }">
              <v-card
                class="ma-2 pa-3"
                :elevation="hover ? 10 : 2"
                @click="confirmEntering(index)"
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

      <Dialog
        :if-show-dialog="if_show_dialog"
        :loading="entering"
        :card-title="dialog_message"
        :accept-needed="true"
        accept-option-string="入室する"
        cancel-option-string="キャンセル"
        @accept="enterRoom"
        @cancel="if_show_dialog = false"
      />

      <Dialog
        :if-show-dialog="if_show_dialog_2"
        :card-title="dialog_message"
        :accept-needed="true"
        accept-option-string="サインイン"
        cancel-option-string="閉じる"
        @accept="$router.push('/sign_in')"
        @cancel="if_show_dialog_2=false"
      />
    </v-main>

    <Footer />
  </v-app>
</template>

<script>
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'
import common from '@/plugins/common'
import firebase from 'firebase/app'
import Dialog from '~/components/Dialog'

export default {
  name: 'AllRooms',
  components: {
    NavigationDrawer,
    ToolBar,
    Dialog,
  },
  data: () => ({
    rooms: null,
    if_show_dialog: false,
    if_show_dialog_2: false,
    dialog_message: '',
    selected_index: -1,
    selected_room_name: '',
    entering: false,
    loading: false,
  }),
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
      this.dialog_message = this.selected_room_name + 'の部屋 に入室しますか？'
      this.if_show_dialog = true
    },
    async enterRoom() {
      if (this.$store.state.isSignedIn) {
        const vm = this

        this.entering = true

        const selected_room_id = this.rooms[this.selected_index].room_id
        const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/enter_room'
        const params = {
          user_id: firebase.auth().currentUser.uid,
          room_id: selected_room_id,
          id_token: await firebase.auth().currentUser.getIdToken(false),
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
          await this.$router.push('/in/' + selected_room_id)
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
}
</script>

<style scoped></style>
