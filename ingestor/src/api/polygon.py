import os
from enum import StrEnum

from polygon import RESTClient

from core.candle import Candle
from core.timeframe import Timeframe

class PolygonTimespan(StrEnum):
    MINUTE = "minute"
    HOUR = "hour"
    DAY = "day"
    WEEK = "week"
    MONTH = "month"
    QUARTER = "quarter"
    YEAR = "year"


class PolygonTimeframe(Timeframe):
    def __init__(self, timeframe):
        super().__init__(timeframe)

    def get_timespan_string(self):
        return {
            "m": PolygonTimespan.MINUTE,
            "h": PolygonTimespan.HOUR,
            "d": PolygonTimespan.DAY,
            "w": PolygonTimespan.WEEK,
            "M": PolygonTimespan.MONTH,
            "q": PolygonTimespan.QUARTER,
            "y": PolygonTimespan.YEAR,
        }[self.timespan]


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