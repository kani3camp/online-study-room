// モジュールモードで。
// index.jsに書いたものはグローバル名前空間に登録

export const state = () => ({
  isSignedIn: false,
  room_id: '',
  room_name: '',
  seat_id: 0,
  drawer: false,
})

export const mutations = {
  setSignInState(state, isSignedIn) {
    state.isSignedIn = isSignedIn
  },
  setRoomId(state, room_id) {
    state.room_id = room_id
  },
  setRoomName(state, room_name) {
    state.room_name = room_name
  },
  setSeatId(state, seat_id) {
    state.seat_id = seat_id
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
