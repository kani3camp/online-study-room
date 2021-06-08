<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>ご意見・お問い合わせ</h2>
        </v-flex>
      </v-container>


      <v-container v-show="! ($store.state.isSignedIn)">
        <p>
          <nuxt-link to="/sign_in">サインイン</nuxt-link>するか、
          <a
            href="https://twitter.com/sorarideblog"
            target="_blank"
          >@sorarideblog</a>
          までダイレクトメッセージを送ってください。
        </p>
      </v-container>

      <v-container v-show="$store.state.isSignedIn">
        <v-form class="mx-auto">
          <v-select
            v-model="selected_contact_type"
            :items="contact_types"
            item-value="mode"
            item-text="ss"
            label="問い合わせの種類"
            outlined
          />
          <v-text-field
            v-model="mail_address"
            label="あなたのメールアドレス"
            outlined
          />
          <v-textarea
            v-model="message"
            label="本文"
            outlined
          />
          <div>
            <v-btn
              color="primary"
              block
              elevation="3"
              :disabled="submitting || !selected_contact_type || !message"
              @click="submit"
            >
              送信
            </v-btn>
          </div>
        </v-form>
      </v-container>
    </v-main>

    <Footer />

    <Dialog
      :if-show-dialog="if_show_dialog"
      :card-title="dialog_message"
      :loading="submitting"
      :accept-needed="false"
      cancel-option-string="閉じる"
      @cancel="if_show_dialog=false"
    />
  </v-app>
</template>

<script>
// Regular expression from W3C HTML5.2 input specification:
// https://www.w3.org/TR/html/sec-forms.html#email-state-typeemail
import common from '~/plugins/common'
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'
import firebase from '@/plugins/firebase'
import Dialog from '~/components/Dialog'

export default {
  name: 'ContactForm',
  components: {
    NavigationDrawer,
    ToolBar,
    Dialog,
  },
  data: () => ({
    // アロー関数でdataを定義している場合は中でthisがundefinedになるので注意
    dialog_message: '',
    message: '',
    submitting: false,
    selected_contact_type: null,
    if_show_dialog: false,
    contact_types: [
      { mode: 'feedback', ss: '意見' },
      { mode: 'contact', ss: '問い合わせ' },
    ],
  }),
  computed: {
    mail_address: {
      get() {
        if (this.$store.state.isSignedIn) {
          return firebase.auth().currentUser.email
        } else {
          return ''
        }
      },
    },
  },
  mounted() {
    common.onAuthStateChanged(this)
  },
  methods: {
    goToHomePage() {
      this.$router.push('/')
    },
    goToSettingsPage() {
      this.$router.push('/settings')
    },
    goToNewsPage() {
      this.$router.push('/news')
    },
    async submit() {
      if (this.selected_contact_type || this.mail_address || this.message) {
        this.if_show_dialog = true
        this.dialog_message = '送信中'
        this.submitting = true

        const url = common.apiLink.send_contact_form
        const params = {
          mail_address: this.mail_address.toString(),
          user_id: firebase.auth().currentUser.uid,
          id_token: await firebase.auth().currentUser.getIdToken(false),
          contact_type: this.selected_contact_type,
          message: this.message,
        }
        const resp = await common.httpPost(url, params)
        this.submitting = false

        if (resp.result === 'ok') {
          this.message = null
          this.selected_contact_type = null
          this.dialog_message = '送信が完了しました。お問い合わせ頂きありがとうございます。'
        } else {
          this.dialog_message = '送信に失敗しました。'
        }
      } else {
        alert('未記入の項目があります。')
      }
    },
  },
}
</script>

<style scoped></style>
