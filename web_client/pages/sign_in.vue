<template>
  <v-app>
    <NavigationDrawer></NavigationDrawer>

    <ToolBar></ToolBar>

    <v-main>
      <v-container>
        <v-flex>
          <h2>サインイン</h2>
        </v-flex>
      </v-container>

      <v-container v-show="! ($store.state.isSignedIn)">
        <p>　登録・サインインともに下のボタンから行えます。</p>
        <v-row justify="center" id="google-sign-in-button">
          <img @click="signInWithGoogle" v-bind:src="imageSource"
               v-on:mouseover="changeImageToHovered"
               v-on:mousedown="changeImageToPressed"
               v-on:mouseout="changeImageToNormal"
               alt="google sign in button">
        </v-row>
      </v-container>

      <v-container>
        <p>すでにサインイン済みです。</p>
        <p>アカウント設定ページからサインアウトできます。</p>
      </v-container>


      <v-dialog v-model="if_show_dialog_2" width=500>
        <v-card class="mx-auto" outlined>
          <v-card-title>{{ dialog_message }}</v-card-title>

          <v-card-actions>
            <v-row justify="end">
              <v-btn @click="if_show_dialog_2=false; goToTopPage" text>閉じる</v-btn>
            </v-row>
          </v-card-actions>

        </v-card>
      </v-dialog>


    </v-main>

    <Footer></Footer>
  </v-app>
</template>

<script>
import normalImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_normal_web@2x.png'
import hoveredImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_focus_web@2x.png'
import pressedImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_pressed_web@2x.png'
import firebase from "@/plugins/firebase";

export default {
  name: "sign_in",
  data: () => ({
    imageSource: '',
    if_show_dialog_2: false,
    dialog_message: '',
  }),
  created() {
    this.imageSource = normalImage
  },
  methods: {
    goToTopPage() {
      this.$router.push('/') // todo
    },
    changeImageToHovered() {
      console.log('changeImageToHovered()')
      this.imageSource = hoveredImage
    },
    changeImageToPressed() {
      console.log('changeImageToPressed()')
      this.imageSource = pressedImage
    },
    changeImageToNormal() {
      console.log('changeImageToNormal()')
      this.imageSource = normalImage
    },
    async signInWithGoogle() {
      const vm = this
      const provider = new firebase.auth.GoogleAuthProvider()
      await firebase.auth().signInWithPopup(provider).then(function(result) {
        let token = result.credential['accessToken']
        let user = result.user
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
  },
}
</script>

<style scoped>

#google-sign-in-button img {
  max-height: 3rem;
}

</style>
