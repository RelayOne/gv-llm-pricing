# @relayone/llm-pricing

Canonical per-model LLM pricing table, bundled as an npm package. Source of
truth: [github.com/RelayOne/gv-llm-pricing](https://github.com/RelayOne/gv-llm-pricing).

## Install

```
npm install @relayone/llm-pricing
```

Requires Node.js >= 20.

## Usage

### ESM

```js
import { load, get } from "@relayone/llm-pricing";

const table = load();                 // full pricing table (safe to mutate)
const rates = get("gpt-4o");          // single model lookup, or undefined
console.log(rates.prompt_usd_per_1k_tokens);
```

### CommonJS

```js
const { load, get } = require("@relayone/llm-pricing");

const table = load();
const rates = get("claude-opus-4");
```

### Direct JSON import

```js
import pricing from "@relayone/llm-pricing/pricing.json" with { type: "json" };
```

## Schema

See the [root README](https://github.com/RelayOne/gv-llm-pricing#readme) for
schema fields, versioning policy, and how updates are released.
