<template>
  <v-app>
    <v-app-bar
      app
      flat
    >
      <v-btn
        icon
        @click="goToHomePage"
      >
        <v-icon>mdi-home</v-icon>
      </v-btn>
      <v-layout justify-center>
        <v-toolbar-title>設定</v-toolbar-title>
      </v-layout>
      <v-btn
        v-show="$store.state.isSignedIn"
        outlined
        @click="signOut"
      >
        サインアウト
      </v-btn>

      <v-dialog
        v-model="if_show_dialog_2"
        width="500"
      >
        <v-card>
          <v-card-title>{{ sign_out_result }}</v-card-title>
          <v-card-actions>
            <v-spacer />
            <v-btn
              text
              @click="goToHomePage"
            >
              閉じる
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-app-bar>

    <v-main>
      <v-layout justify-center>
        <v-list
          id="setting-list"
          two-line
          subheader
        >
          <v-subheader>設定</v-subheader>
          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>表示名</v-list-item-title>
              </v-list-item-content>
            </v-flex>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>
                  <v-text-field v-model="display_name" />
                </v-list-item-title>
              </v-list-item-content>
            </v-flex>
          </v-list-item>

          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>ひとこと</v-list-item-title>
              </v-list-item-content>
            </v-flex>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>
                  <v-text-field v-model="status_message" />
                </v-list-item-title>
              </v-list-item-content>
            </v-flex>
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

          <v-divider />

          <v-subheader>情報</v-subheader>

          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>ログイン中のアカウント</v-list-item-title>
              </v-list-item-content>
            </v-flex>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>{{ provider_id }}</v-list-item-title>
              </v-list-item-content>
            </v-flex>
          </v-list-item>

          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>メールアドレス</v-list-item-title>
              </v-list-item-content>
            </v-flex>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>{{ mail_address }}</v-list-item-title>
              </v-list-item-content>
            </v-flex>
          </v-list-item>

          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>合計学習時間</v-list-item-title>
              </v-list-item-content>
            </v-flex>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>{{ total_study_time }}</v-list-item-title>
              </v-list-item-content>
            </v-flex>
          </v-list-item>

          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-list-item-title>登録日</v-list-item-title>
              </v-list-item-content>
            </v-flex>
            <v-flex>
              <v-list-item-content>
                {{ registration_date_str }}
              </v-list-item-content>
            </v-flex>
          </v-list-item>

          <v-list-item>
            <v-flex>
              <v-list-item-content>
                <v-btn
                  color="primary"
                  :disabled="
                    !is_some_value_changed || is_some_value_blank || saving
                  "
                  @click="saveNewValues"
                >
                  保存
                </v-btn>
              </v-list-item-content>
            </v-flex>
          </v-list-item>
        </v-list>
      </v-layout>
    </v-main>

    <Footer />
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import firebase from 'firebase'
// import common from '~/plugins/common'

export default Vue.extend({
  name: 'Settings',
  data: () => ({
    display_name: null,
    status_message: null,
    if_show_dialog_2: false,
    sign_out_result: null,
    saving: false,
  }),
  computed: {
    is_some_value_changed() {
      const bool1 = this.display_name !== this.firebase_display_name
      const bool2 = this.status_message !== this.firebase_status_message
      return bool1 || bool2
    },
    is_some_value_blank() {
      return !this.display_name || !this.status_message
    },
    firebase_display_name() {
      return firebase.auth().currentUser.displayName
    },
    firebase_status_message() {
      return this.$store.state.user.status_message
    },
    registration_date_str() {
      const registrationDate = this.$store.state.user.registration_date
      if (registrationDate) {
        return (
          registrationDate.getFullYear() +
          '年' +
          (registrationDate.getMonth() + 1) +
          '月' +
          registrationDate.getDate() +
          '日'
        )
      } else {
        return null
      }
    },
    mail_address() {
      return firebase.auth().currentUser.email
    },
    provider_id() {
      return firebase.auth().currentUser.providerData[0].providerId
    },
    total_study_time() {
      const totalSeconds = this.$store.state.user.total_study_time
      if (totalSeconds) {
        const hours = Math.floor(totalSeconds / 3600)
        const totalMinutes = Math.floor(totalSeconds / 60)
        const minutes = totalMinutes % 60
        return hours + '時間' + minutes + '分'
      }
      return null
    },
  },
  watch: {
    firebase_display_name(newValue, oldValue) {
      console.log('watch: ' + this.display_name + ': new: ' + newValue + ', old: ' + oldValue)
      if (oldValue === null && newValue !== null) {
        this.display_name = newValue
      } else if (oldValue !== newValue) {
        this.display_name = newValue
      }
    },
    firebase_status_message(newValue, oldValue) {
      console.log('watch')
      if ((oldValue === null && newValue !== null) || oldValue !== newValue) {
        this.status_message = newValue
      } else if (oldValue !== newValue) {
        this.display_name = newValue
      }
    },
  },
  async created() {
    await this.$onAuthStateChanged(this)
  },
  mounted() {
    this.display_name = this.firebase_display_name
    this.status_message = this.firebase_status_message
  },
  methods: {
    goToHomePage() {
      this.$router.push('/')
    },
    async signOut() {
      const vm = this
      await firebase
        .auth()
        .signOut()
        .then(function () {
          console.log('Sign-out successful.')
          vm.sign_out_result = 'サインアウトしました。'
          vm.if_show_dialog_2 = true
        })
        .catch(function (error) {
          console.log(error)
          vm.sign_out_result = 'サインアウトに失敗しました。'
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
        user_id: firebase.auth().currentUser.uid,
        id_token: await firebase.auth().currentUser.getIdToken(false),
        display_name: this.display_name,
        status_message: this.status_message,
      }
      const resp = await this.$httpPost(url, params)
      if (resp.result === 'ok') {
        console.log('設定変更成功')
        const newDisplayName = this.display_name
        await firebase.auth().currentUser.updateProfile({
          displayName: newDisplayName,
        })
        this.$store.commit('user/setStatusMessage', this.status_message)
      } else {
        console.log(resp)
        this.display_name = this.firebase_display_name
        this.status_message = this.firebase_status_message
      }
      this.saving = false
    },
  },
})
</script>

<style scoped>
/* <v-list-item-title>の左padding */
#setting-list div div div div {
  padding-left: 1rem;
}
</style>
