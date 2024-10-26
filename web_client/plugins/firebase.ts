import firebase from 'firebase/app'
import 'firebase/analytics'
require('firebase/auth')

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
const testFirebaseConfig = {
  apiKey: 'AIzaSyAE5Hm9U0pQDxXAPNkmthO_z--CFDiSn6w',
  authDomain: 'test-online-study-space.firebaseapp.com',
  projectId: 'test-online-study-space',
  storageBucket: 'test-online-study-space.appspot.com',
  messagingSenderId: '788678710894',
  appId: '1:788678710894:web:2cb434fbd57cdc9bfd787f',
  measurementId: 'G-450CVZ6FWP',
}

const firebaseConfig = testFirebaseConfig // TODO: デプロイ時に変更
firebase.initializeApp(firebaseConfig)
export default firebase
// export default !firebase.apps.length ? firebase.initializeApp(firebaseConfig) : firebase.app();
