from redis import Redis

from src.candle.models import TickerDTO, CandleDTO

REDIS_TICKER_KEY = "ticker:{ticker}"
REDIS_CANDLE_INDEX_KEY = "candle:{ticker}:{timeframe}"
REDIS_CANDLE_KEY = "{index_key}:{timestamp}"

class RedisTickerRepository:
    def __init__(self, rdb: Redis):
        self.rdb = rdb


    def save_ticker(self, ticker: TickerDTO) -> None:
        pipe = self.rdb.pipeline()
        key = REDIS_TICKER_KEY.format(ticker=ticker.abbreviation)
        pipe.set(key, ticker.model_dump_json())
        pipe.expire(key, 60*60)
        return pipe.execute()


    def fetch_ticker(self, ticker: str) -> TickerDTO | None:
        key = f"ticker:{ticker}"
        raw_value = self.rdb.get(key)
        if not raw_value:
            return None

        return TickerDTO.model_validate_json(raw_value)


class RedisCandleRepository:
    def __init__(self, rdb: Redis):
        self.rdb = rdb


    async def save_candles(self, candles: list[CandleDTO]) -> None:
        pipe = self.rdb.pipeline()
        for c in candles:
            key = c.redis_key()
            pipe.set(key, c.model_dump_json())
            pipe.expire(key, 60*60)
        await pipe.execute()


    async def fetch_candles(self, ticker: str, timeframe: str, limit: int = 100) -> list[CandleDTO]:
        key = f"candle:{ticker}:{timeframe}:*"
        keys = await self.rdb.keys(key)
        keys = sorted(keys)[-limit:]

        if not keys:
            return []

        raw = await self.rdb.mget(keys)
        return [CandleDTO.model_validate_json(c) for c in raw if c]
    