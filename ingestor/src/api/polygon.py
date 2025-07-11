import os
from enum import StrEnum

from polygon import RESTClient

from core.candle import Candle
from core.timeframe import Timeframe
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



def get_polygon_client():
    api_key = os.getenv("POLYGON_API_KEY")
    if not api_key:
        raise ValueError("POLYGON_API_KEY environment variable is required")

    return RESTClient(api_key=api_key)


async def fetch_candles(ticker, timeframe, from_, to):
    client = get_polygon_client()
    tf = PolygonTimeframe(timeframe)

    try:
        aggs = client.list_aggs(
            ticker=ticker,
            multiplier=tf.multiplier,
            timespan=tf.get_timespan_string(),
            from_=from_,
            to=to,
            limit=10000
        )
    except Exception as e:
        raise RuntimeError(f"Failed to fetch candles from Polygon API: {e}") from e

    candles = []
    for agg in aggs:
        candles.append(Candle(
            ticker=ticker,
            timestamp=agg.timestamp,
            open=agg.open,
            high=agg.high,
            low=agg.low,
            close=agg.close,
            volume=agg.volume,
            timeframe=tf.timeframe,
            source="polygon",
        ))

    return candles