<template>
  <v-app>
    <NavigationDrawer></NavigationDrawer>

    <ToolBar></ToolBar>

    <v-main>
      <v-container>
        <v-flex>
          <h2>サインイン</h2>
        </v-flex>
      </v-container>

      <v-container >
        <p>　登録・サインインともに下のボタンから行えます。</p>
        <v-row justify="center" class="google-sign-in-button">
          <img v-bind:src="imageSource"
               v-on:mouseover="changeImageToHovered"
               v-on:mousedown="changeImageToPressed"
               v-on:mouseout="changeImageToNormal"
               alt="google sign in button">
        </v-row>
      </v-container>

    </v-main>

    <Footer></Footer>
  </v-app>
</template>

<script>
import normalImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_normal_web@2x.png'
import hoveredImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_focus_web@2x.png'
import pressedImage from '~/assets/google_signin_buttons/web/2x/btn_google_signin_light_pressed_web@2x.png'

export default {
  name: "sign_in",
  data: () => ({
    imageSource: '',
  }),
  created() {
    this.imageSource = normalImage
  },
  methods: {
    changeImageToHovered() {
      console.log('fda')
      this.imageSource = hoveredImage
    },
    changeImageToPressed() {
      this.imageSource = pressedImage
    },
    changeImageToNormal() {
      this.imageSource = normalImage
    }
  },
}
</script>

<style scoped>

.google-sign-in-button img,  :hover, :active {
  max-height: 3rem;
}

</style>
