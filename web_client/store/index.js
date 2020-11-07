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
    state.user.user_id = null
    state.user.display_name = null
    state.user.mail_address = null
    state.user.sum_study_time = null
    state.user.registration_date = null
    state.user.status_message = null
    state.user.id_token = null
    state.user.provider_id = null
    state.user.last_entered = null
  },
}

export const actions = {}
