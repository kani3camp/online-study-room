<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      app
    >
      <v-list dense>
        <v-list-item @click="goToHomePage" link>
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
        <v-list-item @click="drawer=false" link>
          <v-list-item-action>
            <v-icon>mdi-email</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>ご意見・お問い合わせ</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar
      app
      flat
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-layout justify-center>
        <v-toolbar-title>利用規約</v-toolbar-title>
      </v-layout>
    </v-app-bar>

    <v-main></v-main>
  </v-app>
</template>

<script>
export default {
  name: "terms_of_service",
  data: () => ({
    drawer: null
  }),
  methods: {
    goToHomePage() {
      this.$router.push('/')
    },
    goToSettingsPage() {
      this.$router.push('/settings')
    },
  }
}
</script>

<style scoped>

</style>
