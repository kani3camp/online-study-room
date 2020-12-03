import { Module, VuexModule, Mutation, Action } from 'vuex-module-decorators'

type UserData = {
  isSignedIn: boolean
  roomId: string
  drawer: boolean
  totalStudyTime: number
  registrationDate: Date
  statusMessage: string
  lastEntered: Date
}

@Module({
  name: 'user',
  stateFactory: true,
  namespaced: true,
})
class User extends VuexModule implements UserData {
  isSignedIn = false
  roomId = ''
  drawer = false
  totalStudyTime = 0
  registrationDate = new Date()
  statusMessage = ''
  lastEntered = new Date()

  @Mutation
  public setSignInState(isSignedIn: boolean): void {
    this.isSignedIn = isSignedIn
  }

  @Mutation
  public setRoomId(roomId: string) {
    this.roomId = roomId
  }

  @Mutation
  public setDrawer(newDrawer: boolean) {
    this.drawer = newDrawer
  }

  @Action
  public signOut() {
    this.isSignedIn = false
    this.totalStudyTime = 0
    this.registrationDate = new Date()
    this.statusMessage = ''
    this.lastEntered = new Date()
  }

  @Mutation
  public setTotalStudyTime(totalStudyTime: number) {
    this.totalStudyTime = totalStudyTime
  }

  @Mutation
  public setRegistrationDate(registrationDate: Date) {
    this.registrationDate = registrationDate
  }

  @Mutation
  public setStatusMessage(statusMessage: string) {
    this.statusMessage = statusMessage
  }

  @Mutation
  public setLastEntered(lastEntered: Date) {
    this.lastEntered = lastEntered
  }
}

export default User
