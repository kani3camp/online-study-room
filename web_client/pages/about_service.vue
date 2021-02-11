<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>はじめての方へ</h2>
        </v-flex>
        <v-flex>
          <p>
            作業がんばってください。
          </p>
        </v-flex>
      </v-container>
    </v-main>

    <Footer />
  </v-app>
</template>

<script>
import NavigationDrawer from '@/components/NavigationDrawer'
import ToolBar from '@/components/ToolBar'

export default {
  name: 'AboutService',
  components: {
    NavigationDrawer,
    ToolBar,
  },
}
</script>

<style scoped></style>
