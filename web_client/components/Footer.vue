<template>
  <v-footer
    id="footer"
    app
    absolute
  >
    <v-container>
      <v-row>
        <v-col
          cols="12"
          sm="6"
          md="4"
          lg="3"
          xl="3"
          class="child-center"
        >
          <h4>SERVICE</h4>
          <ul>
            <li>
              <nuxt-link to="/">
                サイトトップ
              </nuxt-link>
            </li>
            <li>
              <nuxt-link to="/about_service">
                このサイトについて
              </nuxt-link>
            </li>
            <li>
              <nuxt-link to="/news">
                お知らせ
              </nuxt-link>
            </li>
            <li>
              <nuxt-link to="/terms_of_service">
                利用規約
              </nuxt-link>
            </li>
            <li>
              <nuxt-link to="/privacy_policy">
                プライバシーポリシー
              </nuxt-link>
            </li>
          </ul>
        </v-col>

        <v-col
          cols="12"
          sm="6"
          md="4"
          lg="3"
          xl="3"
          class="child-center"
        >
          <h4>APP</h4>
          <ul>
            <li>Android</li>
          </ul>
        </v-col>

        <v-col
          cols="12"
          sm="6"
          md="4"
          lg="3"
          xl="3"
          class="child-center"
        >
          <h4>FOLLOW</h4>
          <ul>
            <li>
              <a
                :href="youtubeLink"
                target="_blank"
                rel="noopener noreferrer"
              ><v-icon>mdi-youtube</v-icon></a>
            </li>
            <li>
              <a
                :href="twitterLink"
                target="_blank"
                rel="noopener noreferrer"
              ><v-icon>mdi-twitter</v-icon></a>
            </li>
          </ul>
        </v-col>

        <v-col
          cols="12"
          sm="6"
          md="4"
          lg="3"
          xl="3"
          class="child-center"
        >
          <h4>ABOUT US</h4>
          <ul>
            <li>
              <a
                href="https://twitter.com/sorarideblog"
                target="_blank"
                rel="noopener noreferrer"
              >開発・運営者について</a>
            </li>
            <li>
              <nuxt-link to="/contact_form">
                お問い合わせ
              </nuxt-link>
            </li>
          </ul>
        </v-col>
      </v-row>

      <v-row justify="center">
        <span>&copy; Online Study Room {{ new Date().getFullYear() }}</span>
      </v-row>
    </v-container>
  </v-footer>
</template>

<script>
import common from '@/plugins/common'

export default {
  name: 'Footer',
  data: () => ({
    youtubeLink: common.key.youtubeLink,
    twitterLink: common.key.twitterLink,
  }),
  methods: {},
}
</script>

<style scoped>
#footer {
  background-color: #d3e6f6;
}

.child-center > * {
  text-align: center;
  justify-content: center;
  list-style-type: none;
}

.child-center h4 {
  color: #006cb8;
  font-family: sans-serif;
}

.child-center ul li a {
  color: #647a7f;
}

.child-center ul li a,
:visited {
  color: #647a7f;
}
.child-center ul li :hover {
  color: #a8babe;
}

.child-center ul li a {
  text-decoration: none;
}
</style>
