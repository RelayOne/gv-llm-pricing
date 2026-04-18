import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

const __dirname = dirname(fileURLToPath(import.meta.url));

/**
 * Load the embedded pricing table. Parses on every call so the returned
 * object is safe to mutate.
 */
export function load() {
  const path = join(__dirname, "v1", "pricing.json");
  return JSON.parse(readFileSync(path, "utf8"));
}

/**
 * Convenience: load(), then look up. Returns undefined on miss.
 * @param {string} model
 */
export function get(model) {
  return load().models[model];
}

export default { load, get };
