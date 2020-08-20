/* @flow */

process.env.VUE_ENV = 'server'

import { extend } from 'shared/util'
import modules from 'vue/src/platforms/web/server/modules'
import baseDirectives from 'vue/src/platforms/web/server/directives'
import { isUnaryTag, canBeLeftOpenTag } from 'vue/src/platforms/web/compiler/util'

import { createRenderer as _createRenderer } from 'server/create-renderer'
import { createBundleRendererCreator } from 'server/bundle-renderer/create-bundle-renderer'

export function createRenderer (options?: Object = {}): {
  renderToString: Function,
  renderToStream: Function
} {
  return _createRenderer(extend(extend({}, options), {
    isUnaryTag,
    canBeLeftOpenTag,
    modules,
    // user can provide server-side implementations for custom directives
    // when creating the renderer.
    directives: extend(baseDirectives, options.directives)
  }))
}

export const createBundleRenderer = createBundleRendererCreator(createRenderer)
