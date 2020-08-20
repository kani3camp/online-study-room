import klass from 'vue/src/platforms/weex/compiler/modules/class'
import style from 'vue/src/platforms/weex/compiler/modules/style'
import props from 'vue/src/platforms/weex/compiler/modules/props'
import append from 'vue/src/platforms/weex/compiler/modules/append'
import recycleList from 'vue/src/platforms/weex/compiler/modules/recycle-list'

export default [
  recycleList,
  klass,
  style,
  props,
  append
]
