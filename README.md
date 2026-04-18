# gv-llm-pricing

**gv-llm-pricing — canonical per-model LLM pricing table for RelayOne / GoodVentures products.**

## What this repo is

This repository is the **JSON source of truth** for per-model LLM pricing used
across the RelayOne / GoodVentures product family. It is **versioned by git
tag** (semver), signed on release with [cosign](https://github.com/sigstore/cosign),
and consumed by:

- **router-core** (Go) — embeds the JSON at build time via a code generator.
- Future **TrustPlane** and **RelayOne** services (any language) — either via
  a thin language wrapper or by fetching the raw tagged URL.
- **Third-party ops dashboards** — read the raw JSON from a pinned tag.

The repo is intentionally tiny: one JSON file per schema major version, a JSON
Schema describing it, and some thin language wrappers. There is **no
application code** here. Pricing changes ship as release tags; consumers pin
to a tag for billing determinism.

## Schema version

Current schema: **v1**.

The schema lives at `v1/schema.json`. The `v1/` directory is a **stability
boundary**: within major version 1, the schema does not change — only price
values, model entries, and optional fields evolve. Any breaking schema change
requires a new top-level directory (`v2/`, `v3/`, ...) and a MAJOR tag bump.

## Consumption patterns

### Go (router-core and related services)

Preferred pattern — **compile-time embed** via the code generator in
router-core:

```bash
# From inside the router-core repo:
go run ./cmd/gen-pricing -tag v1.2.3
```

The generator fetches the JSON at the given tag, verifies the cosign
signature against the expected OIDC identity, and writes a Go file with
the embedded bytes. The build then compiles the prices in — no runtime
fetch, no surprises.

Alternatively, for small tools or scripts, import the thin wrapper:

```go
import "github.com/RelayOne/gv-llm-pricing/go/pricing"

table, err := pricing.Load()
if err != nil {
    // ...
}
```

### Python

Install the packaged wheel (bundles the same JSON):

```bash
pip install relayone-llm-pricing
```

```python
from relayone_llm_pricing import load
table = load()
```

Or fetch the raw JSON at a specific tag:

```bash
curl -fsSL https://raw.githubusercontent.com/RelayOne/gv-llm-pricing/vX.Y.Z/v1/pricing.json
```

### Node

Install the packaged module (bundles the same JSON):

```bash
npm install @relayone/llm-pricing
```

```js
import { load } from "@relayone/llm-pricing";
const table = load();
```

Or fetch the raw JSON at a specific tag:

```bash
curl -fsSL https://raw.githubusercontent.com/RelayOne/gv-llm-pricing/vX.Y.Z/v1/pricing.json
```

## Tag / versioning rules (semver)

This repo follows [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
with pricing-specific interpretation:

- **MAJOR** — schema breaking change (new required field, removed field,
  renamed field). Consumers **must update** their wrappers and regenerate
  embedded bytes.
- **MINOR** — new models added, new optional fields introduced, deprecation
  marks applied to existing entries. Consumers can upgrade without code
  changes.
- **PATCH** — price value updates, `source_notes` edits, or `generated_at`
  refreshes. No schema movement. Billing consumers should still pin and
  re-tag deliberately.

## Stability promise

- The **schema** does NOT change within a major version. Once `v1/schema.json`
  is published under a given major line, fields are additive-only (optional
  new fields) until a MAJOR bump.
- **Price values MAY change on any release**, including PATCH. Consumers
  **pin by tag** for billing determinism — never track `main`, never track a
  mutable branch.

## How to audit a specific tag

Every release is signed with cosign in keyless OIDC mode via GitHub Actions.
To verify a downloaded `pricing.json` against its signature and the expected
signer identity:

```bash
cosign verify-blob \
  --certificate-identity-regexp 'https://github.com/RelayOne/gv-llm-pricing/.github/workflows/sign\.yml@.*' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  --signature pricing.json.sig \
  --certificate pricing.json.pem \
  pricing.json
```

The OIDC identity regex `https://github.com/RelayOne/gv-llm-pricing/.github/workflows/sign\.yml@.*`
pins the signer to **this repo's signing workflow** — any other GitHub
Actions workflow or identity will fail verification.

## Consumer guidance

- **Always pin by tag**, never by branch or `main`. A pinned tag gives you
  byte-for-byte reproducibility and a signed provenance chain.
- **Runtime HTTP fetch is an anti-pattern** for billing data — you inherit
  transient network failures and silently changing prices. Prefer compile-time
  embed (router-core's generator) or an immutable bundled artifact (the
  Python and Node packages). If you must fetch at runtime, fetch once at
  boot, pin the tag, and cache on disk.
- **Verify the cosign signature** at fetch time, at least in CI. The raw
  GitHub CDN is not a trust boundary.

## License

Licensed under the [Apache License 2.0](./LICENSE).

## Contributing / updating prices

See `docs/updating-prices.md` (added in P-5) for the full workflow:
where prices come from, how to propose a diff, and how the signing workflow
cuts a signed release tag.
