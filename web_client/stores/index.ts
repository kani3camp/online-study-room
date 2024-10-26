import { defineStore } from 'pinia'

const userStore = useUserStore()

interface State {
  isSignedIn: boolean
  room_id: string
  room_name: string
  seat_id: number
  drawer: boolean
}

export const useMainStore = defineStore('main', {
  state: (): State => ({
    isSignedIn: false,
    room_id: '',
    room_name: '',
    seat_id: 0,
    drawer: false,
  }),
  actions: {
    setSignInState(isSignedIn: boolean) {
      this.isSignedIn = isSignedIn
    },
    setRoomId(room_id: string) {
      this.room_id = room_id
    },
    setRoomName(room_name: string) {
      this.room_name = room_name
    },
    setSeatId(seat_id: number) {
      this.seat_id = seat_id
    },
    setDrawer(new_drawer: boolean) {
      this.drawer = new_drawer
    },
    signOut() {
      this.isSignedIn = false
      userStore.setTotalStudyTime(null)
      userStore.setRegistrationDate(null)
      userStore.setStatusMessage(null)
      userStore.setLastEntered(null)
    },
  },
})
