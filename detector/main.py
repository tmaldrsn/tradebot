import asyncio
import json
import os

from dotenv import load_dotenv
from src.core.patterns import detect_swing_points
from src.infra.redis_store import get_recent_candles
from redis.asyncio import Redis

load_dotenv('../.env')

CHANNEL = "marketdata:fetched"

async def handle_message(message):
    try:
        data = json.loads(message["data"])
        ticker = data["ticker"]
        timeframe = data["timeframe"]
        print(f"📩 Event received for {ticker} @ {timeframe}")

        # Get candles for analysis
        redis_host = os.getenv('REDIS_HOST')
        if not redis_host:
            raise Exception("Environment variable `REDIS_HOST` is not set.")
        rdb = Redis(host=redis_host, port=6379, decode_responses=True)

        # Detect patterns
        matches = detect_swing_points(candles)
        for m in matches:
            print(f"📈 Detected swing point: {m}")
    except Exception as e:
        print(f"❌ Error handling message: {e}")


async def main():
    redis_host = os.getenv('REDIS_HOST')
    if not redis_host:
        raise Exception("Environment variable `REDIS_HOST` is not set.")
    rdb = Redis(host=redis_host, port=6379, decode_responses=True)

    pubsub = rdb.pubsub()
    await pubsub.subscribe(CHANNEL)
    print(f"👂 Subscribed to Redis channel: {CHANNEL}")

    async for message in pubsub.listen():
        if message["type"] == "message":
            await handle_message(message)

if __name__ == "__main__":
    asyncio.run(main())