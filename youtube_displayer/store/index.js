export const state = {
  roomId: '',
  roomName: '',
}

export const mutations = {
  setRoomId (state, newRoomId) {
    state.roomId = newRoomId
  },
  setRoomName (state, newRoomName) {
    state.roomName = newRoomName
  }
}
