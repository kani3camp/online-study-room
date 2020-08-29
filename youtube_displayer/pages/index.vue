<template>
  <div id="app">
    <div class="block">
      <div id="header">
        <Title msg="Welcome to " />
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
  }
}

</script>

<style>
html {
  padding: 0;
  margin: 0;
}

body {
  margin: 0;
  padding:0;
}

#header {
  width: 100vw;
  height: 18vh;
  display: flex;
}

#main {
  width: 100vw;
  height: 10vh;
  display: flex;
}
.block {
  width: 100vw;
  height: 100vh;
  background-color: pink;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

</style>
