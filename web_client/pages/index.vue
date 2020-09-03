<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      app
    >
      <v-list dense>
        <v-list-item @click="drawer=false" link>
          <v-list-item-action>
            <v-icon>mdi-home</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>ホーム</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="goToSettingsPage" link>
          <v-list-item-action>
            <v-icon>mdi-account-cog</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>設定</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="goToContactFormPage" link>
          <v-list-item-action>
            <v-icon>mdi-email</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>ご意見・お問い合わせ</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="goToNewsPage" link>
          <v-list-item-action>
            <v-icon>mdi-bell</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>お知らせ</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      app
      color="cyan lighten-3"
      flat
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-layout justify-center>
        <v-toolbar-title><h1>オンライン作業部屋</h1></v-toolbar-title>
      </v-layout>
<!--      <v-btn v-show="!($store.state.isSignedIn)" @click="signInWithGoogle" outlined>Googleアカウントでログイン</v-btn>-->
      <img
        v-show="!($store.state.isSignedIn)"
        @click="signInWithGoogle"
        src="~assets/google_signin_buttons/web/2x/btn_google_signin_light_normal_web@2x.png"
        alt="sign in with google"
        height="50" width="200"/>
<!--      <v-btn v-show="!($store.state.isSignedIn)" outlined>登録</v-btn>-->
      <v-btn v-show="$store.state.isSignedIn" @click="goToSettingsPage" icon><v-icon>mdi-account-cog</v-icon></v-btn>
    </v-app-bar>

    <v-main>
      <v-container v-show="loading"
        class="fill-height"
        fluid
      >
        <v-row
          align="center"
          justify="center"
        >
          <v-col class="text-center">
            <div class="big-char">Loading...</div>
          </v-col>
        </v-row>
      </v-container>

      <v-container v-show="! loading">
        <v-row>
          <v-col
            cols="12" sm="4" md="3" lg="3" xl="3"
            v-for="(room, index) in rooms"
            :key="room.room_id"
            @click="confirmEntering(index)"
            dense
          >
            <v-hover v-slot:default="{ hover }">
              <v-card class="ma-2 pa-3" :elevation="hover ? 10 : 2">
                <v-layout justify-center>
                <v-card-title>
                  {{ room.room_body.name }}
                </v-card-title>
                </v-layout>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-container>

      <v-dialog v-model="if_show_dialog" width=500>
        <v-card class="mx-auto" outlined :loading="entering">
          <v-card-title>{{ selected_room_name }}の部屋 に入室しますか？</v-card-title>

          <v-card-actions>
            <v-row justify="end">
              <v-btn :disabled="entering" @click="enterRoom" text color="primary">入室する</v-btn>
              <v-btn :disabled="entering" @click="if_show_dialog=false" text>キャンセル</v-btn>
            </v-row>
          </v-card-actions>

        </v-card>
      </v-dialog>


      <v-dialog v-model="if_show_dialog_2" width=500>
        <v-card class="mx-auto" outlined>
          <v-card-title>{{ dialog_message }}</v-card-title>

          <v-card-actions>
            <v-row justify="end">
              <v-btn @click="if_show_dialog_2=false" text>閉じる</v-btn>
            </v-row>
          </v-card-actions>

        </v-card>
      </v-dialog>

    </v-main>
    <Footer></Footer>
  </v-app>
</template>

<script>
  import common from "~/plugins/common"
  import firebase from '../plugins/firebase'

  export default {
    data: () => ({
      drawer: null,
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
      common.onAuthStateChanged(this)

      this.loading = true
      const url = new URL('https://us-central1-online-study-room-f1f30.cloudfunctions.net/Rooms')
      const response = await fetch(url.toString())
      const resp = await response.json()
      if (resp.result === 'ok') {
        this.rooms = resp.rooms
      } else {
        console.log(resp.message)
      }
      this.loading = false
    },
    methods: {
      goToSettingsPage() {
        this.$router.push('/settings')
      },
      goToContactFormPage() {
        this.$router.push('/contact_form')
      },
      goToNewsPage() {
        this.$router.push('/news')
      },
      signInWithGoogle() {
        const vm = this
        const provider = new firebase.auth.GoogleAuthProvider()
        firebase.auth().signInWithPopup(provider).then(function(result) {
          let token = result.credential.accessToken;
          let user = result.user;
          console.log(user)
          vm.dialog_message = 'ログインに成功しました。'
        }).catch(function(error) {
          let errorCode = error.code
          let errorMessage = error.message
          let email = error.email
          // The firebase.auth.AuthCredential type that was used.
          let credential = error.credential
          vm.dialog_message = 'ログインに失敗しました。'
        })
        vm.if_show_dialog_2 = true
      },
      confirmEntering(index) {
        this.selected_index = index
        this.selected_room_name = this.rooms[this.selected_index].room_body.name
        this.if_show_dialog = true
      },
      async enterRoom() {
        if (this.$store.state.isSignedIn) {

          this.entering = true

          const selected_room_id = this.rooms[this.selected_index].room_id
          const url = 'https://us-central1-online-study-room-f1f30.cloudfunctions.net/EnterRoom'
          const params = new URLSearchParams({
            user_id: this.$store.state.user.user_id,
            room_id: selected_room_id,
            id_token: this.$store.state.user.id_token,
          })
          const res = await fetch(url, {
            method: 'POST',
            body: params
          }).then(response => response.json())

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
          this.dialog_message = 'ログインしてください。'
          this.if_show_dialog_2 = true
        }
      }
    }
  }
</script>

<style>
  .big-char {
    font-size: 3rem;
    color: #7F828B;
  }
</style>
