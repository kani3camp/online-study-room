<template>
  <div id="message">
    現在、全 {{ numRooms }}ルームで {{ numOnlineUsers }}人が作業中
  </div>
</template>

<script>
export default {
  name: 'Message',
  data: () => ({
    numOnlineUsers: 0,
    numRooms: 0,
    timeout: null
  }),
  async created () {
    await this.updateNums()

    const vm = this
    this.timeout = setInterval(vm.updateNums, 10 * 1000)
  },
  methods: {
    updateNums () {
      this.updateNumRooms()
      this.updateNumOnlineUsers()
    },
    async updateNumRooms () {
      const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms')
      const resp = await fetch(url.toString(), { method: 'GET' }).then(r => r.json())
      if (resp.result === 'ok') {
        if (resp.rooms) {
          this.numRooms = Number(resp.rooms.length)
        } else {
          this.numRooms = 0
        }
      } else {
        console.log(resp.message)
      }
    },
    async updateNumOnlineUsers () {
      const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/online_users')
      const resp = await fetch(url.toString(), { method: 'GET' }).then(r => r.json())
      if (resp.result === 'ok') {
        if (resp.online_users) {
          this.numOnlineUsers = Number(resp.online_users.length)
        } else {
          this.numOnlineUsers = 0
        }
      } else {
        console.log(resp.message)
      }
    }
  }
}
</script>

<style scoped>
#message {
  font-size: 2rem;
  background-color: #d9f5d9;
  width: 1920px;
  height: 100px;
}
</style>
