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

      <Dialog
        :if-show-dialog="if_show_dialog"
        :accept-needed="true"
        :loading="exiting"
        accept-option-string="出る"
        cancel-option-string="キャンセル"
        card-title="ルームを出ますか？"
        @accept="closeSocket"
        @cancel="if_show_dialog=false"
      />

      <Dialog
        :if-show-dialog="if_show_result"
        :accept-needed="false"
        cancel-option-string="閉じる"
        card-title="切断されました！ トップページに戻ります。"
        @cancel="exitRoom"
      />

      <Dialog
        :if-show-dialog="if_show_error_dialog"
        :accept-needed="false"
        cancel-option-string="閉じる"
        card-title="エラーがおきました(泣) トップページに戻ります。"
        @cancel="exitRoom"
      />

      <v-layout justify-center>
        <v-toolbar-title>{{ room_name }} の部屋</v-toolbar-title>
      </v-layout>
    </v-app-bar>

    <v-main v-show="is_entered">
      <v-list
        subheader
      >
        <v-list-item>
          <v-layout justify-center>
            入室時刻：{{ entered_time }}
          </v-layout>
        </v-list-item>

        <v-divider />

        <v-subheader>同じ部屋の他のユーザー</v-subheader>

      </v-list>

      <v-container
        style="max-width: 800px"
      >
        <RoomLayout
          :room-id="room_id"
          :layout="room_layout"
        />
      </v-container>

    </v-main>

    <v-main v-show="!is_entered">
      <v-container
        v-show="!is_entered"
        class="fill-height"
        fluid
      >
        <v-row
          align="center"
          justify="center"
        >
          <v-col class="text-center">
            <div class="big-char">
              ルームに入室中...
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-main>

  </v-app>
</template>

<script>
import common from '~/plugins/common'
import firebase from '@/plugins/firebase'
import Dialog from '~/components/Dialog'
import RoomLayout from '~/components/RoomLayout'

export default {
  name: 'Room',
  components: {
    Dialog,
    RoomLayout,
  },
  beforeRouteLeave(to, from, next) {
    window.onbeforeunload = null
    next()
  },
  data() {
    return {
      room_id: '',
      room_name: '',
      // room_layoutの初期値はnull
      room_layout: null,
      entered_time: new Date().getHours() + '時' + new Date().getMinutes() + '分',
      room_status: null,
      if_show_dialog: false,
      if_show_result: false,
      if_show_error_dialog: false,
      exiting: false,
      stay_awake_timeout: null,
      socket: null,
      is_socket_open: false,
      is_entered: false,
    }
  },
  async created() {
    const vm = this
    // これ意味ある？↓
    common.onAuthStateChanged(vm)

    this.room_id = this.$store.state.room_id
    this.room_name = this.$store.state.room_name

    if (vm.$store.state.isSignedIn) {
      await vm.startStudying()
    } else {
      await vm.$router.push('/')
    }
  },
  mounted() {
    window.onbeforeunload = (e) => this.showAlert(e)
  },
  beforeDestroy() {
    clearTimeout(this.stay_awake_timeout)
    this.socket.close()
    if (this.$store.state.room_id != null) {
      this.exitRoom()
    }
  },
  methods: {
    showAlert(e) {
      e.returnValue = '退室する場合は退室ボタンを押してください。'
    },
    async startStudying() {
      console.log('startStudying')
      // websocket
      const vm = this
      vm.socket = new WebSocket('wss://0ieer51ju9.execute-api.ap-northeast-1.amazonaws.com/production')
      vm.socket.onopen = async () => {
        console.log('socket opened.')
        const params = {
          action: 'connect',
          user_id: firebase.auth().currentUser.uid,
          id_token: await firebase.auth().currentUser.getIdToken(false),
          room_id: vm.$store.state.room_id,
          seat_id: vm.$store.state.seat_id,
        }
        vm.socket.send(JSON.stringify(params))
      }
      vm.socket.onmessage = async (event) => {
        console.log('message received.')
        const resp = JSON.parse(event.data)
        if (resp['is_ok']) {
          if (!vm.is_entered) {
            console.log('入室成功！！')
            vm.is_entered = true
            vm.is_socket_open = true
            await vm.stayStudying()
          }
          vm.room_layout = resp['room_layout']
          let amIin = false
          for (const user of resp['users']) {
            if (user.user_id !== firebase.auth().currentUser.uid) {
            } else {
              amIin = true
            }
          }
          if (!amIin) {
            console.log('部屋に自分がいないので退室')
            this.if_show_error_dialog = true
          }
        } else {
          console.error(resp.message)
          this.if_show_error_dialog = true
        }
      }
      vm.socket.onclose = async () => {
        console.log('socket closed.')
        vm.exiting = false
        vm.$store.commit('setRoomId', '')
        vm.if_show_result = true
        vm.is_socket_open = false
      }
      vm.socket.onerror = async () => {
        console.error('socket error.')
        this.if_show_error_dialog = true
      }
    },
    async stayStudying() {
      console.log('stayStudying')
      if (this.$store.state.isSignedIn) {
        const vm = this
        if (vm.is_socket_open) {
          const params = {
            action: 'stay',
            user_id: firebase.auth().currentUser.uid,
            id_token: await firebase.auth().currentUser.getIdToken(false),
            room_id: vm.$store.state.room_id,
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
    closeSocket() {
      if (this.is_socket_open) {
        this.exiting = true
        this.socket.close()
      } else {
        this.exitRoom()
      }
    },
    async exitRoom() {
      await this.$router.push('/')
    },
  },
}
</script>
<style scoped></style>

