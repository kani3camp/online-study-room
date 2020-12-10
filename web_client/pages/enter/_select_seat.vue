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
        <v-toolbar-title>{{ room_name }} の部屋</v-toolbar-title>
      </v-layout>

    </v-app-bar>

    <v-main>
      <v-container style="max-width: 700px">
        <h2>席を選ぼう。</h2>

        <div id="seat-selector">
          <v-form
            class="mx-auto"
          >
            <v-select
              v-model="selected_seat_id"
              :items="seats"
              item-value="id"
              item-text="id"
              label="座席番号を選んでください"
              outlined
            />
            <div>
              <v-btn
                color="primary"
                block
                elevation="3"
                :disabled="! selected_seat_id"
                @click="enterRoom()"
              >
                決定
              </v-btn>
            </div>
          </v-form>
        </div>

        <RoomLayout
          :room-id="room_id"
          :layout="room_layout"
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
// import roomLayoutJson from 'assets/mathematics-rom-layout.json'

export default {
  name: 'EnterRoom',
  components: {
    RoomLayout,
    Dialog,
  },
  data: () => ({
    room_id: '',
    room_name: '',
    room_layout: null,
    seats: [],
    if_show_dialog: false,
    entering: false,
    dialog_message: '',
    selected_seat_id: null,
  }),
  watch: {
    room_layout: function (newValue, oldValue) {
      if (newValue) {
        this.seats = newValue.seats
      }
    },
  },
  async created() {
    // room_idは$storeからも読み込める
    this.room_id = this.$route.params.select_seat
    this.room_name = this.$store.state.room_name

    // fetch layout
    await this.fetchRoomLayout()
  },
  methods: {
    async fetchRoomLayout() {
      if (this.$store.state.isSignedIn) {
        const vm = this
        let url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_layout'
        let params = { room_id: vm.room_id }
        const resp = await common.httpGet(url, params)

        if (resp.result === 'ok') {
          this.room_layout = resp.room_layout_data
        } else {
          console.log(resp.message)
          // todo top page へ戻る
        }
      } else {
        // todo dialog
        await this.$router.push('/')
      }
    },
    async enterRoom() {
      if (this.$store.state.isSignedIn) {
        const vm = this
        const seatId = this.selected_seat_id
        this.$store.commit('setSeatId', seatId)

        await this.$router.push('/in/' + vm.$store.state.room_id)
      } else {
        // todo check
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
