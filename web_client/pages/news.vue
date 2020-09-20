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
        <v-list-item @click="goToContactFormPage" link>
          <v-list-item-action>
            <v-icon>mdi-email</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>ご意見・お問い合わせ</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="drawer=false" link>
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
      flat
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-layout justify-center>
        <v-toolbar-title>お知らせ</v-toolbar-title>
      </v-layout>
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
        <v-row dense>
          <v-col v-for="(news, index) in newsList" :key="news.news_id" cols="12" dense>
            <v-card>
              <v-card-title v-text="news.news_body.title"></v-card-title>
              <v-card-subtitle v-text="formatDateString(news.news_body.updated)"></v-card-subtitle>
              <v-card-text v-text="news.news_body.text_body"></v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>

    </v-main>

    <Footer></Footer>
  </v-app>
</template>

<script>
import common from "@/plugins/common";

export default {
  name: "news",
  data: () => ({
    drawer: null,
    newsList: [],
    loading: false,
  }),
  async created() {
    this.loading = true
    const url = 'https://us-central1-online-study-room-f1f30.cloudfunctions.net/News'
    const params = {
      num_news: 10
    }
    const resp = await common.httpGet(url, params)
    if (resp.result === 'ok') {
      this.newsList = resp['news_list']
    } else {
      console.log(resp.message)
    }
    this.loading = false

  },
  methods: {
    goToHomePage() {
      this.$router.push('/')
    },
    goToSettingsPage() {
      this.$router.push('/settings')
    },
    goToContactFormPage() {
      this.$router.push('/contact_form')
    },
    formatDateString(_date) {
      const date = new Date(_date)
      const y = date.getFullYear()
      const mo = date.getMonth() + 1
      const d = date.getDate()
      const h = date.getHours()
      const mi = date.getMinutes()
      return `${y}/${mo}/${d} ${h}:${mi}`
    }
  }
}
</script>

<style scoped>

</style>
