import asyncio
import datetime
import sys

import yaml
from dotenv import load_dotenv
from src.api.polygon import PolygonClient, PolygonTimeframeDTO
from src.candle.repository.redis import RedisCandleRepository
from src.candle.models import parse_timeframe_string_to_dto
from src.infra.redis import get_redis_connection, publish_event

load_dotenv('../.env')

# instantiate polygon client
polygon_client = PolygonClient()

# instantiate redis connection
rdb = get_redis_connection()
redis_candle_repository = RedisCandleRepository(rdb)


async def poll_ticker(source_name, ticker_cfg):
    while True:
        print(f"🔄 [async] Polling {ticker_cfg['ticker']} [{ticker_cfg['timeframe']}] from {source_name}")
    
        timeframe_dto = parse_timeframe_string_to_dto(ticker_cfg["timeframe"])
        timeframe = PolygonTimeframeDTO.from_generic_timeframe_dto(timeframe_dto)

        from_ = datetime.datetime.now() - datetime.timedelta(days=2)
        to = from_ + timeframe.to_timedelta()

        try:
            # Check ticker
            ticker = polygon_client.get_ticker(ticker_cfg['ticker'])
            if not ticker:
                raise Exception(f"Ticker `{ticker_cfg['ticker']}` not found!")

            # Fetch candles
            candles = polygon_client.fetch_candles(
                ticker=ticker,
                timeframe=timeframe,
                from_=from_,
                to=to
            )

            # Store candles
            await redis_candle_repository.save_candles(candles)

            # Publish event
            event = {
                "ticker": ticker.abbreviation,
                "timeframe": timeframe.model_dump(),
                "count": len(candles)
            }
            await publish_event(rdb, "marketdata:fetched", event)
        except Exception as e:
            print(f"❌ Error in poller for {ticker_cfg['ticker']}: {e}", file=sys.stderr)

        await asyncio.sleep(timeframe.to_seconds())


async def main():
    print("🚀 Starting Async Ingestor")

    with open("config.yaml") as f:
        config = yaml.safe_load(f)

    # Create async polling tasks
    tasks = []
    for source in config["sources"]:
        for ticker_cfg in source["tickers"]:
            task = asyncio.create_task(poll_ticker(source["name"], ticker_cfg))
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