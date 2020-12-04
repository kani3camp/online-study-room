<template>
  <v-navigation-drawer
    v-model="drawer"
    app
    right
    disable-resize-watcher
    overlay-color="#36479f"
  >
    <v-list dense>
      <v-list-item
        link
        @click="goToTopPage"
      >
        <v-list-item-action>
          <v-icon>mdi-home</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>サイトトップ</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToAllRoomsPage"
      >
        <v-list-item-action>
          <v-icon>mdi-door-open</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>ルーム一覧</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToAboutServicePage"
      >
        <v-list-item-action>
          <v-icon>mdi-account-question</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>はじめての方</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToYoutubeLive"
      >
        <v-list-item-action>
          <v-icon>mdi-youtube</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>ライブ配信</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToNewsPage"
      >
        <v-list-item-action>
          <v-icon>mdi-bulletin-board</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>お知らせ</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToTermsOfServicePage"
      >
        <v-list-item-action>
          <v-icon>mdi-file-document-outline</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>利用規約</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToPrivacyPolicyPage"
      >
        <v-list-item-action>
          <v-icon>mdi-account-lock</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>プライバシーポリシー</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        link
        @click="goToContactFormPage"
      >
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

<script>
import common from '@/plugins/common'

export default {
  name: 'NavigationDrawer',
  data: () => ({
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
  methods: {
    CloseDrawer() {
      this.drawer = false
    },
    goToTopPage() {
      this.$router.push('/')
      this.CloseDrawer()
    },
    goToAllRoomsPage() {
      this.$router.push('/all_rooms')
      this.CloseDrawer()
    },
    goToAboutServicePage() {
      this.$router.push('/about_service')
      this.CloseDrawer()
    },
    goToYoutubeLive() {
      window.open(this.youtubeLink, '_blank')
      this.CloseDrawer()
    },
    goToNewsPage() {
      this.$router.push('/news')
      this.CloseDrawer()
    },
    goToTermsOfServicePage() {
      this.$router.push('/terms_of_service')
      this.CloseDrawer()
    },
    goToPrivacyPolicyPage() {
      this.$router.push('/privacy_policy')
      this.CloseDrawer()
    },
    goToSettingsPage() {
      this.$router.push('/settings')
      this.CloseDrawer()
    },
    goToContactFormPage() {
      this.$router.push('/contact_form')
      this.CloseDrawer()
    },
  },
}
</script>

<style scoped></style>
