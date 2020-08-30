<template>
  <div id="roomInfo">
    <h2>{{ room_name }}</h2>
    <div>
      <div id="users">
        <div v-for="user in userIds" :key="user" class="user">
          <p><i class="mdi mdi-account mdi-48px" /></p>
          {{ user }}
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
    userIds: [
      'fafda',
      'fajowef',
      'jifdoajfdoi',
      'jda',
      'jfaio',
      'jofaji',
      'jfiwoaj',
      'wour',
      'hvanonvaL',
      'fafda',
      'fajowef',
      'jifdoajfdoi',
      'jda',
      'jfaio',
      'jofaji',
      'jfiwoaj',
      'wour',
      'hvanonvaL'
    ]
  }),
  computed: {
    userNames () {
      let list = []
      this.userIds.forEach((user) => {
        const url = new URL('')
        list.push()
      })
      return list
    }
  },
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
        // this.userIds = resp.room_status.room_body.userIds
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
  width: 100vw;
  height: 72vh;
  overflow: auto;
}

#users {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
}

.user {
  background-color: azure;
  height: 12vh;
  width: 20vw;
  font-size: x-large;
  margin: 1rem;
  overflow: hidden;
}

</style>
