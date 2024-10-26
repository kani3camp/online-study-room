import type { Pinia } from 'pinia'
import { useMainStore } from '~/stores'

export default defineNuxtPlugin(({ $pinia }) => {
  return {
    provide: {
      store: useMainStore($pinia as Pinia),
    },
  }
})
