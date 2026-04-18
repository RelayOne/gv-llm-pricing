from relayone_llm_pricing import load, get, Record


def test_load_succeeds():
    t = load()
    assert t.schema_version == "v1"
    assert len(t.models) >= 19


def test_unknown_model_returns_none():
    assert get("does-not-exist") is None


def test_gpt_4o_is_populated():
    rec = get("gpt-4o")
    assert rec is not None
    assert isinstance(rec, Record)
    assert rec.prompt_usd_per_1k_tokens > 0
    assert rec.completion_usd_per_1k_tokens > rec.prompt_usd_per_1k_tokens
