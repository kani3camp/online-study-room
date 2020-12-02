<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>はじめての方へ</h2>
        </v-flex>
      </v-container>
    </v-main>

    <Footer />
  </v-app>
</template>

<script lang="ts">
import NavigationDrawer from '@/components/NavigationDrawer.vue'
import ToolBar from '@/components/ToolBar.vue'
import Vue from 'vue'

export default Vue.extend({
  name: 'AboutService',
  components: {
    NavigationDrawer,
    ToolBar,
  },
})
</script>

<style scoped></style>
