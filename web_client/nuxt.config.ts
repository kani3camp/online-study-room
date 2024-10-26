export default defineNuxtConfig({
  ssr: false,

  /*
   ** Nuxt target
   ** See https://nuxtjs.org/api/configuration-target
   */
  // target: 'static',
  /*
   ** Headers of the page
   ** See https://nuxtjs.org/api/configuration-head
   */
  typescript: {
    typeCheck: true,
  },

  css: [],

  plugins: [
    {
      src: '~/plugins/ga.ts',
      mode: 'client',
    },
  ],

  /*
   ** Auto import components
   ** See https://nuxtjs.org/api/configuration-components
   */
  components: true,

  /*
   ** Nuxt.js modules
   */
  modules: ['vuetify-nuxt-module', '@pinia/nuxt'],

  /*
   ** Build configuration
   ** See https://nuxtjs.org/api/configuration-build/
   */
  // router: {
  //   base: './'
  // },
  watchers: {
    webpack: {
      // poll: true,
    },
  },

  build: {},

  generate: {
    // fallback: true,
  },

  vuetify: {
    moduleOptions: {
      /* module specific options */
    },
    vuetifyOptions: {
      /* vuetify options */
    },
  },

  compatibilityDate: '2024-10-21',
})
