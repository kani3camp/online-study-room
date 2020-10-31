<template>
  <v-app-bar
    app
    flat
    color="gray"
    min-height="60"
  >
    <v-layout fill-height pb-0 pt-0 pl-0>
      <v-flex align-self-center><Logo></Logo></v-flex>
    </v-layout>

    <v-spacer></v-spacer>

    <v-layout fill-height pb-0 pt-0 pl-0 wrap id="tool-right">
      <div class="tool-menu">
        <div class="tool-content">ルーム一覧</div>
        <div class="tool-content">はじめての方</div>
        <div class="tool-content">YouTube</div>
        <div class="tool-content">お知らせ</div>

        <!--          <img-->
        <!--          v-show="!($store.state.isSignedIn)"-->
        <!--          @click="signInWithGoogle"-->
        <!--          src="~assets/google_signin_buttons/web/2x/btn_google_signin_light_normal_web@2x.png"-->
        <!--          alt="sign in with google"-->
        <!--          height="50" width="200"/>-->
        <div class="tool-content">
          <v-btn v-show="!($store.state.isSignedIn)" outlined>サインイン</v-btn>
        </div>
        <div>
          <v-btn v-show="$store.state.isSignedIn" @click="goToSettingsPage" icon><v-icon>mdi-account-cog</v-icon></v-btn>
        </div>
      </div>
      <div>
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
    goToSettingsPage() {
      this.$router.push('/settings')
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

</style>
