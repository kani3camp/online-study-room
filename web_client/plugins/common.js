import firebase from '~/plugins/firebase'

const common = {
  key: {
    youtubeLink: 'https://www.youtube.com/channel/UCXuD2XmPTdpVy7zmwbFVZWg',
    twitterLink: 'https://twitter.com/osr_soraride',
  },
}

common.c = (m) => {
  console.log(m)
}

common.onAuthStateChanged = (vm) => {
  firebase.auth().onAuthStateChanged(async (user) => {
    if (user) {
      console.log('User is signed in.')
      vm.$store.commit('setSignInState', true)

      await common.getUserData(vm)

      await firebase.auth().currentUser.getIdToken(true) // refresh idToken
    } else {
      console.log('User is signed out.')
      vm.$store.commit('signOut')
    }
  })
}

common.getUserData = async (vm) => {
  const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/user_status')
  const params = {
    user_id: firebase.auth().currentUser.uid,
  }
  const user_data = await common.httpGet(url, params)
  if (user_data.result !== 'ok') {
    console.log(user_data)
  } else {
    const user_body = user_data['user_status']['user_body']
    vm.$store.commit('user/setStatusMessage', user_body.status)
    vm.$store.commit('user/setTotalStudyTime', user_body.total_study_time)
    vm.$store.commit('user/setRegistrationDate', new Date(user_body.registration_date))
    vm.$store.commit('user/setLastEntered', new Date(user_body.last_entered))
  }
}

common.httpGet = async (url_str, params) => {
  const url = new URL(url_str)
  url.search = new URLSearchParams(params).toString()
  const response = await fetch(url.toString(), { method: 'GET' })
  return await response.json()
}

common.httpPost = async (url_str, _params) => {
  const response = await fetch(url_str, {
    method: 'POST',
    body: JSON.stringify(_params),
  })
  return await response.json()
}

export default common
