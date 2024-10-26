import { defineStore } from 'pinia'

interface State {
  total_study_time: number | null
  registration_date: Date | null
  status_message: string | null
  last_entered: Date | null
}

export const useUserStore = defineStore('user', {
  state: (): State => ({
    total_study_time: null,
    registration_date: null,
    status_message: null,
    last_entered: null,
  }),
  actions: {
    setTotalStudyTime(total_study_time: number | null) {
      this.total_study_time = total_study_time
    },
    setRegistrationDate(registration_date: Date | null) {
      this.registration_date = registration_date
    },
    setStatusMessage(status_message: string | null) {
      this.status_message = status_message
    },
    setLastEntered(last_entered: Date | null) {
      this.last_entered = last_entered
    },
  },
})
