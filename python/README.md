# relayone-llm-pricing

Canonical per-model LLM pricing table as a Python package. Bundles `v1/pricing.json` at install time.

## Install

```
pip install relayone-llm-pricing
```

## Usage

```python
from relayone_llm_pricing import load, get

# Look up a single model
rec = get("gpt-4o")
if rec is not None:
    print(rec.provider, rec.prompt_usd_per_1k_tokens, rec.completion_usd_per_1k_tokens)

# Or grab the whole table
table = load()
print(table.schema_version, len(table.models))
for name, rec in table.models.items():
    print(name, rec.provider)
```

## Source of truth

The canonical data lives at `v1/pricing.json` in the repo root at
[github.com/RelayOne/gv-llm-pricing](https://github.com/RelayOne/gv-llm-pricing). This
package vendors a copy of that file into `src/relayone_llm_pricing/v1/pricing.json`;
the release workflow syncs it before every build. See the root `README.md` for the
schema, versioning rules, and contribution process.

## License

Apache-2.0. See `LICENSE`.
