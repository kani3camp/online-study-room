<template>
  <div id="roomInfo">
    <div id="room-name"><h2>{{ room_name }}</h2></div>
    <div>
      <div id="studying-users-title">
        <h3>勉強中のユーザー</h3>
      </div>
      <div id="users">
        <div v-for="user in users" :key="user.userId" class="user">
          <p><i class="mdi mdi-account mdi-48px" /></p>
          {{ user.userName }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'RoomInfo',
  props: {
    // roomId: { type: String, required: true }
  },
  data: () => ({
    room_name: '',
    userNames: [],
    users: []
  }),
  created () {
    const vm = this
    const refreshRoomInfoInterval = 3 * 1000
    setInterval(() => {
      vm.fetchRoomInfo()
    }, refreshRoomInfoInterval)
  },
  methods: {
    async fetchRoomInfo () {
      const roomId = this.$store.state.roomId
      const url = new URL('https://us-central1-online-study-room-f1f30.cloudfunctions.net/RoomStatus')
      url.search = new URLSearchParams({ room_id: roomId }).toString()
      const resp = await fetch(url.toString(), { method: 'GET' }).then(response =>
        response.json()
      )
      if (resp.result === 'ok') {
        this.room_name = resp.room_status.room_body.name
        const userIds = resp.room_status.room_body.users
        const userNames = resp.user_names
        if (userIds && userNames) {
          if (userIds.length === userNames.length) {
            this.users = []
            for (let i = 0; i < userIds.length; i++) {
              this.users.push({
                userId: userIds[i],
                userName: userNames[i]
              })
            }
          }
        }
      } else {
        console.log(resp.message)
      }
    }
  }
}
</script>

<style scoped>
#roomInfo {
  background-color: lightgoldenrodyellow;
  width: 1920px;
  height: 818px;
  overflow: auto;
}

#room-name {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-color: black;
  border-style: solid;
  border-radius: 1rem;
}

#studying-users-title {
  text-align: left;
  margin-left: 1rem;
  font-size: larger;
}

#users {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
}

.user {
  /*background-color: azure;*/
  height: 6rem;
  width: 12rem;
  font-size: x-large;
  margin: 1rem;
  overflow: hidden;
}

</style>
