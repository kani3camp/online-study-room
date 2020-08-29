import { initMixin } from 'vue/src/core/instance/init'
import { stateMixin } from 'vue/src/core/instance/state'
import { renderMixin } from 'vue/src/core/instance/render'
import { eventsMixin } from 'vue/src/core/instance/events'
import { lifecycleMixin } from 'vue/src/core/instance/lifecycle'
import { warn } from 'vue/src/core/util'

function Vue (options) {
  if (process.env.NODE_ENV !== 'production' &&
    !(this instanceof Vue)
  ) {
    warn('Vue is a constructor and should be called with the `new` keyword')
  }
  this._init(options)
}

initMixin(Vue)
stateMixin(Vue)
eventsMixin(Vue)
lifecycleMixin(Vue)
renderMixin(Vue)

export default Vue
