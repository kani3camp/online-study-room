<template>
  <h1>index.vue</h1>
  <div>
    <!-- Markup shared across all pages, ex: NavBar -->
    <NuxtPage />  
  </div>
</template>


<script>
import common from '~/plugins/common'
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'
import firebase from '@/plugins/firebase'
import RoomLayout from '~/components/RoomLayout'
import Dialog from '~/components/Dialog'

export default {
  components: {
    Dialog,
    NavigationDrawer,
    ToolBar,
  },
  data: () => ({
    rooms: null,
    if_show_dialog: false,
    if_show_dialog_2: false,
    dialog_message: '',
    loading: false,
    youtubeLink: common.key.youtubeLink,
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
    async enterRoom(roomIndex) {
      const selected_room_id = this.rooms[roomIndex].room_id
      const selected_room_name = this.rooms[roomIndex].room_body.name

      this.dialog_message = selected_room_name + 'の部屋 に入室しますか？'

      if (this.$store.state.isSignedIn) {
        const vm = this
        this.$store.commit('setRoomId', selected_room_id)
        this.$store.commit('setRoomName', selected_room_name)
        await this.$router.push('/enter/' + vm.rooms[roomIndex].room_id)
      } else {
        this.dialog_message = 'サインインしてください。'
        this.if_show_dialog_2 = true
      }
    },
    async loadRooms() {
      this.loading = true
      const url = common.apiLink.rooms
      const resp = await common.httpGet(url, {})
      if (resp.result === 'ok') {
        this.rooms = resp.rooms
      } else {
        console.log(resp.message)
      }
      this.loading = false
    },
  },
}
</script>

<style>
h2 {
  color: #36479f;
}

iframe {
  border-width: 0;
}

.big-char {
  font-size: 2rem;
  color: #7f828b;
}
</style>
