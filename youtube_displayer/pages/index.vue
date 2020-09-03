<template>
  <div id="app">
    <div class="block">
      <div id="header">
        <Title />
        <Time />
      </div>

      <div id="main">
        <RoomInfo />
        <!--        <Log></Log>-->
      </div>

      <div id="footer">
        <Message />
      </div>
    </div>
  </div>
</template>

<script>
import Title from '~/components/Title'
import Time from '~/components/Time'
import RoomInfo from '~/components/RoomInfo'
import Message from '~/components/Message'

export default {
  components: {
    Title,
    Time,
    RoomInfo,
    Message
  },
  data: () => ({
  }),
  created () {
  },
  destroyed () {
  },
  methods: {
  }
}

</script>

<style>
html {
  font-size: 2rem;
  padding: 0;
  margin: 0;
  width: 1920px;
  height: 1080px;
  overflow: auto;
  /*cursor: none; !* todo *!*/
}

body {
  margin: 0;
  padding:0;
}

#header {
  width: 1920px;
  height: 160px;
  display: flex;
}

#main {
  width: 1920px;
  display: flex;
}

#app {
  font-family: 'Hannari', Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#footer {
  width: 1920px;
  /*height: 100px;*/
}

.block {
  width: 1920px;
  /*background-color: pink;*/
}

</style>
