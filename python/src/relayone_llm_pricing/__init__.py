"""relayone_llm_pricing - canonical per-model LLM pricing table.

Bundles the JSON at import time via importlib.resources. The raw data
is the source of truth; this package is a convenience wrapper.
"""
from __future__ import annotations

import json
from dataclasses import dataclass
from importlib.resources import files
from typing import Optional

__all__ = ["Record", "Table", "load", "get"]


@dataclass(frozen=True)
class Record:
    provider: str
    prompt_usd_per_1k_tokens: float
    completion_usd_per_1k_tokens: float
    cached_prompt_usd_per_1k_tokens: Optional[float] = None
    context_window_tokens: Optional[int] = None
    deprecated: bool = False
    deprecation_notes: Optional[str] = None


@dataclass(frozen=True)
class Table:
    schema_version: str
    generated_at: str
    models: dict[str, Record]
    source_notes: Optional[str] = None

    def get(self, model: str) -> Optional[Record]:
        return self.models.get(model)


def load() -> Table:
    """Return the embedded pricing table. Parses on every call."""
    raw = files(__package__).joinpath("v1/pricing.json").read_text(encoding="utf-8")
    doc = json.loads(raw)
    models = {
        name: Record(
            provider=rec["provider"],
            prompt_usd_per_1k_tokens=float(rec["prompt_usd_per_1k_tokens"]),
            completion_usd_per_1k_tokens=float(rec["completion_usd_per_1k_tokens"]),
            cached_prompt_usd_per_1k_tokens=(
                float(rec["cached_prompt_usd_per_1k_tokens"])
                if "cached_prompt_usd_per_1k_tokens" in rec else None
            ),
            context_window_tokens=(
                int(rec["context_window_tokens"])
                if "context_window_tokens" in rec else None
            ),
            deprecated=bool(rec.get("deprecated", False)),
            deprecation_notes=rec.get("deprecation_notes"),
        )
        for name, rec in doc["models"].items()
    }
    return Table(
        schema_version=doc["schema_version"],
        generated_at=doc["generated_at"],
        source_notes=doc.get("source_notes"),
        models=models,
    )


def get(model: str) -> Optional[Record]:
    """Convenience: load() then look up. Returns None on miss."""
    return load().get(model)
