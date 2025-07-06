import asyncio
import datetime
import os
import sys

from dotenv import load_dotenv
from redis.asyncio import Redis
from src.api.polygon import fetch_candles
from src.core.config import load_config
from src.core.timeframe import Timeframe
from src.infra.pubsub import publish_event
from src.infra.redis_store import store_candles

load_dotenv('../.env')

# Graceful shutdown event
shutdown = asyncio.Event()

async def poll_ticker(source_name, rdb, ticker_cfg):
    while not shutdown.is_set():
        print(f"🔄 [async] Polling {ticker_cfg['ticker']} [{ticker_cfg['timeframe']}] from {source_name}")
        
        tf = Timeframe(ticker_cfg["timeframe"])

        from_ = datetime.datetime.now() - datetime.timedelta(days=2)
        to = from_ + tf.to_timedelta()

        try:
            # Fetch candles
            candles = await fetch_candles(
                ticker=ticker_cfg["ticker"],
                timeframe=ticker_cfg["timeframe"],
                from_=from_,
                to=to
            )

            # Store candles
            await store_candles(rdb, candles)

            # Publish event
            event = {
                "ticker": ticker_cfg["ticker"],
                "timeframe": ticker_cfg["timeframe"],
                "count": len(candles)
            }
            await publish_event(rdb, "marketdata:fetched", event)
        except Exception as e:
            print(f"❌ Error in poller for {ticker_cfg['ticker']}: {e}", file=sys.stderr)

        await asyncio.sleep(tf.to_seconds())

async def main():
    print("🚀 Starting Async Ingestor")

    # Load config
    config = load_config("config.yaml")

    # Initialize async Redis client
    redis_host = os.getenv("REDIS_HOST")
    if not redis_host:
        raise Exception("Environment variable `REDIS_HOST` is not set.")
    rdb = Redis(host=redis_host, port=6379, decode_responses=True)

    # Create async polling tasks
    tasks = []
    for source in config["sources"]:
        for ticker_cfg in source["tickers"]:
            task = asyncio.create_task(poll_ticker(source["name"], rdb, ticker_cfg))
            tasks.append(task)

    try:
        print("🔁 Ingestor running. Press Ctrl+C to exit.")
        await asyncio.gather(*tasks)
    except asyncio.CancelledError:
        print("🛑 Cancelling tasks...")
    except KeyboardInterrupt:
        print("\n🛑 KeyboardInterrupt received. Cancelling tasks...")
        for task in tasks:
            task.cancel()
        await asyncio.gather(*tasks, return_exceptions=True)

    print("✅ Async Ingestor exited cleanly.")

if __name__ == "__main__":
    asyncio.run(main())