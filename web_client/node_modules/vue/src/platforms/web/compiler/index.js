/* @flow */

import { baseOptions } from 'vue/src/platforms/web/compiler/options'
import { createCompiler } from 'compiler/index'

const { compile, compileToFunctions } = createCompiler(baseOptions)

export { compile, compileToFunctions }
