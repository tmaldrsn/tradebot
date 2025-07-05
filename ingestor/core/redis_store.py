import json


async def store_candles(rdb, candles):
    pipe = rdb.pipeline()
    for c in candles:
        key = f"candle:{c.ticker}:{c.timeframe}:{c.timestamp}"
        pipe.set(key, json.dumps(c.to_dict()))
    await pipe.execute()