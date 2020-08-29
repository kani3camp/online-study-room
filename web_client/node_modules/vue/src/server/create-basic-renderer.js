/* @flow */

import { createWriteFunction } from 'vue/src/server/write'
import { createRenderFunction } from 'vue/src/server/render'
import type { RenderOptions } from 'vue/src/server/create-renderer'

export function createBasicRenderer ({
  modules = [],
  directives = {},
  isUnaryTag = (() => false),
  cache
}: RenderOptions = {}) {
  const render = createRenderFunction(modules, directives, isUnaryTag, cache)

  return function renderToString (
    component: Component,
    context: any,
    done: any
  ): void {
    if (typeof context === 'function') {
      done = context
      context = {}
    }
    let result = ''
    const write = createWriteFunction(text => {
      result += text
      return false
    }, done)
    try {
      render(component, write, context, () => {
        done(null, result)
      })
    } catch (e) {
      done(e)
    }
  }
}
