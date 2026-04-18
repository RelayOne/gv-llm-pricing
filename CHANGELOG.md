# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

Versioning rules specific to pricing data:

- **MAJOR** — breaking schema changes (new required field, removed field, renamed field).
- **MINOR** — new models added, new optional fields, deprecation marks.
- **PATCH** — price value updates, `source_notes` or `generated_at` refreshes.

## [Unreleased]

> This section becomes v1.0.0 when tagged. All content-complete for 1.0.0 per the spec.

### Added

- Initial v1 schema at `v1/schema.json`.
- Seed pricing data at `v1/pricing.json` for the big-5 providers:
  - **OpenAI:** gpt-4o, gpt-4o-mini, gpt-4-turbo, o1, o1-mini, text-embedding-3-small, text-embedding-3-large.
  - **Anthropic:** claude-opus-4, claude-sonnet-4, claude-haiku-4, claude-3-5-sonnet-20241022, claude-3-5-haiku-20241022.
  - **Google:** gemini-1.5-pro, gemini-1.5-flash, gemini-2.0-flash.
  - **Meta:** llama-3.1-405b, llama-3.1-70b.
  - **Mistral:** mistral-large-2, codestral-latest.

### Changed

### Removed
