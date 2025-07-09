import json
import redis


async def store_candles(rdb: redis.Redis, candles):
    pipe = rdb.pipeline()
    for c in candles:
        key = f"candle:{c.ticker}:{c.timeframe}:{c.timestamp}"
        pipe.set(key, json.dumps(c.to_dict()))
        pipe.expire(key, 60*60)
    await pipe.execute()