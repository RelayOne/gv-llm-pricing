const { readFileSync } = require("node:fs");
const { join } = require("node:path");

function load() {
  const path = join(__dirname, "v1", "pricing.json");
  return JSON.parse(readFileSync(path, "utf8"));
}

function get(model) {
  return load().models[model];
}

module.exports = { load, get, default: { load, get } };
