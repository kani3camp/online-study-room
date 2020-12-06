<template>
  <v-app>
    <v-app-bar
      app
      flat
    >
      <v-btn
        icon
        @click="if_show_dialog=true"
      >
        <v-icon>mdi-close</v-icon>
      </v-btn>

    </v-app-bar>

    <v-main>
      席を選ぼう。

      <RoomLayout />
    </v-main>
  </v-app>
</template>
<script>
import RoomLayout from '~/components/RoomLayout'
export default {
  name: 'EnterRoom',
  components: {
    RoomLayout,
  },
}
</script>

<style scoped></style>
