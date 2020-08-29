/* @flow */

import { toNumber, toString, looseEqual, looseIndexOf } from 'shared/util'
import { createTextVNode, createEmptyVNode } from 'core/vdom/vnode'
import { renderList } from 'vue/src/core/instance/render-helpers/render-list'
import { renderSlot } from 'vue/src/core/instance/render-helpers/render-slot'
import { resolveFilter } from 'vue/src/core/instance/render-helpers/resolve-filter'
import { checkKeyCodes } from 'vue/src/core/instance/render-helpers/check-keycodes'
import { bindObjectProps } from 'vue/src/core/instance/render-helpers/bind-object-props'
import { renderStatic, markOnce } from 'vue/src/core/instance/render-helpers/render-static'
import { bindObjectListeners } from 'vue/src/core/instance/render-helpers/bind-object-listeners'
import { resolveScopedSlots } from 'vue/src/core/instance/render-helpers/resolve-scoped-slots'
import { bindDynamicKeys, prependModifier } from 'vue/src/core/instance/render-helpers/bind-dynamic-keys'

export function installRenderHelpers (target: any) {
  target._o = markOnce
  target._n = toNumber
  target._s = toString
  target._l = renderList
  target._t = renderSlot
  target._q = looseEqual
  target._i = looseIndexOf
  target._m = renderStatic
  target._f = resolveFilter
  target._k = checkKeyCodes
  target._b = bindObjectProps
  target._v = createTextVNode
  target._e = createEmptyVNode
  target._u = resolveScopedSlots
  target._g = bindObjectListeners
  target._d = bindDynamicKeys
  target._p = prependModifier
}
