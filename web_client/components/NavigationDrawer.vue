<template>
  <v-navigation-drawer
    v-model="drawer"
    app
    right
  >
    <v-list dense>
      <v-list-item @click="drawer = false" link>
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
</template>

<script>
export default {
  name: "NavigationDrawer",
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
  },
  computed: {
    drawer: {
      get() {
        return this.$store.state.drawer
      },
      set(value) {
        this.$store.commit('setDrawer', value)
      }
    }
  }
}
</script>

<style scoped>

</style>
