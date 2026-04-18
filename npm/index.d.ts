export interface ModelRecord {
  provider: "openai" | "anthropic" | "google" | "meta" | "mistral" | "cohere" | "other";
  prompt_usd_per_1k_tokens: number;
  completion_usd_per_1k_tokens: number;
  cached_prompt_usd_per_1k_tokens?: number;
  context_window_tokens?: number;
  deprecated?: boolean;
  deprecation_notes?: string;
}

export interface Table {
  schema_version: string;
  generated_at: string;
  source_notes?: string;
  models: Record<string, ModelRecord>;
}

export function load(): Table;
export function get(model: string): ModelRecord | undefined;

declare const _default: { load: typeof load; get: typeof get };
export default _default;
