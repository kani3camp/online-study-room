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
        :accept-needed="accept_needed"
        :loading="exiting"
        :accept-option-string="accept_string"
        :cancel-option-string="cancel_string"
        :card-title="dialog_message"
        @accept="on_accept"
        @cancel="on_cancel"
      />

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
          <RoomLayout
            :room-id="room_id"
            :layout="room_layout"
          />
        </v-list-item>
      </v-list>
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
    console.log('beforeRouteLeave: to=', to.path, ',from=', from.path)
    if (this.$store.state.room_id !== '') {
      // todo
      // window.alert('退室する場合は退室ボタンを押してください。')
      next()
    } else {
      window.onbeforeunload = null
      console.log('remove beforeunload')
      next()
    }
  },
  data() {
    return {
      room_id: '',
      room_name: '',
      room_layout: null,
      entered_time: new Date().getHours() + '時' + new Date().getMinutes() + '分',
      room_status: null,
      if_show_dialog: false,
      dialog_message: 'ルームを出ますか？',
      accept_needed: true,
      accept_string: '出る',
      cancel_string: 'キャンセル',
      on_accept: null,
      on_cancel: null,
      exiting: false,
      other_users_info: [],
      stay_awake_timeout: null,
      user_timeout: null,
      socket: null,
      is_socket_open: false,
      is_entered: false,
    }
  },
  async created() {
    const vm = this
    // todo これ意味ある？↓
    common.onAuthStateChanged(vm)

    this.room_name = this.$store.state.room_name
    this.on_accept = vm.closeSocket()
    this.on_cancel = vm.closeDialog()

    if (vm.$store.state.isSignedIn) {
      // 入室時刻を取得
      vm.user_timeout = setTimeout(() => {
        vm.updateUserData()
      }, 5000)

      await vm.startStudying()
    } else {
      await vm.$router.push('/')
    }
  },
  mounted() {
    window.onbeforeunload = (e) => this.showAlert(e)
    console.log('add beforeunload')
    console.log(window.onbeforeunload)
  },
  beforeDestroy() {
    clearTimeout(this.stay_awake_timeout)
    clearTimeout(this.user_timeout)
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
        vm.is_socket_open = true
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
            await vm.stayStudying()
          }
          let info = []
          let amIin = false
          for (const user of resp['users']) {
            if (user.user_id !== firebase.auth().currentUser.uid) {
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
          console.error(resp.message)
          await vm.exitRoom()
        }
      }
      vm.socket.onclose = async () => {
        console.log('socket closed.')
        vm.$store.commit('setRoomId', '')
        vm.dialog_message = '切断されました！ トップページに戻ります。'
        vm.accept_needed = false
        vm.cancel_string = '閉じる'
        vm.on_cancel = vm.exitRoom()
      }
      vm.socket.onerror = async () => {
        console.error('socket error.')
        // todo
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
    async updateUserData() {
      // const date_time = this.$store.state.user.last_entered
      // if (date_time) {
      //   this.entered_time = date_time.getHours() + '時' + date_time.getMinutes() + '分'
      // }
    },
    closeSocket() {
      if (this.socket) {
        this.exiting = true
        this.socket.close()
      } else {
        // todo
      }
    },
    closeDialog() {
      this.if_show_dialog = false
    },
    async exitRoom() {
      await this.$router.push('/')
    },
  },
}
</script>
<style scoped></style>

