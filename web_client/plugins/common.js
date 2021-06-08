import firebase from '~/plugins/firebase'

const productionApiLink = {
  rooms: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms',
  room_status: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_status',
  online_users: '',
  user_status: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/user_status',
  change_user_info: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/change_user_info',
  send_contact_form: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/send_contact_form',
  create_new_room: '',
  news: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/news',
  room_layout: '',
  websocket: 'wss://0ieer51ju9.execute-api.ap-northeast-1.amazonaws.com/release',
}
const productionFirebaseConfig = {
  apiKey: 'AIzaSyCIY7RCArApaRpg4nkVUJCze88OmF6-wGM',
  authDomain: 'online-study-space.firebaseapp.com',
  databaseURL: 'https://online-study-space.firebaseio.com',
  projectId: 'online-study-space',
  storageBucket: 'online-study-space.appspot.com',
  messagingSenderId: '486366182751',
  appId: '1:486366182751:web:dc94794117a6268a050aea',
  measurementId: 'G-KH6EFLWF7T',
}

const testApiLink = {
  rooms: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_rooms',
  room_status: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_room_status',
  online_users: '',
  user_status: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_user_status',
  change_user_info: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_change_user_info',
  send_contact_form: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_send_contact_form',
  create_new_room: '',
  news: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_news',
  room_layout: '',
  websocket: 'wss://jukb7tixp7.execute-api.ap-northeast-1.amazonaws.com/test',
}
const testFirebaseConfig = {
  apiKey: 'AIzaSyAE5Hm9U0pQDxXAPNkmthO_z--CFDiSn6w',
  authDomain: 'test-online-study-space.firebaseapp.com',
  projectId: 'test-online-study-space',
  storageBucket: 'test-online-study-space.appspot.com',
  messagingSenderId: '788678710894',
  appId: '1:788678710894:web:2cb434fbd57cdc9bfd787f',
  measurementId: 'G-450CVZ6FWP',
}

const common = {
  key: {
    youtubeLink: 'https://www.youtube.com/channel/UCXuD2XmPTdpVy7zmwbFVZWg',
    twitterLink: 'https://twitter.com/osr_soraride',
  },
  // todo デプロイ時に変更
  apiLink: productionApiLink,
  firebaseConfig: productionFirebaseConfig,

  // apiLink: testApiLink,
  // firebaseConfig: testFirebaseConfig,
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
  const url = new URL(common.apiLink.user_status)
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
