"use strict";

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

/*
	MIT License http://www.opensource.org/licenses/mit-license.php
	Author Tobias Koppers @sokra
*/
var ExternalModule = require("webpack/lib/ExternalModule");

var ExternalModuleFactoryPlugin =
/*#__PURE__*/
function () {
  function ExternalModuleFactoryPlugin(type, externals) {
    _classCallCheck(this, ExternalModuleFactoryPlugin);

    this.type = type;
    this.externals = externals;
  }

  _createClass(ExternalModuleFactoryPlugin, [{
    key: "apply",
    value: function apply(normalModuleFactory) {
      var _this = this;

      var globalType = this.type;
      normalModuleFactory.hooks.factory.tap("ExternalModuleFactoryPlugin", function (factory) {
        return function (data, callback) {
          var context = data.context;
          var dependency = data.dependencies[0];

          var handleExternal = function handleExternal(value, type, callback) {
            if (typeof type === "function") {
              callback = type;
              type = undefined;
            }

            if (value === false) return factory(data, callback);
            if (value === true) value = dependency.request;

            if (type === undefined && /^[a-z0-9]+ /.test(value)) {
              var idx = value.indexOf(" ");
              type = value.substr(0, idx);
              value = value.substr(idx + 1);
            }

            var externalModule = new ExternalModule(value, type || globalType, dependency.request);
            externalModule.id = dependency.userRequest;
            callback(null, externalModule);
            return true;
          };

          var handleExternals = function handleExternals(externals, callback) {
            if (typeof externals === "string") {
              if (externals === dependency.request) {
                return handleExternal(dependency.request, callback);
              }
            } else if (Array.isArray(externals)) {
              var i = 0;

              var next = function next() {
                var asyncFlag;

                var handleExternalsAndCallback = function handleExternalsAndCallback(err, module) {
                  if (err) return callback(err);

                  if (!module) {
                    if (asyncFlag) {
                      asyncFlag = false;
                      return;
                    }

                    return next();
                  }

                  callback(null, module);
                };

                do {
                  asyncFlag = true;
                  if (i >= externals.length) return callback();
                  handleExternals(externals[i++], handleExternalsAndCallback);
                } while (!asyncFlag);

                asyncFlag = false;
              };

              next();
              return;
            } else if (externals instanceof RegExp) {
              if (externals.test(dependency.request)) {
                return handleExternal(dependency.request, callback);
              }
            } else if (typeof externals === "function") {
              externals.call(null, context, dependency.request, function (err, value, type) {
                if (err) return callback(err);

                if (value !== undefined) {
                  handleExternal(value, type, callback);
                } else {
                  callback();
                }
              });
              return;
            } else if (_typeof(externals) === "object" && Object.prototype.hasOwnProperty.call(externals, dependency.request)) {
              return handleExternal(externals[dependency.request], callback);
            }

            callback();
          };

          handleExternals(_this.externals, function (err, module) {
            if (err) return callback(err);
            if (!module) return handleExternal(false, callback);
            return callback(null, module);
          });
        };
      });
    }
  }]);

  return ExternalModuleFactoryPlugin;
}();

module.exports = ExternalModuleFactoryPlugin;