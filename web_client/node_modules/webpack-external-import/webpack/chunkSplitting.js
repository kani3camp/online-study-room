"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.interleaveConfig = interleaveConfig;
exports.interleaveStyleConfig = interleaveStyleConfig;
exports.interleaveStyleJsConfig = interleaveStyleJsConfig;
exports.hasExternalizedModuleViaJson = void 0;

var _mem = _interopRequireDefault(require("mem"));

var _utils = require("./utils");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

var interleaveMap = (0, _utils.getInterleaveConfig)();
var hasExternalizedModuleViaJson = (0, _mem["default"])(function (moduleResource) {
  if (!moduleResource || !interleaveMap) return;
  var interleaveKeys = Object.keys(interleaveMap || {});

  if (interleaveKeys) {
    var foundMatch = interleaveKeys.find(function (item) {
      return moduleResource.includes(item);
    });
    return interleaveMap[foundMatch] || false;
  }
});
exports.hasExternalizedModuleViaJson = hasExternalizedModuleViaJson;

function interleaveConfig(_ref) {
  var testPath = _ref.testPath,
      manifestName = _ref.manifestName;
  return {
    test: function test(module) {
      // check if module has a resource path (not virtual modules)
      if (module.resource) {
        return module.resource.includes(testPath) && !!hasExternalizedModuleViaJson(module.resource, manifestName);
      }
    },
    name: function name(module) {
      // Check if module is listed in the interleave interface
      var foundValue = hasExternalizedModuleViaJson(module.resource, manifestName);
      if (foundValue) return foundValue;
    },
    // force module into a chunk regardless of how its used
    enforce: true // might need for next.js
    // reuseExistingChunk: false,

  };
}

function interleaveStyleConfig(_ref2) {
  var manifestName = _ref2.manifestName;
  return {
    test: function test(module) {
      // check if module has a resource path (not virtual modules)
      if (module.constructor.name === "CssModule") {
        return !!hasExternalizedModuleViaJson(module.identifier(), manifestName);
      }
    },
    // eslint-disable-next-line no-unused-vars
    name: function name(module, chunks, cacheGroupKey) {
      // Check if module is listed in the interleave interface
      if (!module.resource) {
        var foundValue = hasExternalizedModuleViaJson(module.resource || module.identifier(), manifestName);
        if (foundValue) return "".concat(foundValue, "-style");
      }

      return "styles";
    },
    // force module into a chunk regardless of how its used
    chunks: "all",
    enforce: true,
    reuseExistingChunk: false
  };
}

function interleaveStyleJsConfig(_ref3) {
  var manifestName = _ref3.manifestName;
  return {
    test: function test(module) {
      if (module.request && module.request.match(/\\(style-loader|css-loader|sass-loader)$/)) {
        return true;
      }

      if (module.constructor.name === "CssModule") {
        return false;
      }

      if (module.resource && module.resource.match(/\.(scss|css)$/)) {
        return true;
      }

      return false;
    },
    name: function name(module, chunks) {
      return "".concat(manifestName, "-stylejs");
    },
    chunks: "all",
    enforce: true,
    reuseExistingChunk: false
  };
}