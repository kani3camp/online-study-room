import firebase from '~/plugins/firebase'
import { useMainStore } from '#imports'
import { useUserStore } from '#imports'

const mainStore = useMainStore()
const userStore = useUserStore()

const productionApiLink = {
  rooms: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms',
  room_status: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_status',
  online_users: '',
  user_status: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/user_status',
  change_user_info: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/change_user_info',
  send_contact_form:
    'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/send_contact_form',
  create_new_room: '',
  news: 'https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/news',
  room_layout: '',
  websocket: 'wss://0ieer51ju9.execute-api.ap-northeast-1.amazonaws.com/release',
}

const testApiLink = {
  rooms: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_rooms',
  room_status: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_room_status',
  online_users: '',
  user_status: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_user_status',
  change_user_info:
    'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_change_user_info',
  send_contact_form:
    'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_send_contact_form',
  create_new_room: '',
  news: 'https://fquvf774q5.execute-api.ap-northeast-1.amazonaws.com/test_news',
  room_layout: '',
  websocket: 'wss://jukb7tixp7.execute-api.ap-northeast-1.amazonaws.com/test',
}

const common = {
  key: {
    youtubeLink: 'https://www.youtube.com/channel/UCXuD2XmPTdpVy7zmwbFVZWg',
    twitterLink: 'https://twitter.com/osr_soraride',
  },
  // TODO: デプロイ時に変更
  // apiLink: productionApiLink,

  apiLink: testApiLink,

  onAuthStateChanged: () => {
    firebase.auth().onAuthStateChanged(async (user) => {
      if (user) {
        console.log('User is signed in.')
        mainStore.setSignInState(true)

        await common.getUserData()

        await firebase.auth().currentUser!.getIdToken(true) // refresh idToken
      } else {
        console.log('User is signed out.')
        mainStore.signOut()
      }
    })
  },
  getUserData: async () => {
    const params = {
      user_id: firebase.auth().currentUser!.uid,
    }
    const user_data = await common.httpGet(common.apiLink.user_status, params)
    if (user_data.result !== 'ok') {
      console.log(user_data)
    } else {
      const user_body = user_data['user_status']['user_body']
      userStore.status_message = user_body.status
      userStore.total_study_time = user_body.total_study_time
      userStore.registration_date = new Date(user_body.registration_date)
      userStore.last_entered = new Date(user_body.last_entered)
    }
  },
  httpGet: async (url_str: string, params: any) => {
    const url = new URL(url_str)
    url.search = new URLSearchParams(params).toString()
    const response = await fetch(url.toString(), { method: 'GET' })
    return await response.json()
  },
  httpPost: async (url_str: string, _params: any) => {
    const response = await fetch(url_str, {
      method: 'POST',
      body: JSON.stringify(_params),
    })
    return await response.json()
  },
}

export default common
