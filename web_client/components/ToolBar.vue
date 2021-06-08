<template>
  <v-app-bar
    app
    flat
    color="white"
    min-height="60"
  >
    <v-layout
      fill-height
      pb-0
      pt-0
      pl-0
    >
      <nuxt-link to="/">
        <v-flex
          align-self-center
          @click="goToTopPage"
        >
          <Logo />
        </v-flex>
      </nuxt-link>
    </v-layout>

    <v-spacer />

    <v-layout
      id="tool-right"
      fill-height
      pb-0
      pt-0
      pl-0
      wrap
    >
      <div
        v-show="!($vuetify.breakpoint.mobile)"
        class="tool-menu"
      >
        <div class="tool-content">
          <nuxt-link to="/about_service">
            このサイトについて
          </nuxt-link>
        </div>

        <!-- <div class="tool-content">
          <a
            :href="youtubeLink"
            target="_blank"
            rel="noopener noreferrer"
          >
            <v-icon>mdi-youtube</v-icon>
            YouTube
          </a> -->

        </div>
        <div class="tool-content">
          <nuxt-link to="/news">
            お知らせ
          </nuxt-link>
        </div>
      </div>
      <div
        v-show="! ($store.state.isSignedIn)"
        class="tool-content"
      >
        <v-btn
          v-show="!($store.state.isSignedIn)"
          outlined
          @click="goToSignInPage"
        >
          サインイン
        </v-btn>
      </div>
      <div v-show="($store.state.isSignedIn)">
        <v-btn
          icon
          @click="goToSettingsPage"
        >
          <v-icon>mdi-account-cog</v-icon>
        </v-btn>
      </div>
      <div v-show="($vuetify.breakpoint.mobile)">
        <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      </div>
    </v-layout>
  </v-app-bar>
</template>

<script>
import Logo from '@/components/Logo'
import common from '@/plugins/common'

export default {
  name: 'ToolBar',
  components: {
    Logo,
  },
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
    goToTopPage() {
      this.$router.push('/')
    },
    goToSettingsPage() {
      this.$router.push('/settings')
    },
    goToSignInPage() {
      this.$router.push('/sign_in')
    },
  },
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
.tool-content a,
:visited {
  color: #000;
  text-decoration: none;
}
.tool-content :hover {
  color: #006cb8;
}
</style>
