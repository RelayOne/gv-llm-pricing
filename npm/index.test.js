import { describe, it, expect } from "vitest";
import { load, get } from "./index.js";

describe("relayone-llm-pricing", () => {
  it("load() returns a table with 19+ models at schema v1", () => {
    const t = load();
    expect(t.schema_version).toBe("v1");
    expect(Object.keys(t.models).length).toBeGreaterThanOrEqual(19);
  });

  it("get() returns undefined for unknown model", () => {
    expect(get("does-not-exist")).toBeUndefined();
  });

  it("get('gpt-4o') has positive prompt and completion rates", () => {
    const r = get("gpt-4o");
    expect(r).toBeDefined();
    expect(r.prompt_usd_per_1k_tokens).toBeGreaterThan(0);
    expect(r.completion_usd_per_1k_tokens).toBeGreaterThan(r.prompt_usd_per_1k_tokens);
  });
});
