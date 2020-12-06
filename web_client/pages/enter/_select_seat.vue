<template>
  <v-app>
    <v-app-bar
      app
      flat
    >
      <v-btn
        icon
        @click="goToTopPage"
      >
        <v-icon>mdi-close</v-icon>
      </v-btn>

      <v-layout justify-center>
        <v-toolbar-title>{{ propRoomName }} の部屋</v-toolbar-title>
      </v-layout>

    </v-app-bar>

    <v-main>
      <v-container>
        <h2>席を選ぼう。</h2>

        <RoomLayout
          :room-id="room_id"
          @selected="enterRoom"
        />

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
      </v-container>
    </v-main>
  </v-app>
</template>
<script>
import RoomLayout from '~/components/RoomLayout'
import firebase from '~/plugins/firebase'
import common from '~/plugins/common'
import Dialog from '~/components/Dialog'

export default {
  name: 'EnterRoom',
  components: {
    RoomLayout,
    Dialog,
  },
  props: {
    propRoomName: {
      type: String,
      required: true,
    },
  },
  data: () => ({
    room_id: '',
    room_layout: null,
    if_show_dialog: false,
    entering: false,
    dialog_message: '',
  }),
  created() {
    this.room_id = this.$route.params.select_seat
    console.log('渡されたroomIdは', this.room_id)
  },
  methods: {
    async enterRoom(seatId) {
      if (this.$store.state.isSignedIn) {
        const vm = this

        this.entering = true

        const selected_room_id = this.rooms[this.selected_index].room_id
        const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/enter_room'
        const params = {
          user_id: firebase.auth().currentUser.uid,
          id_token: await firebase.auth().currentUser.getIdToken(false),
          room_id: selected_room_id,
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
    goToTopPage() {
      this.$router.push('/')
    },
  },
}
</script>

<style scoped></style>
