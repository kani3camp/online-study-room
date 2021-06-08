<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <!-- <v-container>
        <v-flex>
          <h2 style="display: inline-block">
            <v-icon>mdi-youtube</v-icon>
            YouTubeライブ
          </h2>
        </v-flex>
        <v-flex>
          <v-col>
            <v-row justify="center">
              <a
                target="_blank"
                :href="youtubeLink"
              ><h3>ライブ配信を見に行く</h3></a>
            </v-row>
          </v-col>
        </v-flex>
      </v-container> -->

      <v-container>
        <v-flex>
          <v-row>
            <v-flex>
              <v-col>
                <h2 style="display: inline-block">
                  <v-icon>mdi-door-open</v-icon>
                  ルーム
                </h2>
                に入室して作業開始！
              </v-col>
            </v-flex>
            <v-flex>
              <v-col class="d-flex flex-row-reverse">
                <v-btn
                  :disabled="loading"
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
                ロード中...
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
            >
              <v-hover v-slot="{ hover }">
                <v-card
                  class="ma-2 pa-3"
                  :elevation="hover ? 10 : 2"
                  @click="enterRoom(index)"
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
      </v-container>

      <Dialog
        :if-show-dialog="if_show_dialog"
        :card-title="dialog_message"
        :accept-needed="false"
        cancel-option-string="OK"
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
