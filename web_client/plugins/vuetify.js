// src/plugins/vuetify.js

import '@mdi/font/css/materialdesignicons.css'
import Vue from 'vue'
import Vuetify from 'vuetify/lib'

Vue.use(Vuetify, )

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: '#64f3ef',
        secondary: '#b0bec5',
        accent: '#8c9eff',
        error: '#b71c1c',
        anchor: '#5be0d4',
      },
    },
  },
  icons: {
    iconfont: 'mdiSvg', // 'mdi' || 'mdiSvg' || 'md' || 'fa' || 'fa4' || 'faSvg'
  },
})
