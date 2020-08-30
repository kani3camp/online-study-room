<template>
  <div id="message">
    {{ message }}
  </div>
</template>

<script>
export default {
  name: 'Message',
  props: {
    message: { type: String, required: true }
  }
}
</script>

<style scoped>
#message {
  background-color: lightgreen;
  width: 1920px;
  height: 100px;
}
</style>
