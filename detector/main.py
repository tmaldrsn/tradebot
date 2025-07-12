import asyncio
import json

from dotenv import load_dotenv
from src.patterns.patterns import detect_swing_points
from src.candle.repository.redis import RedisCandleRepository
from src.infra.redis import get_redis_connection

load_dotenv('../.env')

# instantiate redis connection
rdb = get_redis_connection()
redis_candle_repository = RedisCandleRepository(rdb)


CHANNEL = "marketdata:fetched"

async def handle_message(message):
    try:
        data = json.loads(message["data"])
        ticker = data["ticker"]
        timeframe = data["timeframe"]
        print(f"📩 Event received for {ticker} @ {timeframe}")
        
        candles = await redis_candle_repository.fetch_candles(ticker, timeframe, limit=3)

        # Detect patterns
        matches = detect_swing_points(candles)
        for m in matches:
            print(f"📈 Detected swing point: {m}")
    except Exception as e:
        print(f"❌ Error handling message: {e}")


async def main():
    rdb = get_redis_connection()

    pubsub = rdb.pubsub()
    await pubsub.subscribe(CHANNEL)
    print(f"👂 Subscribed to Redis channel: {CHANNEL}")

    async for message in pubsub.listen():
        if message["type"] == "message":
            await handle_message(message)

if __name__ == "__main__":
    asyncio.run(main())