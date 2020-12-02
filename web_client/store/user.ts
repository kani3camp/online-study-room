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
export default class User extends VuexModule {
  public info: UserData = {
    isSignedIn: false,
    roomId: '',
    drawer: false,
    totalStudyTime: 0,
    registrationDate: new Date(),
    statusMessage: '',
    lastEntered: new Date(),
  }

  @Mutation
  public setSignInState(isSignedIn: boolean): void {
    this.info.isSignedIn = isSignedIn
  }

  @Mutation
  public setRoomId(roomId: string) {
    this.info.roomId = roomId
  }

  @Mutation
  public setDrawer(newDrawer: boolean) {
    this.info.drawer = newDrawer
  }

  @Action
  public signOut() {
    this.info.isSignedIn = false
    this.info.totalStudyTime = 0
    this.info.registrationDate = new Date()
    this.info.statusMessage = ''
    this.info.lastEntered = new Date()
  }

  @Mutation
  public setTotalStudyTime(totalStudyTime: number) {
    this.info.totalStudyTime = totalStudyTime
  }

  @Mutation
  public setRegistrationDate(registrationDate: Date) {
    this.info.registrationDate = registrationDate
  }

  @Mutation
  public setStatusMessage(statusMessage: string) {
    this.info.statusMessage = statusMessage
  }

  @Mutation
  public setLastEntered(lastEntered: Date) {
    this.info.lastEntered = lastEntered
  }
}
