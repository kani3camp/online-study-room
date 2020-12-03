import Vue from 'vue'
import firebase from 'firebase'
import { UserStore } from '../store'

declare module 'vue/types/vue' {
  interface Vue {
    $onAuthStateChanged(vm: any): void
    $key: { youtubeLink: string; twitterLink: string }
    $getUserData(vm: any): void
    $httpGet(urlStr: string, params: Object): string
    $httpPost(urlStr: string, _params: Object): string
  }
}

Vue.prototype.$key = {
  youtubeLink: 'https://www.youtube.com/channel/UCXuD2XmPTdpVy7zmwbFVZWg',
  twitterLink: 'https://twitter.com/osr_soraride',
}

Vue.prototype.$onAuthStateChanged = (vm: any) => {
  firebase.auth().onAuthStateChanged(async (user: any) => {
    if (user) {
      console.log('User is signed in.')
      UserStore.setSignInState(true)

      await Vue.prototype.$getUserData(vm)

      await firebase.auth().currentUser.getIdToken(true) // refresh idToken
    } else {
      console.log('User is signed out.')
      UserStore.signOut()
    }
  })
}

Vue.prototype.$getUserData = async (vm: any) => {
  const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/user_status')
  const params = {
    user_id: firebase.auth().currentUser.uid,
  }
  const userData = await Vue.prototype.$httpGet(url, params)
  if (userData.result !== 'ok') {
    console.log(userData)
  } else {
    const userBody = userData.user_status.user_body
    UserStore.setStatusMessage(userBody.status)
    UserStore.setTotalStudyTime(userBody.total_study_time)
    UserStore.setRegistrationDate(new Date(userBody.registration_date))
    UserStore.setLastEntered(new Date(userBody.last_entered))
  }
}

Vue.prototype.$httpGet = async (urlStr: string, params: Object) => {
  const url = new URL(urlStr)
  url.search = new URLSearchParams(params.toString()).toString()
  const response = await fetch(url.toString(), { method: 'GET' })
  return await response.json()
}

Vue.prototype.$httpPost = async (urlStr: string, _params: Object) => {
  const response = await fetch(urlStr, {
    method: 'POST',
    body: JSON.stringify(_params),
  })
  return await response.json()
}
