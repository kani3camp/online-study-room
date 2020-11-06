<template>
  <v-footer
    app
  >
    <v-container>
      <v-row>
        <v-col cols="auto">
          <h4>SERVICE</h4>
          <ul>
            <li>サイトトップ</li>
            <li>このサイトについて</li>
            <li>お知らせ</li>
            <li>利用規約</li>
            <li>プライバシーポリシー</li>
          </ul>
        </v-col>

        <v-col cols="auto">
          <h4>APP</h4>
          <ul>
            <li>Android</li>
          </ul>
        </v-col>

        <v-col cols="auto">
          <h4>FOLLOW</h4>
          <ul>
            <li>
              <a v-bind:href="youtubeLink"><v-icon>mdi-youtube</v-icon></a>
            </li>
            <li>
              <a v-bind:href="twitterLink"><v-icon>mdi-twitter</v-icon></a>
            </li>
          </ul>
        </v-col>

        <v-col cols="auto">
          <h4>ABOUT US</h4>
          <ul>
            <li>開発・運営者について</li>
            <li>お問い合わせ</li>
          </ul>
        </v-col>
      </v-row>

      <v-row justify="center"><span>&copy; Online Study Room {{ new Date().getFullYear() }}</span></v-row>
    </v-container>

  </v-footer>

</template>

<script>
import common from "@/plugins/common";

    export default {
      name: "Footer",
      data: () => ({
        youtubeLink: common.key.youtubeLink,
        twitterLink: common.key.twitterLink,
      }),
      methods: {
      }
    }
</script>

<style scoped>

</style>
