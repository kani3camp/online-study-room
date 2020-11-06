<template>
  <v-app>
    <v-app-bar
      app
      flat
    >
      <v-btn @click="goToHomePage" icon><v-icon>mdi-home</v-icon></v-btn>
      <v-layout justify-center>
        <v-toolbar-title>設定</v-toolbar-title>
      </v-layout>
      <v-btn v-show="($store.state.isSignedIn)" @click="signOut" outlined>サインアウト</v-btn>

      <v-dialog v-model="if_show_dialog_2" width=500>
        <v-card>
          <v-card-title>{{ sign_out_result }}</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn @click="goToHomePage" text>閉じる</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

    </v-app-bar>

    <v-main>
      <v-layout justify-center>
      <v-list two-line subheader min-width="600px" max-width="800px">
        <v-subheader inset>設定</v-subheader>
        <v-list-item>
          <v-list-item-avatar>
            <v-icon></v-icon>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>表示名</v-list-item-title>
          </v-list-item-content>
          <v-list-item-content>
            <v-list-item-title>
              <v-text-field v-model="display_name"></v-text-field>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-avatar>
            <v-icon></v-icon>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>ひとこと</v-list-item-title>
          </v-list-item-content>
          <v-list-item-content>
            <v-list-item-title>
              <v-text-field v-model="status_message"></v-text-field>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

<!--        <v-list-item>-->
<!--          <v-list-item-avatar>-->
<!--            <v-icon></v-icon>-->
<!--          </v-list-item-avatar>-->

<!--          <v-list-item-content>-->
<!--            <v-list-item-title>パスワード</v-list-item-title>-->
<!--          </v-list-item-content>-->
<!--          <v-list-item-content>-->
<!--            <v-list-item-title><v-btn color="primary" @click="confirmChangingPassword" outlined>変更する</v-btn></v-list-item-title>-->
<!--          </v-list-item-content>-->
<!--        </v-list-item>-->

        <v-divider inset></v-divider>

        <v-subheader inset>情報</v-subheader>

        <v-list-item>
          <v-list-item-avatar>
            <v-icon></v-icon>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>ログイン中のアカウント</v-list-item-title>
          </v-list-item-content>
          <v-list-item-content>
            <v-list-item-title>{{ provider_id }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-avatar>
            <v-icon></v-icon>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>メールアドレス</v-list-item-title>
          </v-list-item-content>
          <v-list-item-content>
            <v-list-item-title>{{ mail_address }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-avatar>
            <v-icon></v-icon>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>合計学習時間</v-list-item-title>
          </v-list-item-content>
          <v-list-item-content>
            <v-list-item-title>{{ sum_study_time }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-avatar>
            <v-icon></v-icon>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>登録日</v-list-item-title>
          </v-list-item-content>
          <v-list-item-content>
            {{ registration_date }}
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-content>
            <v-btn color="primary" @click="saveNewValues" :disabled="!is_some_value_changed || is_some_value_blank || saving">保存</v-btn>
          </v-list-item-content>
        </v-list-item>

      </v-list>
      </v-layout>

    </v-main>

    <Footer></Footer>
  </v-app>
</template>

<script>
  import firebase from '../plugins/firebase'
  import common from "@/plugins/common";

  export default {
    name: "settings",
    data: () => ({
      display_name: null,
      status_message: null,
      mail_address: null,
      sum_study_time: null,
      registration_date: null,
      // if_show_dialog_1: false,
      if_show_dialog_2: false,
      sign_out_result: null,
      provider_id: null,
      saving: false,
    }),
    created() {
      this.display_name = this.$store.state.user.display_name
      this.status_message = this.$store.state.user.status_message
      this.mail_address = this.$store.state.user.mail_address
      this.sum_study_time = this.$store.state.user.sum_study_time
      this.registration_date = this.$store.state.user.registration_date
      this.provider_id = this.$store.state.user.provider_id
    },
    computed: {
      is_some_value_changed: function () {
        const bool1 = this.display_name !== this.$store.state.user.display_name
        const bool2 = this.status_message !== this.$store.state.user.status_message
        return bool1 || bool2;
      },
      is_some_value_blank: function () {
        return !this.display_name || !this.status_message
      }
    },
    methods: {
      goToHomePage() {
        this.$router.push('/')
      },
      async signOut() {
        const vm = this
        await firebase.auth().signOut().then(function() {
          console.log('Sign-out successful.')
          vm.sign_out_result = 'ログアウトしました。'
          vm.if_show_dialog_2 = true
        }).catch(function(error) {
          console.log(error)
          vm.sign_out_result = 'ログアウトに失敗しました。'
          vm.if_show_dialog_2 = true
        })
      },
      // confirmChangingPassword() {
      //   this.if_show_dialog_1 = true
      // },
      // sendPasswordChangeEmail() {
      //   this.if_show_dialog_1 = false
      // },
      async saveNewValues() {
        console.log('saveNewValues()')
        this.saving = true

        const url = 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/change_user_info'
        const params = {
          user_id: this.$store.state.user.user_id,
          id_token: this.$store.state.user.id_token,
          display_name: this.display_name,
          status_message: this.status_message,
        }
        const resp = await common.httpPost(url, params)
        if (resp.result === 'ok') {
          console.log('設定変更成功')
          this.$store.commit('user/setDisplayName', this.display_name)
          this.$store.commit('user/setStatusMessage', this.status_message)
        } else {
          console.log(resp)
          this.display_name = this.$store.state.user.display_name
          this.status_message = this.$store.state.user.status_message
        }
        this.saving = false
      }
    }
  }
</script>

<style scoped>

</style>
