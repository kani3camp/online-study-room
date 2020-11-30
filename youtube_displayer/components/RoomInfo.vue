<template>
  <div id="roomInfo">
    <transition name="fade">
      <div v-show="show">
        <div id="room-name">
          <h2>{{ room_name }}</h2>
        </div>
        <h2 id="room-category">
          ルーム
        </h2>
        <div>
          <div id="studying-users-title">
            <p>作業中のユーザー</p>
          </div>
          <div id="users">
            <div v-for="user in users" :key="user.userId" class="user">
              <p><i class="mdi mdi-account" /></p>
              <p class="user-name">
                {{ user.userName }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  name: 'RoomInfo',
  data: () => ({
    roomIdList: [],
    timeout1: null,
    timeout2: null,
    room_name: '　　',
    users: [],
    show: true,
    switchRoomInterval: 12 * 1000,
    // refreshRoomInfoInterval: 3 * 1000
  }),
  created () {
    const vm = this

    vm.switchRoom()

    // vm.timeout1 = setTimeout(() => {
    //   vm.switchRoom()
    // }, vm.switchRoomInterval)

    // vm.timeout2 = setTimeout(() => {
    //   vm.fetchRoomInfo()
    // }, vm.refreshRoomInfoInterval)
  },
  destroyed () {
    clearTimeout(this.timeout1)
    // clearTimeout(this.timeout2)
  },
  methods: {
    async switchRoom () {
      const vm = this

      vm.show = false

      // 全てのroom_idのリストを更新
      await vm.retrieveRoomIdList()

      // room_idを次のものに進める。ない場合はリストの先頭から
      const currentIndex = vm.roomIdList.indexOf(vm.$store.state.roomId)
      if (currentIndex === -1) {
        await vm.$store.commit('setRoomId', vm.roomIdList[0])
      } else {
        const nextIndex = (currentIndex + 1) % this.roomIdList.length
        await vm.$store.commit('setRoomId', vm.roomIdList[nextIndex])
      }
      await vm.fetchRoomInfo()

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
        const list = []
        resp.rooms.forEach((room) => {
          list.push(room.room_id)
        })
        vm.roomIdList = list
      } else {
        console.log(resp.message)
      }
    },
    async fetchRoomInfo () {
      const vm = this
      const roomId = vm.$store.state.roomId
      const url = new URL('https://io551valj4.execute-api.ap-northeast-1.amazonaws.com/room_status')
      url.search = new URLSearchParams({ room_id: roomId }).toString()
      const resp = await fetch(url.toString(), { method: 'GET' }).then(response =>
        response.json()
      )
      if (resp.result === 'ok') {
        this.room_name = resp.room_status.room_body.name
        const userIds = resp.room_status.room_body.users
        const users = resp.users
        if (userIds && users) {
          if (userIds.length === users.length) {
            const list = []
            for (let i = 0; i < userIds.length; i++) {
              list.push({
                userId: userIds[i],
                userName: users[i].display_name,
              })
            }
            this.users = list
          } else {
            console.log('userIds.length !== users.length')
          }
        } else {
          this.users = []
        }
      } else {
        console.log(resp.message)
      }

      // this.timeout2 = setTimeout(() => {
      //   vm.fetchRoomInfo()
      // }, vm.refreshRoomInfoInterval)
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
.fade-enter-active {
  transition: .3s;
}
.fade-leave-active {
  transition: 1.5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

#room-name {
  display: inline-block;
  padding: .3rem .6rem;
  border: solid 0.2rem #36479f;
  border-radius: 1rem;
}
#room-category {
  display: inline-block;
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
