<template>
  <div id="title">
    <h1>{{ title }}</h1>
  </div>
</template>

<script>

export default {
  name: 'Title',
  data: () => ({
    title: 'オンライン勉強部屋'
  })
}
</script>

<style scoped>
#title {
  padding-top: 1rem;
  background-color: #ecc6a1;
  width: 1536px;
}

</style>
