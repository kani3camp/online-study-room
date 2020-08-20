/* @flow */

import modules from 'vue/src/platforms/web/server/modules'
import directives from 'vue/src/platforms/web/server/directives'
import { isUnaryTag, canBeLeftOpenTag } from 'vue/src/platforms/web/compiler/util'
import { createBasicRenderer } from 'server/create-basic-renderer'

export default createBasicRenderer({
  modules,
  directives,
  isUnaryTag,
  canBeLeftOpenTag
})
