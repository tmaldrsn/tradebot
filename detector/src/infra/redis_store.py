import json


async def get_recent_candles(rdb, ticker, timeframe, limit=100):
    # Use scan or sorted keys (your data model may vary)
    pattern = f"candle:{ticker}:{timeframe}:*"
    keys = await rdb.keys(pattern)
    keys = sorted(keys)[-limit:]

    if not keys:
        return []

    raw = await rdb.mget(keys)
    return [json.loads(c) for c in raw if c]