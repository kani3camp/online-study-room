declare interface RoomsResponse {
  result: string
  message: string
  rooms: {
    roomId: string
    roomBody: {
      created: Date
      name: string
      users: string[]
      type: string
      themeColorHex: string
    }
  }
}
