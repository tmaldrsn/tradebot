import asyncio
import json
from redis.asyncio import Redis
from core.redis_store import get_recent_candles
from core.patterns import detect_swing_points

CHANNEL = "marketdata:fetched"

async def handle_message(message):
    try:
        data = json.loads(message["data"])
        ticker = data["ticker"]
        timeframe = data["timeframe"]
        print(f"📩 Event received for {ticker} @ {timeframe}")

        # Get candles for analysis
        rdb = Redis(host="redis", port=6379, decode_responses=True)
        candles = await get_recent_candles(rdb, ticker, timeframe, limit=100)

        # Detect patterns
        matches = detect_swing_points(candles)
        for m in matches:
            print(f"📈 Detected swing point: {m}")
    except Exception as e:
        print(f"❌ Error handling message: {e}")

async def main():
    rdb = Redis(host="redis", port=6379, decode_responses=True)
    pubsub = rdb.pubsub()
    await pubsub.subscribe(CHANNEL)
    print(f"👂 Subscribed to Redis channel: {CHANNEL}")

    async for message in pubsub.listen():
        if message["type"] == "message":
            await handle_message(message)

if __name__ == "__main__":
    asyncio.run(main())