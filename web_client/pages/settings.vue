<template>
  <v-app>
    <v-app-bar app flat>
      <v-btn icon @click="goToHomePage">
        <v-icon>mdi-home</v-icon>
      </v-btn>
      <v-layout justify-center>
        <v-toolbar-title>設定</v-toolbar-title>
      </v-layout>
      <v-btn v-show="store.isSignedIn" outlined @click="signOut"> サインアウト </v-btn>
    </v-app-bar>

    <v-main>
      <v-card class="mx-auto" max-width="500px" outlined>
        <v-list id="setting-list" two-line subheader>
          <v-subheader>設定</v-subheader>
          <v-list-item two-line>
            <v-list-item-content>
              <v-list-item-title>表示名</v-list-item-title>
              <v-list-item-subtitle>
                <v-text-field v-model="display_name" />
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>

          <v-list-item two-line>
            <v-list-item-content>
              <v-list-item-title>ひとこと</v-list-item-title>
              <v-list-item-subtitle>
                <v-text-field v-model="status_message" />
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>

          <v-divider />

          <v-subheader>情報</v-subheader>

          <v-list-item two-line>
            <v-list-item-content>
              <v-list-item-title>ログイン中のアカウント</v-list-item-title>
              <v-list-item-subtitle class="text-center">
                {{ provider_id }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>

          <v-list-item two-line>
            <v-list-item-content>
              <v-list-item-title>メールアドレス</v-list-item-title>
              <v-list-item-subtitle class="text-center">
                {{ mail_address }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>

          <v-list-item two-line>
            <v-list-item-content>
              <v-list-item-title>合計学習時間</v-list-item-title>
              <v-list-item-subtitle class="text-center">
                {{ total_study_time }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>

          <v-list-item two-line>
            <v-list-item-content>
              <v-list-item-title>登録日</v-list-item-title>
              <v-list-item-subtitle class="text-center">
                {{ registration_date_str }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>

          <v-list-item>
            <v-list-item-content>
              <v-btn
                color="primary"
                :disabled="!is_some_value_changed || is_some_value_blank || saving"
                @click="saveNewValues"
              >
                保存
              </v-btn>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card>

      <Dialog
        :if-show-dialog="if_show_dialog"
        :card-title="dialog_message"
        :loading="saving"
        :accept-needed="false"
        cancel-option-string="閉じる"
        @cancel="goToHomePage"
      />
    </v-main>

    <Footer />
  </v-app>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import firebase from '../plugins/firebase'
import common from '~/plugins/common'
import Dialog from '#components'

const router = useRouter()
const store = useMainStore()
const userStore = useUserStore()

const display_name = ref<string | null>(null)
const status_message = ref<string | null>(null)
const if_show_dialog = ref(false)
const dialog_message = ref<string | null>(null)
const saving = ref(false)

const is_some_value_changed = computed(() => {
  const bool1 = display_name !== firebase_display_name
  const bool2 = status_message !== firebase_status_message
  return bool1 || bool2
})

const is_some_value_blank = computed(() => {
  return !display_name || !status_message
})

const firebase_display_name = computed(() => {
  return firebase.auth().currentUser?.displayName
})

const firebase_status_message = computed(() => {
  return userStore.status_message
})

const registration_date_str = computed(() => {
  const registration_date = userStore.registration_date
  if (registration_date) {
    return (
      registration_date.getFullYear() +
      '年' +
      (registration_date.getMonth() + 1) +
      '月' +
      registration_date.getDate() +
      '日'
    )
  } else {
    return null
  }
})

const mail_address = computed(() => {
  return firebase.auth().currentUser?.email
})

const provider_id = computed(() => {
  return firebase.auth().currentUser?.providerData[0]?.providerId
})

const total_study_time = computed(() => {
  const total_seconds = userStore.total_study_time
  if (total_seconds !== null) {
    const hours = Math.floor(total_seconds / 3600)
    const total_minutes = Math.floor(total_seconds / 60)
    const minutes = total_minutes % 60
    return hours + '時間' + minutes + '分'
  }
  return null
})

watch(firebase_status_message, (newValue, oldValue) => {
  console.log('watch: ' + display_name.value + ': new: ' + newValue + ', old: ' + oldValue)
  if (oldValue === null && newValue !== null) {
    display_name.value = newValue
  } else if (oldValue !== newValue && newValue !== null) {
    display_name.value = newValue
  }
})

watch(firebase_status_message, (newValue, oldValue) => {
  console.log('watch')
  if ((oldValue === null && newValue !== null) || oldValue !== newValue) {
    status_message.value = newValue
  } else if (oldValue !== newValue && newValue !== null) {
    display_name.value = newValue
  }
})

onMounted(() => {
  common.onAuthStateChanged()
})
onMounted(() => {
  display_name.value = firebase_display_name.value
  status_message.value = firebase_status_message.value
})

const goToHomePage = () => {
  router.push('/')
}

const signOut = async () => {
  await firebase
    .auth()
    .signOut()
    .then(function () {
      console.log('Sign-out successful.')
      dialog_message.value = 'サインアウトしました。'
      if_show_dialog.value = true
    })
    .catch(function (error) {
      console.log(error)
      dialog_message.value = 'サインアウトに失敗しました。'
      if_show_dialog.value = true
    })
}

const saveNewValues = async () => {
  console.log('saveNewValues()')
  saving.value = true
  dialog_message.value = '保存中'
  if_show_dialog.value = true

  const url = common.apiLink.change_user_info
  const params = {
    user_id: firebase.auth().currentUser?.uid,
    id_token: await firebase.auth().currentUser?.getIdToken(false),
    display_name: display_name,
    status_message: status_message,
  }
  const resp = await common.httpPost(url, params)
  if (resp.result === 'ok') {
    dialog_message.value = '完了！'
    const new_display_name = display_name
    await firebase.auth().currentUser?.updateProfile({
      displayName: new_display_name.value ?? '',
    })
    userStore.setStatusMessage(status_message.value)
  } else {
    console.log(resp)
    dialog_message.value = 'エラー。もう一度試してみてください。'
    display_name.value = firebase_display_name.value
    status_message.value = firebase_status_message.value
  }
  saving.value = false
}
</script>

<style scoped>
/* <v-list-item-title>の左padding */
#setting-list div div div div {
  padding-left: 1rem;
}
</style>
