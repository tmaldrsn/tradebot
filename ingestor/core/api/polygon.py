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


async def fetch_candles(ticker, timeframe, from_, to):
    api_key = os.getenv("POLYGON_API_KEY")
    client = RESTClient(api_key=api_key)
    tf = PolygonTimeframe(timeframe)

    aggs = client.list_aggs(
        ticker=ticker,
        multiplier=tf.multiplier,
        timespan=tf.get_timespan_string(),
        from_=from_,
        to=to,
        limit=10000
    )

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

    return list(candles)