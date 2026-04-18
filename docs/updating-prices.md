# Updating Prices

A step-by-step guide for maintainers of `gv-llm-pricing`. If you have never
published a pricing update before, read this document end-to-end once, then
follow the numbered recipe in **The update flow** section.

## Overview

This repository is the canonical source of list prices for the LLM providers
that `router-core` (and other RelayOne billing consumers) price against.

You update it when a provider changes their pricing page. The trigger is
almost always external: OpenAI, Anthropic, Google, Meta (via hosts), Mistral,
or Cohere announces a new rate card, a new model, or deprecates an old one.
You fetch the new numbers, edit `v1/pricing.json`, and cut a release.

Pricing *values* are not breaking changes. A price going from
`$3.00 / 1M input tokens` to `$2.50 / 1M input tokens` is a PATCH bump
(`1.0.0` → `1.0.1`). Consumers pin to a tag, so they will not pick up the new
number until they bump their pin — that is the intended audit trail.

The semver rules in one sentence: **PATCH** for price value edits, **MINOR**
for additive changes (new model, new optional schema field, deprecation
marker), **MAJOR** for removals or schema-shape changes that break existing
consumers.

Because PATCH releases are cheap, do not batch price updates. Ship each
provider's update as its own release when you can. Smaller diffs are easier to
review and easier to roll back via a superseding release if a number is wrong.

## Before you start

Make sure you have all of the following:

- **Write access** to `RelayOne/gv-llm-pricing` on GitHub, and permission to
  push tags to `origin`. If `git push origin vX.Y.Z` fails for you, stop and
  ask an org admin — tagging without push access leaves the repo in a bad
  state.
- **A local clone** of the repo with `main` up to date:

  ```
  git clone git@github.com:RelayOne/gv-llm-pricing.git
  cd gv-llm-pricing
  git checkout main
  git pull --ff-only origin main
  ```

- **Validators installed locally** so you can check your edits before pushing:
  - `jq` — for sanity-checking `v1/pricing.json` parses at all.
  - `python3` with `check-jsonschema` (`pip install check-jsonschema`) — for
    validating `v1/pricing.json` against `v1/schema.json`.
  - Alternatively, the Go test suite: `go test ./tests/` runs
    `tests/schema_test.go` which does the same validation natively. If you
    already have Go installed, this is the least setup.
- **Provider pricing pages open in a browser** so you can copy exact numbers.
  URLs are listed in step 2 below.
- **A text editor** that will not mangle JSON formatting (trailing commas,
  tabs-vs-spaces). The existing file uses two-space indentation.

## The update flow

Follow these steps in order. Do not skip validation.

1. **Pull latest `main`.**

   ```
   git checkout main
   git pull --ff-only origin main
   ```

   Start from a clean working tree. If `git status` shows anything, stash or
   commit it before continuing.

2. **Fetch current list prices from each provider you are updating.** Use the
   official public pricing page for the provider. Do not copy numbers from
   blog posts, third-party aggregators, or cached search results — those lag
   the real rate card.

   - OpenAI: https://openai.com/api/pricing/
   - Anthropic: https://www.anthropic.com/pricing
   - Google (Gemini): https://ai.google.dev/pricing
   - Meta (Llama models): pulled via hosting providers. Together AI
     (https://www.together.ai/pricing) is the reference for list prices in
     this repo.
   - Mistral: https://mistral.ai/pricing
   - Cohere: https://cohere.com/pricing

   Record the exact URL you pulled from for each provider — you will paste it
   into `source_notes` in step 4.

3. **Edit `v1/pricing.json`.** Update prices on existing models, add any new
   models, and mark deprecated ones. Keep the field shape exactly as the
   schema defines — do not invent new keys. If a provider has added a
   genuinely new pricing axis (for example, caching discounts on a new
   dimension) that does not fit the current schema, stop here and file an
   issue. Schema changes are their own PR, separate from price edits.

   Prices are per-million-token rates in USD unless the schema says
   otherwise. Double-check the unit on the provider page — some pages default
   to per-1K, some to per-1M.

4. **Update `source_notes`** with the fetch date and the URLs you actually
   pulled from. Example:

   ```
   "source_notes": "Fetched 2026-04-18 from: OpenAI https://openai.com/api/pricing/, Anthropic https://www.anthropic.com/pricing"
   ```

   The date is the day you did the fetch, not the release date.

5. **Update `generated_at`** to an RFC 3339 timestamp in UTC. For example:
   `"2026-04-18T14:30:00Z"`. On Linux/macOS, `date -u +%Y-%m-%dT%H:%M:%SZ`
   prints exactly this format.

6. **Run validation locally.** Use whichever of these you have installed:

   ```
   go test -count=1 ./tests/
   ```

   or

   ```
   check-jsonschema --schemafile v1/schema.json v1/pricing.json
   ```

   Both must pass. If either fails, fix the JSON before going further. Common
   causes of failure are trailing commas, missing required fields on a newly
   added model, and wrong types (string where a number is expected).

7. **Update `CHANGELOG.md`.** Open the top of the file and you will see an
   `## [Unreleased]` section. Do this:

   - Rename `## [Unreleased]` to `## [X.Y.Z] — YYYY-MM-DD` with the version
     you are about to release and today's date.
   - Add a fresh empty `## [Unreleased]` above it so the next maintainer has
     somewhere to write.
   - The version number follows semver (see **Versioning recipes** below):
     - PATCH for a price-value-only update.
     - MINOR if you added a new model or a new optional schema field or a
       deprecation marker.
     - MAJOR if you removed a model or changed the schema shape.

8. **Open a PR against `main`.** Push your branch, open the PR, ask for at
   least one review from another `@RelayOne/router-core-maintainers` member.
   Do not push price changes directly to `main`. The PR is the audit log.
   Merge via squash-or-merge once reviewed.

9. **Tag the merge commit.** After merge, pull `main` again, then:

   ```
   git checkout main
   git pull --ff-only origin main
   git tag vX.Y.Z -m "release vX.Y.Z"
   git push origin vX.Y.Z
   ```

   The tag must match the version you wrote into `CHANGELOG.md`. Do not
   prefix the tag with anything else (no `release/`, no `pricing/`).

10. **Watch GitHub Actions.** Three workflows run on tag push:
    - `validate.yml` — schema + tests. Must go green.
    - `release.yml` — publishes to PyPI (`relayone-llm-pricing`) and npm
      (`@relayone/llm-pricing`). Must go green.
    - `sign.yml` — signs the release artifacts. Must go green.

    If any fails, do not retag. Diagnose the failure, push a fix PR, merge,
    and cut a superseding version (e.g. `1.2.4` instead of trying to re-use
    `1.2.3`). Tags should never be force-pushed.

11. **Verify the release landed** on both package registries:

    ```
    pip install --upgrade relayone-llm-pricing
    python -c "import relayone_llm_pricing; print(relayone_llm_pricing.__version__)"
    ```

    ```
    npm view @relayone/llm-pricing version
    ```

    Both should print the version you just tagged. If one is missing even
    after all three workflows went green, check the workflow logs — usually a
    transient registry 5xx that will resolve on retry.

## Versioning recipes

Pick the bump by the *shape* of the diff, not the number of lines changed.

- **Price value change on an existing model** → PATCH. Example:
  `1.0.0` → `1.0.1`. This is the common case.
- **New model added, schema unchanged** → MINOR. Example: `1.0.1` → `1.1.0`.
  Existing consumers keep working, new consumers get the new model.
- **Model deprecation marker added** (model stays in the file, flagged as
  deprecated) → MINOR. Deprecation is additive information.
- **Model fully removed from the file** → MAJOR. A consumer might have
  pinned to a version that referenced this model and then upgraded; the
  lookup would now return `not found`. Treat it as breaking. Example:
  `1.5.2` → `2.0.0`.
- **New optional schema field added** (existing consumers ignore it) →
  MINOR.
- **Schema field renamed, removed, or type-changed** → MAJOR. Existing
  consumers will fail to parse or will read the wrong value.

When in doubt, bump higher. A spurious MAJOR is annoying; a missed MAJOR is
a production incident.

## How consumers notice

Consumers (`router-core` and any other tool that embeds this data) pin to a
specific tag of this repo. They do not auto-upgrade. To pick up new prices,
a consumer bumps the pin in their own repo, runs tests, and redeploys. That
upgrade is the audit trail — the git history on the *consumer side* records
exactly when they adopted a given pricing snapshot.

For `router-core` specifically, `router-core --version` prints the pinned
pricing tag alongside the router-core version. Operators can sanity-check
which pricing data is in effect in production without reading config files.

If you make a pricing change that consumers should pick up urgently (for
example, a provider cut prices materially and billing systems are now
over-charging), post in the release notes. Do not try to push the change
into consumers for them — forcing an upgrade across consumers is a
router-core concern, not a gv-llm-pricing concern.

## Common mistakes

- **Editing `v1/schema.json` casually.** The schema is a stability boundary.
  Changing it breaks every consumer that validates against it. Schema
  changes are a separate PR, a separate review cycle, and (for shape
  changes) a MAJOR release.
- **Hand-editing `python/src/.../v1/pricing.json` or `npm/v1/pricing.json`.**
  These are synced from the root `v1/pricing.json` by the release workflow.
  Edit the root file only; the wrappers pick it up automatically at release
  time. If you hand-edit the copies, your changes get overwritten on the
  next release and the wrappers drift from the canonical source.
- **Pushing price changes straight to `main`.** No. PR review is cheap and
  catches typos before they hit PyPI. The PR log is also the
  auditable-history-of-who-changed-what-and-why.
- **Forgetting to update `generated_at` and `source_notes`.** These are how
  downstream auditors (and your future self) know when the data was
  captured and from where. If they are stale, the file is suspect.
- **Reusing a tag after a failed release.** Once a tag is pushed, leave it.
  Fix forward with a superseding version. Force-pushing tags corrupts
  anything that has already fetched the old tag (including caches on CI
  runners and developer machines).

## Emergency fix (a published price is wrong)

If a release contains a wrong price — whether a typo, a fat-finger, or a
misread of the provider's page — do not try to delete or reset the bad tag.
Consumers may have already pinned to it. Deleting it would break their
builds and corrupt any cache that had fetched the artifact. The bad tag
stays.

Instead:

1. Publish a PATCH release immediately with the correct number. Follow the
   normal update flow, same day.
2. In the `CHANGELOG.md` entry for the new release, add a `### Changed` (or
   `### Security` if the error caused billing-correctness harm) note
   explaining the error, and include the phrase
   **`supersedes v1.2.3`** (with the bad version number) so downstream
   consumers searching for why prices moved can find the explanation.
3. Open an issue on `gv-llm-pricing` with the bad tag, linking to the
   superseding release. Close it as `fixed-by vX.Y.Z`.
4. If the bad release made it into a production consumer, that is a
   router-core incident, not a pricing-repo incident. Hand it to the
   router-core oncall. Your job here is to get the corrected numbers onto
   PyPI and npm fast.

The shape is "append, never rewrite." That is the only way consumers can
trust a pin means what it meant yesterday.
