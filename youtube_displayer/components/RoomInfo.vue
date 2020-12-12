<template>
  <div id="roomInfo">
    <transition name="fade">
      <div v-show="show">
        <RoomLayout
          :room-id="room_id"
          :layout="room_layout"
        />
      </div>
    </transition>
  </div>
</template>

<script>
import RoomLayout from '@/components/RoomLayout'
export default {
  name: 'RoomInfo',
  components: {
    RoomLayout,
  },
  data: () => ({
    roomIdList: [],
    roomNameList: [],
    timeout1: null,
    room_id: '',
    room_name: '　　',
    room_layout: null,
    users: [],
    show: true,
    switchRoomInterval: 3 * 1000,
  }),
  created () {
    this.switchRoom()
  },
  destroyed () {
    clearTimeout(this.timeout1)
  },
  methods: {
    async switchRoom () {
      const vm = this

      vm.show = false

      // 全てのroom_idのリストを更新
      await vm.retrieveRoomIdList()

      // room_idを次のものに進める。ない場合はリストの先頭から
      const currentIndex = vm.roomIdList.indexOf(vm.$store.state.roomId)
      console.log('current index = ', currentIndex)
      let nextRoomName = ''
      if (currentIndex === -1) {
        await vm.$store.commit('setRoomId', vm.roomIdList[0])
        nextRoomName = vm.roomNameList[0]
      } else {
        const nextIndex = (currentIndex + 1) % this.roomIdList.length
        await vm.$store.commit('setRoomId', vm.roomIdList[nextIndex])
        nextRoomName = vm.roomNameList[nextIndex]
      }
      await vm.fetchRoomInfo(nextRoomName)

      vm.timeout1 = setTimeout(() => {
        vm.switchRoom()
      }, vm.switchRoomInterval)

      vm.show = true
    },
    async retrieveRoomIdList () {
      const vm = this
      const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/rooms')
      const resp = await fetch(url.toString(), { method: 'GET' }).then(response => response.json())
      if (resp.result === 'ok') {
        const idList = []
        const nameList = []
        resp.rooms.forEach((room) => {
          idList.push(room.room_id)
          nameList.push(room.room_body.name)
        })
        vm.roomIdList = idList
        vm.roomNameList = nameList
      } else {
        console.error(resp.message)
      }
    },
    async fetchRoomInfo (nextRoomName) {
      const vm = this
      const roomId = vm.$store.state.roomId
      const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_status')
      url.search = new URLSearchParams({ room_id: roomId }).toString()
      const resp = await fetch(url.toString(), { method: 'GET' }).then(response =>
        response.json()
      )
      this.$store.commit('setRoomName', nextRoomName)

      if (resp.result === 'ok') {
        this.room_id = roomId
        this.room_layout = resp.room_layout
        const userSeats = resp.room_status.room_body.users
        const users = resp.users
        if (userSeats && users) {
          if (userSeats.length === users.length) {
            const list = []
            for (let i = 0; i < userSeats.length; i++) {
              list.push({
                userId: userSeats[i].user_id,
                userName: users[i].display_name,
              })
            }
            this.users = list
          } else {
            console.error('userSeats.length !== users.length')
          }
        } else {
          this.users = []
        }
      } else {
        console.error(resp.message)
        this.room_layout = null
      }
    }
  }
}
</script>

<style scoped>

h2 {
  color: #36479f;
}

#roomInfo {
  padding-top: 0.5rem;
  /*background-color: #fcfcf2;*/
  width: 1920px;
  height: 818px;
  overflow: auto;
}

/*Vue.jsの書き方(<transition>に作用)*/
.fade-enter-active {
  transition: .3s;
}
.fade-leave-active {
  transition: 2s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

#studying-users-title {
  text-align: left;
  margin-left: 1rem;
  font-size: 1.5rem;
}

#users {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
}

.user {
  /*background-color: pink;*/
  height: 5rem;
  width: 10rem;
  margin: 1rem;
  overflow: hidden;
}
.user i {
  font-size: 2rem;
}
.user .user-name {
  font-size: 1.5rem;
}

</style>
