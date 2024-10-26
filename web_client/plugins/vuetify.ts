import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
// import Vue from 'vue'

const app = createApp()
const vuetify = createVuetify({
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
app.use(vuetify)

// Vue.use(Vuetify)

// export default new Vuetify({
//   theme: {
//     themes: {
//       light: {
//         primary: '#64f3ef',
//         secondary: '#b0bec5',
//         accent: '#8c9eff',
//         error: '#b71c1c',
//         anchor: '#5be0d4',
//       },
//     },
//   },
//   icons: {
//     iconfont: 'mdiSvg', // 'mdi' || 'mdiSvg' || 'md' || 'fa' || 'fa4' || 'faSvg'
//   },
// })
