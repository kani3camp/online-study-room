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
        <v-list-item @click="drawer=false" link>
          <v-list-item-action>
            <v-icon>mdi-email</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>ご意見・お問い合わせ</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="goToNewsPage" link>
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
        <v-toolbar-title>意見・お問い合わせ</v-toolbar-title>
      </v-layout>
    </v-app-bar>

    <v-main>
    <v-form class="mx-auto">
      <v-container>
        <v-select
          v-model="selected_contact_type"
          :items="contact_types"
          item-value="mode"
          item-text="ss"
          label="問い合わせの種類"
          outlined
        ></v-select>
        <v-text-field
          v-model="mail_address"
          label="あなたのメールアドレス"
          outlined
        ></v-text-field>
        <v-textarea label="本文" v-model="message" outlined></v-textarea>
        <div>
          <v-btn
            @click="submit"
            color="primary"
            block elevation="3"
            :disabled="submitting"
          >送信
          </v-btn>
        </div>
      </v-container>
    </v-form>
    </v-main>

    <Footer></Footer>

    <v-dialog v-model="if_show_dialog" width=500>
      <v-card>
        <v-card-title>{{ dialog_message }}</v-card-title>
        <v-card-actions>
          <v-row justify="end">
            <v-btn @click="if_show_dialog=false" text>閉じる</v-btn>
          </v-row>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-app>
</template>

<script>
// Regular expression from W3C HTML5.2 input specification:
// https://www.w3.org/TR/html/sec-forms.html#email-state-typeemail
import common from "~/plugins/common"

const emailRegExp = /^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/

export default {
  name: "contact_form",
  data: function () {
    return {
      drawer: null,
      mail_address: this.$store.state.user.mail_address,
      message: null,
      submitting: false,
      contact_types: [
        {mode: 'feedback', ss: '意見'},
        {mode: 'contact', ss: '問い合わせ'},
      ],
      selected_contact_type: null,
      if_show_dialog: false,
      dialog_message: null,
    };
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

        const url = new URL('https://us-central1-online-study-room-f1f30.cloudfunctions.net/SendContactForm')
        const params = new URLSearchParams({
          mail_address: this.mail_address,
          user_id: this.$store.state.user.user_id,
          id_token: this.$store.state.user.id_token,
          contact_type: this.selected_contact_type,
          message: this.message,
        })
        const resp = await fetch(url.toString(), {
          method: 'POST',
          body: params
        }).then(response => response.json())

        if (resp.result === 'ok') {
          this.message = null
          this.selected_contact_type = null
          this.dialog_message = '送信が完了しました。お問い合わせ頂きありがとうございます。'
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

<style scoped>
</style>
