
export const state = () => ({
  user_id: null,
  user_name: null,
  mail_address: null,
  sum_study_time: null,
  registration_date: null,
  status_message: null,
  id_token: null,
  provider_id: null,
  last_entered: null,
})


export const mutations = {
  setUserId(state, user_id) {
    state.user_id = user_id
  },
  setUserName(state, user_name) {
    state.user_name = user_name
  },
  setMailAddress(state, mail_address) {
    state.mail_address = mail_address
  },
  setSumStudyTime(state, sum_study_time) {
    state.sum_study_time = sum_study_time
  },
  setRegistrationDate(state, registration_date) {
    state.registration_date = registration_date
  },
  setStatusMessage(state, status_message) {
    state.status_message = status_message
  },
  setIdToken(state, id_token) {
    state.id_token = id_token
  },
  setProviderId(state, provider_id) {
    state.provider_id = provider_id
  },
  setLastEntered(state, last_entered) {
    state.last_entered = last_entered
  }
}

export const actions = {

}
