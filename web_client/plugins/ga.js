import firebase from '~/plugins/firebase'

export default ({ app }) => {
  // if (process.env.NODE_ENV !== 'production') return

  firebase.analytics()
}
