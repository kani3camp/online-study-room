"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.addWebpackRegister = addWebpackRegister;

var Template = require("webpack/lib/Template"); // eslint-disable-next-line import/prefer-default-export


function addWebpackRegister(source) {
  if (source) {
    var splitSource = source.split("jsonpArray.push = webpackJsonpCallback;");
    return Template.asString([splitSource[0].replace("var oldJsonpFunction = jsonpArray.push.bind(jsonpArray);", Template.asString(['var webpackRegister = window["webpackRegister"] = window["webpackRegister"] || [];', "var oldJsonpFunction = jsonpArray.push.bind(jsonpArray);"])), "jsonpArray.push = function(data) {", Template.indent("webpackJsonpCallback(data)"), Template.indent(["data[0].forEach(function(chunkId) {", Template.indent(["if (interleaveDeferred[chunkId]) {", Template.indent("interleaveDeferred[chunkId].resolver[0](interleaveDeferred);"), "}"]), "});"]), "};", "webpackRegister.push = registerLocals;", splitSource[1]]);
  }
}