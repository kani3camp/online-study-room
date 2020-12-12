<template>
  <div id="title">
    <!--    <h1>{{ title }}</h1>-->
    <div id="room-name">
      <h2>{{ room_name }}</h2>
    </div>
    <h2 id="room-category">
      ルーム
    </h2>
  </div>
</template>

<script>

export default {
  name: 'Title',
  data: () => ({
    title: 'オンライン作業部屋'
  }),
  computed: {
    room_name () {
      return this.$store.state.roomName
    }
  },
}
</script>

<style scoped>
#title {
  padding-top: 0.5rem;
  /*background-color: #e2fcfc;*/
  width: 1536px;
  font-size: 1.3rem;
  color: #36479f;
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

</style>
