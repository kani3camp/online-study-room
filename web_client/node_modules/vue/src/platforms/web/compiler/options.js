/* @flow */

import {
  isPreTag,
  mustUseProp,
  isReservedTag,
  getTagNamespace
} from 'vue/src/platforms/web/util'

import modules from 'vue/src/platforms/web/compiler/modules'
import directives from 'vue/src/platforms/web/compiler/directives'
import { genStaticKeys } from 'shared/util'
import { isUnaryTag, canBeLeftOpenTag } from 'vue/src/platforms/web/compiler/util'

export const baseOptions: CompilerOptions = {
  expectHTML: true,
  modules,
  directives,
  isPreTag,
  isUnaryTag,
  mustUseProp,
  canBeLeftOpenTag,
  isReservedTag,
  getTagNamespace,
  staticKeys: genStaticKeys(modules)
}
