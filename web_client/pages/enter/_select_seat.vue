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

    <Dialog
      :if-show-dialog="if_show_dialog"
      :accept-needed="false"
      :card-title="dialog_message"
      cancel-option-string="戻る"
      @cancel="$router.push('/')"
    />

    <v-main v-show="!loading">
      <v-container style="max-width: 700px">
        <h2>席を選ぼう。</h2>

        <div id="seat-selector">
          <v-form
            class="mx-auto"
          >
            <v-select
              v-model="selected_seat_id"
              :items="vacant_seats"
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

      </v-container>
    </v-main>

    <v-main v-show="loading">
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
              ロード中...
            </div>
          </v-col>
        </v-row>
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
  data: () => ({
    room_id: '',
    room_name: '',
    room_layout: null,
    loading: true,
    vacant_seats: [],
    if_show_dialog: false,
    dialog_message: '',
    selected_seat_id: null,
  }),
  async created() {
    // room_idは$storeからも読み込める
    this.room_id = this.$route.params.select_seat
    this.room_name = this.$store.state.room_name

    // fetch layout
    await this.fetchRoomInfo()
  },
  methods: {
    async fetchRoomInfo() {
      if (this.$store.state.isSignedIn) {
        const vm = this
        let url = common.apiLink.room_status
        let params = { room_id: vm.room_id }
        const resp = await common.httpGet(url, params)

        if (resp.result === 'ok') {
          this.room_layout = resp.room_layout
          this.vacant_seats = this.room_layout.seats.filter((item) => {
            console.log(item.is_vacant)
            return item.is_vacant
          })
          this.loading = false
        } else {
          console.log(resp.message)
          this.dialog_message = 'エラーなりました（泣）'
          this.if_show_dialog = true
        }
      } else {
        this.dialog_message = 'サインインしてください。'
        this.if_show_dialog = false
      }
    },
    async enterRoom() {
      if (this.$store.state.isSignedIn) {
        const vm = this
        const seatId = this.selected_seat_id
        this.$store.commit('setSeatId', seatId)

        await this.$router.push('/in/' + vm.$store.state.room_id)
      } else {
        this.dialog_message = 'サインインしてください。'
        this.if_show_dialog = false
      }
    },
    goToTopPage() {
      this.$router.push('/')
    },
  },
}
</script>

<style scoped></style>
