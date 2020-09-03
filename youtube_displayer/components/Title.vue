<template>
  <div id="title">
    <h1>{{ title }}</h1>
  </div>
</template>

<script>

export default {
  name: 'Title',
  data: () => ({
    title: 'オンライン作業部屋'
  })
}
</script>

<style scoped>
#title {
  padding-top: 0.5rem;
  background-color: #e2fcfc;
  width: 1536px;
  font-size: 1.3rem;
}

</style>
