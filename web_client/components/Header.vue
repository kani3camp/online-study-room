<template>
    <div>
      メニュー
      <ul>
        <li v-for="item in pages" :key="item.name">
          {{ item.name }}
        </li>
      </ul>
    </div>
</template>

<script>
    export default {
        data() {
          return {pages: [
            {name: 'TOP', to: ''},
            {name: 'ROOMS', to: ''},
            {name: 'ONLINE USERS', to: ''},
            {name: 'OPTION', to: ''}
          ]}
        },
        name: "menu"
    }
</script>

<style scoped>

</style>
