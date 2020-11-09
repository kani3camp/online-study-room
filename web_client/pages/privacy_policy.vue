<template>
  <v-app>
    <NavigationDrawer />

    <ToolBar />

    <v-main>
      <v-container>
        <v-flex>
          <h2>プライバシーポリシー</h2>
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
  name: 'PrivacyPolicy',
  components: {
    NavigationDrawer,
    ToolBar,
  },
}
</script>

<style scoped></style>
