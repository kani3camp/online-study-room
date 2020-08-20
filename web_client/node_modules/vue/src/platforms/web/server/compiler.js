/* @flow */

import { baseOptions } from 'vue/src/platforms/web/compiler/options'
import { createCompiler } from 'server/optimizing-compiler/index'

const { compile, compileToFunctions } = createCompiler(baseOptions)

export {
  compile as ssrCompile,
  compileToFunctions as ssrCompileToFunctions
}
