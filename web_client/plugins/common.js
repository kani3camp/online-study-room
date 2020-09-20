import firebase from "~/plugins/firebase"

const common = {
}

common.c = (m) => {
  console.log(m)
}

common.onAuthStateChanged = (vm) => {
  // const vm = this
  firebase.auth().onAuthStateChanged(async (user) => {
    if (user) {
      vm.$store.commit('user/setMailAddress', user.email)
      vm.$store.commit('user/setUserId', user.uid)
      vm.$store.commit('user/setProviderId', user.providerData[0].providerId)
      vm.$store.commit('user/setDisplayName', user.displayName)

      console.log('User is signed in.')
      vm.$store.commit('setSignInState', true)

      await common.getUserData(vm)

      firebase.auth().currentUser.getIdToken(/* forceRefresh */ true).then(function(idToken) {
        vm.$store.commit('user/setIdToken', idToken)
      }).catch(function(error) {
        console.error(error)
      })
    } else {
      console.log('User is signed out.')
      vm.$store.commit('signOut')
    }
  })
}


common.getUserData = async (vm) => {
  const url = new URL('https://us-central1-online-study-room-f1f30.cloudfunctions.net/UserStatus')
  const params = { user_id: vm.$store.state.user.user_id }
  // url.search = new URLSearchParams(params).toString()
  // const response = (await fetch(url.toString(), {method: 'GET'}))
  // const user_data = await response.json()
  const user_data = await common.httpGet(url, params)
  if (user_data.result !== 'ok') {
    console.log(user_data)
  } else {
    const user_body = user_data['user_status']['user_body']
    vm.$store.commit('user/setStatusMessage', user_body.status)
    // this.$store.commit('user/setSumStudyTime', use) // Todo
    vm.$store.commit('user/setRegistrationDate', new Date(user_body.registration_date))
    vm.$store.commit('user/setLastEntered', new Date(user_body.last_entered))
  }
}

common.httpGet = async (url_str, params) => {
  const url = new URL(url_str)
  url.search = new URLSearchParams(params).toString()
  const response = (await fetch(url.toString(), {method: 'GET'}))
  return await response.json()
}

common.httpPost = async (url_str, _params) => {
  const params = new URLSearchParams(_params)
  return await fetch(url_str, {
    method: 'POST',
    body: params
  }).then(r => r.json())
}

export default common
