<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>サインイン</h2>
        </v-flex>
      </v-container>

      <v-container v-show="!store.isSignedIn">
        <p>登録・サインインともに下のボタンから行えます。</p>
        <p>
          <a href="/terms_of_service" target="_blank"> 利用規約 </a>
          および
          <a href="/privacy_policy" target="_blank"> プライバシーポリシー </a>
          に同意したうえで、当サービスをご利用ください。
        </p>
        <v-row id="google-sign-in-button" justify="center">
          <img
            :src="imageSource"
            alt="google sign in button"
            @click="signInWithGoogle"
            @mouseover="changeImageToHovered"
            @mousedown="changeImageToPressed"
            @mouseout="changeImageToNormal"
          />
        </v-row>
      </v-container>

      <v-container v-show="store.isSignedIn">
        <p>すでにサインイン済みです。</p>
        <p>アカウント設定ページからサインアウトできます。</p>
      </v-container>

      <Dialog
        :if-show-dialog="if_show_dialog"
        :card-title="dialog_message"
        :accept-needed="false"
        cancel-option-string="閉じる"
        @cancel="goToTopPage"
      />
    </v-main>

    <Footer />
  </v-app>
</template>

<script setup lang="ts">
import normalImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_normal_web@2x.png'
import hoveredImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_focus_web@2x.png'
import pressedImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_pressed_web@2x.png'
import firebase from '~/plugins/firebase'
import NavigationDrawer from '#components'
import ToolBar from '#components'
import Dialog from '#components'

const router = useRouter()
const store = useMainStore()

const imageSource = ref('')
const if_show_dialog = ref(false)
const dialog_message = ref('')

onMounted(() => {
  imageSource.value = normalImage
})

const goToTopPage = () => {
  router.push('/')
}
const changeImageToHovered = () => {
  console.log('changeImageToHovered()')
  imageSource.value = hoveredImage
}

const changeImageToPressed = () => {
  console.log('changeImageToPressed()')
  imageSource.value = pressedImage
}

const changeImageToNormal = () => {
  console.log('changeImageToNormal()')
  imageSource.value = normalImage
}

const signInWithGoogle = async () => {
  const provider = new firebase.auth.GoogleAuthProvider()
  await firebase
    .auth()
    .signInWithPopup(provider)
    .then(function () {
      dialog_message.value = 'ログインに成功しました。'
    })
    .catch(function (error) {
      console.log(error)
      dialog_message.value = 'ログインに失敗しました。'
    })
  if_show_dialog.value = true
}
</script>

<style scoped>
#google-sign-in-button img {
  max-height: 3rem;
}
</style>
