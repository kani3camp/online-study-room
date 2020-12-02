import firebase from '@/plugins/firebase.ts'

export default ({ app }: any) => {
  // if (process.env.NODE_ENV !== 'production') return

  firebase.analytics()
}
