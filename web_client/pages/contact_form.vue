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

      <v-container v-show="!store.isSignedIn">
        <p>
          <nuxt-link to="/sign_in">サインイン</nuxt-link>するか、
          <a href="https://twitter.com/sorarideblog" target="_blank">@sorarideblog</a>
          までダイレクトメッセージを送ってください。
        </p>
      </v-container>

      <v-container v-show="store.isSignedIn">
        <v-form class="mx-auto">
          <v-select
            v-model="selected_contact_type"
            :items="contact_types"
            item-value="mode"
            item-text="ss"
            label="問い合わせの種類"
            outlined
          />
          <v-text-field v-model="mail_address" label="あなたのメールアドレス" outlined />
          <v-textarea v-model="message" label="本文" outlined />
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
      @cancel="if_show_dialog = false"
    />
  </v-app>
</template>

<script setup lang="ts">
// Regular expression from W3C HTML5.2 input specification:
// https://www.w3.org/TR/html/sec-forms.html#email-state-typeemail
import common from '~/plugins/common'
import NavigationDrawer from '#components'
import ToolBar from '#components'
import firebase from '~/plugins/firebase'
import Dialog from '#components'
import { useMainStore } from '~/stores'

const router = useRouter()
const store = useMainStore()

const dialog_message = ref('')
const message = ref<string | null>('')
const submitting = ref(false)
const selected_contact_type = ref(null)
const if_show_dialog = ref(false)
const contact_types = ref([
  { mode: 'feedback', ss: '意見' },
  { mode: 'contact', ss: '問い合わせ' },
])

const mail_address = computed(() => {
  if (store.isSignedIn) {
    return firebase.auth().currentUser?.email
  } else {
    return ''
  }
})

onMounted(() => {
  common.onAuthStateChanged()
})

const goToHomePage = () => router.push('/')
const goToSettingsPage = () => router.push('/settings')
const goToNewsPage = () => router.push('/news')
const submit = async () => {
  if (selected_contact_type || mail_address || message) {
    if_show_dialog.value = true
    dialog_message.value = '送信中'
    submitting.value = true

    const url = common.apiLink.send_contact_form
    const params = {
      mail_address: mail_address.toString(),
      user_id: firebase.auth().currentUser?.uid,
      id_token: await firebase.auth().currentUser?.getIdToken(false),
      contact_type: selected_contact_type,
      message: message,
    }
    const resp = await common.httpPost(url, params)
    submitting.value = false

    if (resp.result === 'ok') {
      message.value = null
      selected_contact_type.value = null
      dialog_message.value = '送信が完了しました。お問い合わせ頂きありがとうございます。'
    } else {
      dialog_message.value = '送信に失敗しました。'
    }
  } else {
    alert('未記入の項目があります。')
  }
}
</script>

<style scoped></style>
