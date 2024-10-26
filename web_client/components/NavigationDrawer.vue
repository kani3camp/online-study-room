<template>
  <v-navigation-drawer v-model="drawer" app right disable-resize-watcher overlay-color="#36479f">
    <v-list dense>
      <v-list-item link @click="goToTopPage">
        <v-list-item-action>
          <v-icon>mdi-home</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>サイトトップ</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item link @click="goToAboutServicePage">
        <v-list-item-action>
          <v-icon>mdi-account-question</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>このサイトについて</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <!-- <v-list-item
        link
        @click="goToYoutubeLive"
      >
        <v-list-item-action>
          <v-icon>mdi-youtube</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>ライブ配信</v-list-item-title>
        </v-list-item-content>
      </v-list-item> -->

      <v-list-item link @click="goToNewsPage">
        <v-list-item-action>
          <v-icon>mdi-bulletin-board</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>お知らせ</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item link @click="goToTermsOfServicePage">
        <v-list-item-action>
          <v-icon>mdi-file-document-outline</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>利用規約</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item link @click="goToPrivacyPolicyPage">
        <v-list-item-action>
          <v-icon>mdi-account-lock</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>プライバシーポリシー</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item link @click="goToContactFormPage">
        <v-list-item-action>
          <v-icon>mdi-email</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>ご意見・お問い合わせ</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import common from '~/plugins/common'
import { useMainStore } from '#imports'

// ルーターとストアのインスタンスを取得
const router = useRouter()
const store = useMainStore()

// プロパティの定義
const youtubeLink = common.key.youtubeLink

// drawerの状態をPiniaストアと同期させるコンピューテッドプロパティ
const drawer = computed({
  get: () => store.drawer,
  set: (value: boolean) => store.setDrawer(value),
})

// ダイアログを閉じる関数
const closeDrawer = () => {
  store.setDrawer(false)
}

// 各ページへの遷移関数
const goToTopPage = async () => {
  await router.push('/')
  closeDrawer()
}

const goToAboutServicePage = async () => {
  await router.push('/about_service')
  closeDrawer()
}

const goToYoutubeLive = () => {
  window.open(youtubeLink, '_blank')
  closeDrawer()
}

const goToNewsPage = async () => {
  await router.push('/news')
  closeDrawer()
}

const goToTermsOfServicePage = async () => {
  await router.push('/terms_of_service')
  closeDrawer()
}

const goToPrivacyPolicyPage = async () => {
  await router.push('/privacy_policy')
  closeDrawer()
}

const goToContactFormPage = async () => {
  await router.push('/contact_form')
  closeDrawer()
}
</script>

<style scoped></style>
