// this entry is built and wrapped with a factory function
// used to generate a fresh copy of Vue for every Weex instance.

import Vue from 'vue/src/platforms/weex/runtime'

exports.Vue = Vue
