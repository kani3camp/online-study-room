"use strict";

function _slicedToArray(arr, i) { return _arrayWithHoles(arr) || _iterableToArrayLimit(arr, i) || _nonIterableRest(); }

function _nonIterableRest() { throw new TypeError("Invalid attempt to destructure non-iterable instance"); }

function _iterableToArrayLimit(arr, i) { if (!(Symbol.iterator in Object(arr) || Object.prototype.toString.call(arr) === "[object Arguments]")) { return; } var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"] != null) _i["return"](); } finally { if (_d) throw _e; } } return _arr; }

function _arrayWithHoles(arr) { if (Array.isArray(arr)) return arr; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

var path = require("path");

var fse = require("fs-extra");

var createHash = require("webpack/lib/util/createHash");

var fs = require("fs");

var ExternalModuleFactoryPlugin = require("./ExternalModuleFactory"); // const FunctionModuleTemplatePlugin = require("webpack/lib/FunctionModuleTemplatePlugin");


var _require = require("./utils"),
    mergeDeep = _require.mergeDeep;

var _require2 = require("./requireExtentions"),
    addInterleaveExtension = _require2.addInterleaveExtension,
    addInterleaveRequire = _require2.addInterleaveRequire;

var _require3 = require("./beforeStartup"),
    addWebpackRegister = _require3.addWebpackRegister;

var _require4 = require("./chunkSplitting"),
    interleaveStyleConfig = _require4.interleaveStyleConfig,
    interleaveStyleJsConfig = _require4.interleaveStyleJsConfig,
    interleaveConfig = _require4.interleaveConfig,
    hasExternalizedModuleViaJson = _require4.hasExternalizedModuleViaJson;

var _require5 = require("./localVars"),
    addLocalVars = _require5.addLocalVars;

var _require6 = require("./optimizeChunk"),
    wrapChunks = _require6.wrapChunks;

var emitCountMap = new Map();

function getFileType(str) {
  var split = str.replace(/\?.*/, "").split(".");
  return split.pop();
}

var URLImportPlugin =
/*#__PURE__*/
function () {
  function URLImportPlugin(opts) {
    _classCallCheck(this, URLImportPlugin);

    var debug = (typeof v8debug === "undefined" ? "undefined" : _typeof(v8debug)) === "object" || /--debug|--inspect/.test(process.execArgv.join(" "));

    if (!opts.manifestName) {
      throw new Error("URLImportPlugin: You MUST specify a manifestName in your options. Something unique. Like {manifestName: my-special-build}");
    }

    this.opts = _objectSpread({
      publicPath: null,
      debug: debug || false,
      testPath: "src",
      basePath: "",
      manifestName: "unknown-project",
      fileName: "importManifest.js",
      writeToFileEmit: false,
      seed: null,
      filter: null,
      hashDigest: "base64",
      hashDigestLength: 5,
      context: null,
      hashFunction: "md4"
    }, opts || {});
  }

  _createClass(URLImportPlugin, [{
    key: "apply",
    value: function apply(compiler) {
      var _options$optimization,
          _options$optimization2,
          _this = this,
          _options$optimization3;

      if (this.opts.debug) {}

      var options = compiler === null || compiler === void 0 ? void 0 : compiler.options;

      if (options.externals) {
        throw new Error("URLImportPlugin: Externals must be applied via the plugin, not via webpack config object. Please see useExternals on the plugin documentation");
      } // add to the existing webpack config
      // adding a new splitChunks cache group called interleave


      var chunkSplitting = (options === null || options === void 0 ? void 0 : (_options$optimization = options.optimization) === null || _options$optimization === void 0 ? void 0 : (_options$optimization2 = _options$optimization.splitChunks) === null || _options$optimization2 === void 0 ? void 0 : _options$optimization2.cacheGroups) || {};
      chunkSplitting.stylejs = interleaveStyleJsConfig(this.opts);
      chunkSplitting.style = interleaveStyleConfig(this.opts); // interleaveConfig figures out if a file meets the paramaters for interleaving

      chunkSplitting.interleave = interleaveConfig(this.opts);

      if (options.mode === "production") {
        chunkSplitting.vendors = {
          name: "".concat(this.opts.manifestName, "-vendors"),
          test: /node_modules/,
          priority: -10,
          enforce: true,
          maxSize: 50000
        };
        Object.assign(chunkSplitting["default"], {
          maxSize: 50000
        });
      }

      Object.assign(options.optimization || {}, {
        namedChunks: true,
        // dont rename exports when hoisting and tree shaking
        providedExports: false
      }); // likely will be refactored or removed, used for entryManifest.js to map chunks (this is V1 where its outside the runtime still)

      if (this.opts.debug) {} // merge my added splitChunks config into the webpack config object passed in


      Object.assign(options.optimization.splitChunks, {
        chunks: "all"
      });
      mergeDeep(options, {
        optimization: {
          splitChunks: {
            chunks: "all",
            cacheGroups: chunkSplitting
          }
        }
      });
      Object.keys(chunkSplitting).forEach(function (key) {
        var _chunkSplitting$key;

        if (key === "interleave") {
          return;
        }

        chunkSplitting[key].automaticNamePrefix = "".concat(_this.opts.manifestName, "~").concat((chunkSplitting === null || chunkSplitting === void 0 ? void 0 : (_chunkSplitting$key = chunkSplitting[key]) === null || _chunkSplitting$key === void 0 ? void 0 : _chunkSplitting$key.automaticNamePrefix) || "");
      });
      Object.assign(options.optimization, {
        // node debugger breaks with TerserPlugin
        minimizer: this.opts.debug ? [] : options.optimization.minimizer,
        splitChunks: ((_options$optimization3 = options.optimization) === null || _options$optimization3 === void 0 ? void 0 : _options$optimization3.splitChunks) || {}
      }); // forcefully mutate it

      Object.assign(options.optimization.splitChunks, {
        chunks: "all",
        cacheGroups: chunkSplitting
      });

      if (this.opts.debug) {} // eslint-disable-next-line no-unused-vars


      compiler.hooks.thisCompilation.tap("URLImportPlugin", function (compilation) {// TODO: throw warning when changing module ID type
        // if (options.ignoreOrder) {
        //   compilation.warnings.push(
        //     new Error(
        //       `chunk ${chunk.name || chunk.id} [${pluginName}]\n`
        //           + 'Conflicting order between:\n'
        //           + ` * ${fallbackModule.readableIdentifier(
        //             requestShortener,
        //           )}\n`
        //           + `${bestMatchDeps
        //             .map(m => ` * ${m.readableIdentifier(requestShortener)}`)
        //             .join('\n')}`,
        //     ),
        //   );
        // }
      });
      var moduleAssets = {};
      var outputFolder = compiler.options.output.path;
      var outputFile = path.resolve(outputFolder, this.opts.fileName);
      var outputName = path.relative(outputFolder, outputFile);

      var moduleAsset = function moduleAsset(_ref, file) {
        var userRequest = _ref.userRequest;

        if (userRequest) {
          moduleAssets[file] = path.join(path.dirname(file), path.basename(userRequest));
        }
      };

      var emit = function emit(compilation, compileCallback) {
        var emitCount = emitCountMap.get(outputFile) - 1;
        emitCountMap.set(outputFile, emitCount);
        var seed = _this.opts.seed || {};
        var publicPath = _this.opts.publicPath != null ? _this.opts.publicPath : compilation.options.output.publicPath;
        var stats = compilation.getStats().toJson();

        if (_this.opts.debug) {}

        var files = compilation.chunks.reduce(function (f, chunk) {
          return chunk.files.reduce(function (fx, filePath) {
            var name = chunk.id ? chunk.id : null;

            if (name) {
              name = "".concat(name, ".").concat(getFileType(filePath));
            } else {
              // For nameless chunks, just map the files directly.
              name = filePath;
            }

            return fx.concat({
              path: filePath,
              chunk: chunk,
              name: name,
              isChunk: true,
              isAsset: false,
              isModuleAsset: false
            });
          }, f);
        }, []); // module assets don't show up in assetsByChunkName.
        // we're getting them this way;

        files = stats.assets.reduce(function (fx, asset) {
          var name = moduleAssets[asset.name];

          if (name) {
            return fx.concat({
              path: asset.name,
              name: name
            });
          }

          var isEntryAsset = asset.chunks.length > 0;

          if (isEntryAsset) {
            return fx;
          }

          return fx.concat({
            path: asset.name,
            name: asset.name
          });
        }, files);
        files = files.filter(function (file) {
          // Don't add hot updates to manifest
          var isUpdateChunk = file.path.includes("hot-update"); // Don't add manifest from another instance

          var isManifest = emitCountMap.get(path.join(outputFolder, file.name)) !== undefined;
          return !isUpdateChunk && !isManifest;
        });

        if (_this.opts.debug) {} // Append optional basepath onto all references.
        // This allows output path to be reflected in the manifest.


        if (_this.opts.basePath) {
          files = files.map(function (file) {
            return _objectSpread({}, file, {
              name: _this.opts.basePath + file.name
            });
          });
        }

        if (publicPath) {
          // Similar to basePath but only affects the value (similar to how
          // output.publicPath turns require('foo/bar') into '/public/foo/bar', see
          // https://github.com/webpack/docs/wiki/configuration#outputpublicpath
          files = files.map(function (file) {
            return _objectSpread({}, file, {
              path: publicPath + file.path
            });
          });
        }

        files = files.map(function (file) {
          return _objectSpread({}, file, {
            name: file.name.replace(/\\/g, "/"),
            path: file.path.replace(/\\/g, "/")
          });
        });

        if (_this.opts.filter) {
          files = files.filter(_this.opts.filter);
        }

        if (_this.opts.debug) {}

        var manifest = files.reduce(function (m, file) {
          return _objectSpread({}, m, _defineProperty({}, file.name, {
            path: file.path
          }));
        }, seed);

        if (_this.opts.debug) {}

        var isLastEmit = emitCount === 0;

        if (isLastEmit) {
          var cleanedManifest = Object.entries(manifest).reduce(function (acc, _ref2) {
            var _asset$path;

            var _ref3 = _slicedToArray(_ref2, 2),
                key = _ref3[0],
                asset = _ref3[1];

            if (!(asset === null || asset === void 0 ? void 0 : (_asset$path = asset.path) === null || _asset$path === void 0 ? void 0 : _asset$path.includes(".map"))) {
              return Object.assign(acc, _defineProperty({}, key, asset));
            }

            return acc;
          }, {});

          var serialize = function serialize(manifest) {
            return "if(!window.entryManifest) {window.entryManifest = {}}; window.entryManifest[\"".concat(_this.opts.manifestName, "\"] = ").concat(JSON.stringify(manifest));
          };

          var output = serialize(cleanedManifest);

          if (_this.opts.debug) {} // eslint-disable-next-line no-param-reassign


          compilation.assets[outputName] = {
            source: function source() {
              return output;
            },
            size: function size() {
              return output.length;
            }
          };

          if (_this.opts.writeToFileEmit) {
            fse.outputFileSync(outputFile, output);
          }
        }

        if (compiler.hooks) {
          compiler.hooks.webpackURLImportPluginAfterEmit.call(manifest);
        } else {
          compilation.applyPluginsAsync("webpack-manifest-plugin-after-emit", manifest, compileCallback);
        }
      };

      function beforeRun(comp, callback) {
        var emitCount = emitCountMap.get(outputFile) || 0;
        emitCountMap.set(outputFile, emitCount + 1);

        if (callback) {
          callback();
        }
      }

      if (compiler.hooks) {
        var _require7 = require("tapable"),
            SyncWaterfallHook = _require7.SyncWaterfallHook;

        var pluginOptions = {
          name: "URLImportPlugin",
          stage: Infinity
        }; // eslint-disable-next-line no-param-reassign

        compiler.hooks.webpackURLImportPluginAfterEmit = new SyncWaterfallHook(["manifest"]);
        compiler.hooks.compile.tap("ExternalsPlugin", function (_ref4) {
          var normalModuleFactory = _ref4.normalModuleFactory;
          new ExternalModuleFactoryPlugin(options.output.libraryTarget, _this.opts.useExternals).apply(normalModuleFactory);
        });
        compiler.hooks.compilation.tap("URLImportPlugin", function (compilation) {
          var mainTemplate = compilation.mainTemplate; // Add another webpack__require method to the webpack runtime
          // this new method will allow a interleaved component to be required and automatically download its dependencies
          // it returns a promise so the actual interleaved module is not executed until any missing dependencies are loaded

          mainTemplate.hooks.requireExtensions.tap("URLImportPlugin", function (source) {
            return [addInterleaveExtension, addInterleaveRequire].reduce(function (sourceCode, extension) {
              return extension(sourceCode, mainTemplate.requireFn, _this.opts);
            }, source);
          }); // TODO add an option for this

          if (_this.afterOptimizations) {
            // before chunk files are optimized
            compilation.hooks.beforeOptimizeChunkAssets.tap("URLImportPlugin", function (chunks) {
              // access all chunks webpack created, then add some code to each chunk file, which is run when a chunk is
              // loaded on a page as <script>
              wrapChunks(compilation, chunks);
            });
          } else {
            // adfter chunk files are optimized
            compilation.hooks.optimizeChunkAssets.tapAsync("URLImportPlugin", function (chunks, done) {
              wrapChunks(compilation, chunks);
              done();
            });
          } // Expose chunk registration functions and bindings from webpack runtime to the window
          // webpack does this and its how code splitting works. It exposes window.webpackJsonP
          // This registration system works just like webpacks, it exposes a function that allows information to be passed
          // into webpack runtime, because the function is in webpack runtime, i have access to all of webpacks internals


          mainTemplate.hooks.beforeStartup.tap("URLImportPlugin", function (source) {
            return addWebpackRegister(source, options.output.jsonpFunction);
          }); // add variables to webpack runtime which are available throughout all functions and closures within the runtime
          // localVars are like global variables for webpack, anything can access them.

          mainTemplate.hooks.localVars.tap("URLImportPlugin", addLocalVars);
        });
        compiler.hooks.compilation.tap("URLImportPlugin", function (compilation) {
          var usedIds = new Set(); // creates hashed module IDs based on the contents of the file - works like [contenthash] but for each module

          compilation.hooks.beforeModuleIds.tap("URLImportPlugin", function (modules) {
            var _this$opts;

            var provideExternals = Object.keys(((_this$opts = _this.opts) === null || _this$opts === void 0 ? void 0 : _this$opts.provideExternals) || {}); // eslint-disable-next-line no-restricted-syntax

            var _iteratorNormalCompletion = true;
            var _didIteratorError = false;
            var _iteratorError = undefined;

            try {
              for (var _iterator = modules[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
                var _module = _step.value;

                if (provideExternals.includes(_module.rawRequest)) {
                  _module.id = _module.rawRequest;
                }

                if (_module.id === null && _module.resource) {
                  var hash = createHash(_this.opts.hashFunction);
                  var resourcePath = _module.resource;

                  if (resourcePath.indexOf("?") > -1) {
                    var _resourcePath$split = resourcePath.split("?");

                    var _resourcePath$split2 = _slicedToArray(_resourcePath$split, 1);

                    resourcePath = _resourcePath$split2[0];
                  }

                  try {
                    var exports = "";

                    if (Array.isArray(_module.usedExports)) {
                      exports = _module.usedExports.join(".");
                    }

                    hash.update(fs.readFileSync(resourcePath) + exports);
                  } catch (ex) {
                    throw ex;
                  }

                  var hashId = hash.digest(_this.opts.hashDigest);
                  var len = _this.opts.hashDigestLength;

                  while (usedIds.has(hashId.substr(0, len))) {
                    len++;
                  }

                  _module.id = hashId.substr(0, len);
                  usedIds.add(_module.id);
                } else if (_this.opts.debug) {// console.log("Module with no ID", module);
                }

                var externalModule = hasExternalizedModuleViaJson(_module.resource);

                if (externalModule || false) {
                  var _module$buildMeta;

                  // add exports back to usedExports, prevents tree shaking on module
                  Object.assign(_module, {
                    usedExports: (_module === null || _module === void 0 ? void 0 : (_module$buildMeta = _module.buildMeta) === null || _module$buildMeta === void 0 ? void 0 : _module$buildMeta.providedExports) || true
                  });

                  if (externalModule) {
                    _module.id = externalModule;
                  }
                }
              }
            } catch (err) {
              _didIteratorError = true;
              _iteratorError = err;
            } finally {
              try {
                if (!_iteratorNormalCompletion && _iterator["return"] != null) {
                  _iterator["return"]();
                }
              } finally {
                if (_didIteratorError) {
                  throw _iteratorError;
                }
              }
            }
          });
        });
        compiler.hooks.compilation.tap(pluginOptions, function (_ref5) {
          var hooks = _ref5.hooks;
          // TODO: remove in ^2.2
          hooks.moduleAsset.tap(pluginOptions, moduleAsset);
        }); // writes the importManifest file containing a map of Chunk IDs to the cache busted JS files

        compiler.hooks.emit.tap(pluginOptions, emit);
        compiler.hooks.run.tap(pluginOptions, beforeRun);
        compiler.hooks.watchRun.tap(pluginOptions, beforeRun);
      } else {
        compiler.plugin("compilation", function (compilation) {
          compilation.plugin("module-asset", moduleAsset);
        });
        compiler.plugin("emit", emit);
        compiler.plugin("before-run", beforeRun);
        compiler.plugin("watch-run", beforeRun);
      }
    }
  }]);

  return URLImportPlugin;
}();

module.exports = URLImportPlugin;