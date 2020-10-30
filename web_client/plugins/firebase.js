import firebase from "firebase/app"
require('firebase/auth')

const firebaseConfig = {
  apiKey: "AIzaSyCIY7RCArApaRpg4nkVUJCze88OmF6-wGM",
  authDomain: "online-study-space.firebaseapp.com",
  databaseURL: "https://online-study-space.firebaseio.com",
  projectId: "online-study-space",
  storageBucket: "online-study-space.appspot.com",
  messagingSenderId: "486366182751",
  appId: "1:486366182751:web:dc94794117a6268a050aea",
  measurementId: "G-KH6EFLWF7T"
}
firebase.initializeApp(firebaseConfig)
export default firebase
// export default !firebase.apps.length ? firebase.initializeApp(firebaseConfig) : firebase.app();
