export const state = () => ({
  total_study_time: null,
  registration_date: null,
  status_message: null,
  last_entered: null,
})

export const mutations = {
  // setUserId(state, user_id) {
  //   state.user_id = user_id
  // },
  // setDisplayName(state, display_name) {
  //   console.log('表示名変更 => ' + display_name)
  //   state.display_name = display_name
  // },
  // setMailAddress(state, mail_address) {
  //   state.mail_address = mail_address
  // },
  setTotalStudyTime(state, total_study_time) {
    state.total_study_time = total_study_time
  },
  setRegistrationDate(state, registration_date) {
    state.registration_date = registration_date
  },
  setStatusMessage(state, status_message) {
    state.status_message = status_message
  },
  // setIdToken(state, id_token) {
  //   state.id_token = id_token
  // },
  // setProviderId(state, provider_id) {
  //   state.provider_id = provider_id
  // },
  setLastEntered(state, last_entered) {
    state.last_entered = last_entered
  },
}

export const actions = {}
