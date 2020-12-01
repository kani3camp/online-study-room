// モジュールモードで。
// index.jsに書いたものはグローバル名前空間に登録

export const state = () => ({
  isSignedIn: false,
  room_id: null,
  drawer: false,
})

export const mutations = {
  setSignInState(state, isSignedIn) {
    state.isSignedIn = isSignedIn
  },
  setRoomId(state, room_id) {
    state.room_id = room_id
  },
  setDrawer(state, newDrawer) {
    state.drawer = newDrawer
  },
  signOut(state) {
    state.isSignedIn = false
    state.user.total_study_time = null
    state.user.registration_date = null
    state.user.status_message = null
    state.user.last_entered = null
  },
}

export const actions = {}
