<template>
  <v-app>
    <NavigationDrawer></NavigationDrawer>

    <ToolBar></ToolBar>

    <v-main>
      <v-container>
        <v-flex>
          <h2>はじめての方へ</h2>
        </v-flex>
      </v-container>

    </v-main>

    <Footer></Footer>
  </v-app>
</template>

<script>
import NavigationDrawer from "@/components/NavigationDrawer"
import ToolBar from "@/components/ToolBar"


export default {
  name: "about_service",
  components: {
    NavigationDrawer,
    ToolBar,
  }
}
</script>

<style scoped>

</style>
