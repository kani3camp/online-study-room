/* @flow */

import on from 'vue/src/compiler/directives/on'
import bind from 'vue/src/compiler/directives/bind'
import { noop } from 'shared/util'

export default {
  on,
  bind,
  cloak: noop
}
