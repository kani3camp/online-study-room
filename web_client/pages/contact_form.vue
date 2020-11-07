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
        サインインしてください。
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

    <v-dialog
      v-model="if_show_dialog"
      width="500"
    >
      <v-card>
        <v-card-title>{{ dialog_message }}</v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn
            text
            @click="if_show_dialog=false"
          >
            閉じる
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
// Regular expression from W3C HTML5.2 input specification:
// https://www.w3.org/TR/html/sec-forms.html#email-state-typeemail
import common from '~/plugins/common'
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'

// const emailRegExp = /^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/

export default {
  name: 'ContactForm',
  components: {
    NavigationDrawer,
    ToolBar,
  },
  data: () => ({
    // アロー関数でdataを定義している場合は中でthisがundefinedになるので注意
    dialog_message: '',
    message: null,
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
        return this.$store.state.user.mail_address
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
        this.submitting = true

        const url =
          'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/send_contact_form'
        const params = {
          mail_address: this.mail_address.toString(),
          user_id: this.$store.state.user.user_id,
          id_token: this.$store.state.user.id_token,
          contact_type: this.selected_contact_type,
          message: this.message,
        }
        const resp = await common.httpPost(url, params)

        if (resp.result === 'ok') {
          this.message = null
          this.selected_contact_type = null
          this.dialog_message =
            '送信が完了しました。お問い合わせ頂きありがとうございます。'
          this.if_show_dialog = true
        } else {
          this.dialog_message = '送信に失敗しました。'
          this.if_show_dialog = true
        }
        this.submitting = false
      } else {
        alert('未記入の項目があります。')
      }
    },
  },
}
</script>

<style scoped></style>
