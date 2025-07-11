import json

from redis import Redis

from src.candle.models import TickerDTO

REDIS_TICKER_KEY = "ticker:{ticker}"

class RedisTickerRepository:
    def __init__(self, rdb: Redis):
        self.rdb = rdb


    def save_ticker(self, ticker: TickerDTO) -> None:
        pipe = self.rdb.pipeline()
        key = REDIS_TICKER_KEY.format(ticker=ticker.abbreviation)
        pipe.set(key, TickerDTO.model_dump())
        pipe.expire(key, 60*60)
        return pipe.execute()


    def fetch_ticker(self, ticker: str) -> TickerDTO | None:
        pattern = f"ticker:{ticker}"
        keys = self.rdb.keys(pattern)
        if not keys:
            print(f"TickerDTO `{ticker}` not found...")
            return

        raw = self.rdb.mget(keys)
        raw_json = json.loads(raw)
        return TickerDTO(raw_json).model_validate()

