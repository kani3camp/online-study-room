<template>
  <v-app-bar
    app
    flat
    color="white"
    min-height="60"
  >
    <v-layout fill-height pb-0 pt-0 pl-0>
      <nuxt-link to="/"><v-flex @click="goToTopPage" align-self-center><Logo></Logo></v-flex></nuxt-link>
    </v-layout>

    <v-spacer></v-spacer>

    <v-layout fill-height pb-0 pt-0 pl-0 wrap id="tool-right">
      <div class="tool-menu" v-show="!($vuetify.breakpoint.mobile)">
        <div class="tool-content"><nuxt-link to="/all_rooms">ルーム一覧</nuxt-link></div>
        <div class="tool-content"><nuxt-link to="/about_service">はじめての方</nuxt-link></div>
        <div class="tool-content"><nuxt-link to="/">YouTube</nuxt-link></div> <!-- todo -->
        <div class="tool-content"><nuxt-link to="/news">お知らせ</nuxt-link></div>
      </div>
      <div v-show="! ($store.state.isSignedIn)" class="tool-content">
        <v-btn v-show="!($store.state.isSignedIn)" @click="goToSignInPage" outlined>サインイン</v-btn>
      </div>
      <div v-show="($store.state.isSignedIn)">
        <v-btn @click="goToSettingsPage" icon><v-icon>mdi-account-cog</v-icon></v-btn>
      </div>
      <div v-show="($vuetify.breakpoint.mobile)">
        <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      </div>
    </v-layout>
  </v-app-bar>
</template>

<script>
import Logo from "@/components/Logo"

export default {
  name: "ToolBar",
  components: {
    Logo
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
  },
  methods: {
    goToTopPage() {
      this.$router.push('/')
    },
    goToYoutubeLive() {
      // todo
    },
    goToSettingsPage() {
      this.$router.push('/settings')
    },
    goToSignInPage() {
      this.$router.push('/sign_in')
    },
  }
}
</script>

<style scoped>

#tool-right {
  display: flex;
  justify-content: flex-end;
}

.tool-menu {
  display: flex;
}

.tool-content {
  margin: 0 0.5rem;
  align-self: center;
}
.tool-content a, :visited {
  color: #000;
  text-decoration: none;
}
.tool-content :hover {
  color: #006CB8;
}

</style>
