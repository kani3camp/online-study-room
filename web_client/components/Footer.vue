<template>
  <v-footer
    app
  >
    <span>&copy; Online Study Room {{ new Date().getFullYear() }}</span>
    <v-spacer></v-spacer>
    <nuxt-link to="/terms_of_service">利用規約</nuxt-link>
  </v-footer>

</template>

<script>
    export default {
      name: "Footer",
      methods: {
      }
    }
</script>

<style scoped>

</style>
