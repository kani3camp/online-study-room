/* @flow */

import { preTransformRecycleList } from 'vue/src/platforms/weex/compiler/modules/recycle-list/recycle-list'
import { postTransformComponent } from 'vue/src/platforms/weex/compiler/modules/recycle-list/component'
import { postTransformComponentRoot } from 'vue/src/platforms/weex/compiler/modules/recycle-list/component-root'
import { postTransformText } from 'vue/src/platforms/weex/compiler/modules/recycle-list/text'
import { preTransformVBind } from 'vue/src/platforms/weex/compiler/modules/recycle-list/v-bind'
import { preTransformVIf } from 'vue/src/platforms/weex/compiler/modules/recycle-list/v-if'
import { preTransformVFor } from 'vue/src/platforms/weex/compiler/modules/recycle-list/v-for'
import { postTransformVOn } from 'vue/src/platforms/weex/compiler/modules/recycle-list/v-on'
import { preTransformVOnce } from 'vue/src/platforms/weex/compiler/modules/recycle-list/v-once'

let currentRecycleList = null

function shouldCompile (el: ASTElement, options: WeexCompilerOptions) {
  return options.recyclable ||
    (currentRecycleList && el !== currentRecycleList)
}

function preTransformNode (el: ASTElement, options: WeexCompilerOptions) {
  if (el.tag === 'recycle-list') {
    preTransformRecycleList(el, options)
    currentRecycleList = el
  }
  if (shouldCompile(el, options)) {
    preTransformVBind(el)
    preTransformVIf(el, options) // also v-else-if and v-else
    preTransformVFor(el, options)
    preTransformVOnce(el)
  }
}

function transformNode (el: ASTElement, options: WeexCompilerOptions) {
  if (shouldCompile(el, options)) {
    // do nothing yet
  }
}

function postTransformNode (el: ASTElement, options: WeexCompilerOptions) {
  if (shouldCompile(el, options)) {
    // mark child component in parent template
    postTransformComponent(el, options)
    // mark root in child component template
    postTransformComponentRoot(el)
    // <text>: transform children text into value attr
    if (el.tag === 'text') {
      postTransformText(el)
    }
    postTransformVOn(el)
  }
  if (el === currentRecycleList) {
    currentRecycleList = null
  }
}

export default {
  preTransformNode,
  transformNode,
  postTransformNode
}
