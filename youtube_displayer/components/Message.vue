<template>
  <div id="message">
    message
    {{ message }}
  </div>
</template>

<script>
export default {
  name: 'Message',
  props: {
    message: String
  }
}
</script>

<style scoped>
#message {
  background-color: lightgreen;
  width: 100vw;
}
</style>
