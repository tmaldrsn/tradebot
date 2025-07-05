import datetime
import asyncio
import signal
import sys

from core.config import load_config
from core.pubsub import publish_event
from core.redis_store import store_candles
from core.api.polygon import fetch_candles
from core.timeframe import Timeframe
from redis.asyncio import Redis

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
    rdb = Redis(host="redis", port=6379, decode_responses=True)
    await rdb.ping()

    # Create async polling tasks
    tasks = []
    for source in config["sources"]:
        for ticker_cfg in source["tickers"]:
            task = asyncio.create_task(poll_ticker(source["name"], rdb, ticker_cfg))
            tasks.append(task)

    # Shutdown handler
    def handle_shutdown():
        print("\n🛑 Received shutdown signal")
        shutdown.set()

    loop = asyncio.get_running_loop()
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, handle_shutdown)

    await shutdown.wait()
    for task in tasks:
        task.cancel()
    await asyncio.gather(*tasks, return_exceptions=True)

    print("✅ Async Ingestor exited cleanly.")

if __name__ == "__main__":
    asyncio.run(main())