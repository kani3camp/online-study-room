<template>
  <div id="roomInfo">
    room info
    <h2>{{ room_name }}</h2>
    <div>
      <ul>
        <li v-for="user in users" v-bind:key="user">
          {{ user }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'RoomInfo',
  props: {
    room_id: String
  },
  data: () => ({
    room_name: '',
    users: []
  }),
  created () {
    setInterval(this.fetchRoomInfo, 5000)
  },
  methods: {
    async fetchRoomInfo () {
      // todo
    }
  }
}
</script>

<style scoped>
#roomInfo {
  background-color: lightgoldenrodyellow;
  width: 100vw;
}

</style>
