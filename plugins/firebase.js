import firebase from "firebase/app"
require('firebase/auth')

const firebaseConfig = {
  apiKey: "AIzaSyBsMmNzh3qo1X-UpS1T9Q6C07ODlohM-_4",
  authDomain: "online-study-room-f1f30.firebaseapp.com",
  databaseURL: "https://online-study-room-f1f30.firebaseio.com",
  projectId: "online-study-room-f1f30",
  storageBucket: "online-study-room-f1f30.appspot.com",
  messagingSenderId: "77400576490",
  appId: "1:77400576490:web:b10794e40775edf921e293",
  measurementId: "G-4DN92GENJJ"
}
firebase.initializeApp(firebaseConfig)
export default firebase
// export default !firebase.apps.length ? firebase.initializeApp(firebaseConfig) : firebase.app();
