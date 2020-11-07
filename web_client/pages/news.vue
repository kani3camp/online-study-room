<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>お知らせ</h2>
        </v-flex>
      </v-container>

      <v-container
        v-show="loading"
        class="fill-height"
        fluid
      >
        <v-row
          align="center"
          justify="center"
        >
          <v-col class="text-center">
            <div class="big-char">
              Loading...
            </div>
          </v-col>
        </v-row>
      </v-container>

      <v-container v-show="! loading">
        <v-row dense>
          <v-col
            v-for="news in newsList"
            :key="news['news_id']"
            cols="12"
            dense
          >
            <v-card>
              <v-card-title v-text="news['news_body'].title" />
              <v-card-subtitle v-text="formatDateString(news['news_body'].updated)" />
              <v-card-text v-text="news['news_body']['text_body']" />
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>

    <Footer />
  </v-app>
</template>

<script>
import common from '@/plugins/common'
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'

export default {
  name: 'News',
  components: {
    NavigationDrawer,
    ToolBar,
  },
  data: () => ({
    drawer: null,
    newsList: [],
    loading: false,
  }),
  async created() {
    this.loading = true
    const url =
      'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/news'
    const params = {
      num_news: 10,
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
    },
  },
}
</script>

<style scoped></style>
