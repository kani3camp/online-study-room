"use strict";

module.exports.requireInterleaveExtension = function () {
  /* global interleaveDeferredCopy, interleaveDeferred,installedChunks */
  // interleaveDeferredCopy, interleaveDeferred, installedChunks are globals inside the webpack runtime scope
  function interleaveCss(args) {
    // Interleaved CSS loading
    var installedChunks = args.installedChunks,
        chunkId = args.chunkId,
        foundChunk = args.foundChunk,
        finalResolve = args.finalResolve; // 0 means 'already installed'

    if (installedChunks[chunkId] !== 0) {
      installedChunks[chunkId] = new Promise(function (resolve, reject) {
        var fullhref = foundChunk.path;
        var existingLinkTags = document.getElementsByTagName("link");

        for (var i = 0; i < existingLinkTags.length; i++) {
          var tag = existingLinkTags[i];
          var linkDataHref = tag.getAttribute("data-href") || tag.getAttribute("href");
          if (tag.rel === "stylesheet" && linkDataHref === fullhref) resolve();
          return finalResolve[0]();
        }

        var existingStyleTags = document.getElementsByTagName("style");

        for (var _i = 0; _i < existingStyleTags.length; _i++) {
          var _tag = existingStyleTags[_i];

          var styleDataHref = _tag.getAttribute("data-href");

          if (styleDataHref === fullhref) interleaveDeferred[chunkId].resolver[0]();
          interleaveDeferredCopy[chunkId] = interleaveDeferred[chunkId];
          delete interleaveDeferred[chunkId];
          finalResolve[0]();
          return;
        }

        var linkTag = document.createElement("link");
        linkTag.rel = "stylesheet";
        linkTag.type = "text/css";

        linkTag.onload = function () {
          // trigger a promise resolution for anything else waiting
          interleaveDeferred[chunkId].resolver[0](); // remove from object after resolving it

          delete interleaveDeferred[chunkId]; // resolve the promise chain in this function scope

          finalResolve[0]();
        };

        linkTag.onerror = function (event) {
          var request = event && event.target && event.target.src || fullhref;
          var err = new Error("Loading CSS chunk ".concat(chunkId, " failed.\n(").concat(request, ")"));
          err.code = "CSS_CHUNK_LOAD_FAILED";
          err.request = request;
          linkTag.parentNode.removeChild(linkTag);
          reject(err);
          interleaveDeferred[chunkId].resolver[1](err);
          delete interleaveDeferred[chunkId];
          finalResolve[1](err);
        };

        linkTag.href = fullhref;

        if (linkTag.href.indexOf("".concat(window.location.origin, "/")) !== 0) {
          linkTag.crossOrigin = true;
        }

        var target = document.querySelector("body");
        target.insertBefore(linkTag, target.firstChild);
      }).then(function () {
        installedChunks[chunkId] = 0;
      });
    }
  } // registerLocals chunk loading for javascript


  __webpack_require__.interleaved = function (moduleIdWithNamespace, isNested) {
    var initialRequestMap = {};
    var interleavePromises = [];
    var finalResolve;
    var finalPromise = new Promise(function (resolve, reject) {
      finalResolve = [resolve, reject];
    });

    if (!isNested) {}

    if (isNested) {}

    var chunkId = moduleIdWithNamespace.substring(moduleIdWithNamespace.indexOf("/") + 1);
    var namespace = moduleIdWithNamespace.split("/")[0];
    var namespaceObj = window.entryManifest[namespace];
    var foundChunk = namespaceObj[chunkId] || namespaceObj["".concat(chunkId, ".js")];

    if (!foundChunk) {
      finalResolve[1]("webpack-external-import: unable to find ".concat(chunkId));
      return finalPromise;
    }

    var isCSS = chunkId.indexOf(".css") !== -1;

    if (!isNested) {
      initialRequestMap[moduleIdWithNamespace] = chunkId;
    }

    var installedChunkData = installedChunks[chunkId];

    if (installedChunkData !== 0 && !isCSS) {
      // 0 means 'already installed'.
      // a Promise means "currently loading".
      if (installedChunkData) {
        interleavePromises.push(installedChunkData[2]);
      } else {
        if (!interleaveDeferred[chunkId]) {
          // current main chunk
          var resolver;

          var _promise = new Promise(function (resolve, reject) {
            resolver = [resolve, reject];
          });

          interleaveDeferred[chunkId] = {
            promise: _promise,
            resolver: resolver
          };
        } // setup Promise in chunk cache


        var promise = new Promise(function (resolve, reject) {
          installedChunkData = installedChunks[chunkId] = [resolve, reject];
        });
        interleavePromises.push(installedChunkData[2] = promise);
        // start chunk loading
        var script = document.createElement("script");
        script.charset = "utf-8";
        script.timeout = 120;

        if (__webpack_require__.nc) {
          script.setAttribute("nonce", __webpack_require__.nc);
        }

        script.src = foundChunk.path; // create error before stack unwound to get useful stacktrace later

        var error = new Error();

        var onScriptComplete = function onScriptComplete(event) {
          // avoid mem leaks in IE.
          script.onerror = script.onload = null; // eslint-disable-next-line no-use-before-define

          clearTimeout(timeout);
          var chunk = installedChunks[chunkId];

          if (chunk !== 0) {
            if (chunk) {
              var errorType = event && (event.type === "load" ? "missing" : event.type);
              var realSrc = event && event.target && event.target.src;
              error.message = "Loading chunk ".concat(chunkId, " failed. (").concat(errorType, ": ").concat(realSrc, ")");
              error.name = "ChunkLoadError";
              error.type = errorType;
              error.request = realSrc;
              chunk[1](error);
              delete interleaveDeferred[chunkId];
              finalResolve[1](error);
            }

            installedChunks[chunkId] = undefined;
          }

          var interleaveDeferredKeys = Object.keys(interleaveDeferred);
          interleaveDeferredCopy[chunkId] = interleaveDeferred[chunkId];
          delete interleaveDeferred[chunkId];
          var chunksToInstall = interleaveDeferredKeys.filter(function (item) {
            return installedChunks[item] === undefined;
          });

          if (!chunksToInstall.length) {
            finalResolve[0]();
          } // recursively find more chunks to install and push them into the interleave function
          // once all nested calls are done, resolve the current functions promise


          Promise.all(chunksToInstall.map(function (chunk) {
            return __webpack_require__.interleaved("".concat(namespace, "/").concat(chunk), true);
          })).then(finalResolve[0]);
        };

        var timeout = setTimeout(function () {
          onScriptComplete({
            type: "timeout",
            target: script
          });
        }, 120000);
        script.onerror = script.onload = onScriptComplete;
        document.head.appendChild(script);
      }
    }

    if (installedChunks[chunkId] !== 0 && isCSS) {
      interleaveCss({
        installedChunks: installedChunks,
        chunkId: chunkId,
        foundChunk: foundChunk,
        finalResolve: finalResolve
      });
    }

    if (function () {}) {}

    return finalPromise.then(function () {
      if (!isNested) return __webpack_require__(chunkId);
    });
  };
};