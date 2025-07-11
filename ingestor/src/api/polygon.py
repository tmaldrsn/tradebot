import datetime
import os

from polygon import RESTClient
from polygon.exceptions import BadResponse
from src.candle.models import CandleDTO, TickerDTO, TimeframeDTO

TIMESPAN_ABBR_TO_POLYGON_TIMESPAN_ABBR = {
    "m": "minute",
    "h": "hour",
    "d": "day",
    "w": "week",
    "M": "month",
    "q": "quarter",
    "y": "year",
}

class PolygonTimeframeDTO(TimeframeDTO):
    def get_timespan_string(self) -> str:
        return TIMESPAN_ABBR_TO_POLYGON_TIMESPAN_ABBR[self.timespan.value]

    @classmethod
    def from_generic_timeframe_dto(cls, dto: TimeframeDTO):
        if dto.timespan.value == "s": 
            raise NotImplementedError("Polygon does not support `second` timeframes.")

        return cls(
            multiplier=dto.multiplier,
            timespan=dto.timespan
        )


class PolygonClient:
    def __init__(self, api_key=None):
        self.api_key = os.getenv("POLYGON_API_KEY") if not api_key else api_key
        assert self.api_key, "Provide `api_key` argument or set `POLYGON_API_KEY` environment variable."

        self.client = RESTClient(api_key=self.api_key)


    def get_ticker(self, ticker: str) -> TickerDTO | None:
        try:
            deats = self.client.get_ticker_details(ticker)
        except BadResponse:
            return
        
        return TickerDTO(abbreviation=deats.ticker)


    def fetch_candles(self, ticker: TickerDTO, timeframe: PolygonTimeframeDTO, from_: datetime.date, to: datetime.date) -> list[CandleDTO]:
        try:
            aggs = self.client.list_aggs(
                ticker=ticker.abbreviation,
                multiplier=timeframe.multiplier,
                timespan=timeframe.get_timespan_string(),
                from_=from_,
                to=to,
                limit=10000
            )
        except Exception as e:
            raise RuntimeError(f"Failed to fetch candles from Polygon API: {e}") from e

        candles = []
        for agg in aggs:
            candles.append(CandleDTO(
                ticker=ticker,
                timestamp=agg.timestamp,
                open=agg.open,
                high=agg.high,
                low=agg.low,
                close=agg.close,
                volume=agg.volume,
                timeframe=timeframe,
                source="polygon",
            ))

        return candles