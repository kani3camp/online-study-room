import firebase from 'firebase/app'
import 'firebase/analytics'
require('firebase/auth')
import common from '@/plugins/common'

const firebaseConfig = common.firebaseConfig
firebase.initializeApp(firebaseConfig)
export default firebase
// export default !firebase.apps.length ? firebase.initializeApp(firebaseConfig) : firebase.app();
