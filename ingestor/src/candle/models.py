import datetime
import re
from enum import Enum

from pydantic import BaseModel, PositiveFloat, PositiveInt


class Timespan(str, Enum):
    SECOND = "s"
    MINUTE = "m"
    HOUR = "h"
    DAY = "d"
    WEEK = "w"
    MONTH = "m"
    QUARTER = "q"
    YEAR = "y"

TIMESPAN_ABBR_TO_TIMESPAN = {
    "s": Timespan.SECOND,
    "m": Timespan.MINUTE,
    "h": Timespan.HOUR,
    "d": Timespan.DAY,
    "w": Timespan.WEEK,
    "M": Timespan.MONTH,
    "q": Timespan.QUARTER,
    "y": Timespan.YEAR,
}

class TimeframeDTO(BaseModel):
    multiplier: PositiveInt
    timespan: Timespan

    def __str__(self):
        return f"{self.multiplier}{self.timespan.value}"


    def to_seconds(self) -> int:
        if self.timespan in 'Mqy':
            raise NotImplementedError("The number of seconds for months, quarters, or years cannot be calculated.")

        mapping = {
            "s": 1,
            "m": 60,
            "h": 60 * 60,
            "d": 60 * 60 * 24,
            "w": 60 * 60 * 24 * 7,
        }
        return self.multiplier * mapping[self.timespan.value]

    
    def to_timedelta(self) -> datetime.timedelta:
        return datetime.timedelta(seconds=self.to_seconds())


def parse_timeframe_string_to_dto(timeframe_string: str) -> TimeframeDTO:
    RE_TIMEFRAME = re.compile(r'^([0-9]+)([smhdwMqy])$')

    match = RE_TIMEFRAME.match(timeframe_string)
    if not match:
        raise ValueError(f"Invalid timeframe string: '{timeframe_string}'")
    
    return TimeframeDTO(
        multiplier = int(match.groups()[0]),
        timespan = TIMESPAN_ABBR_TO_TIMESPAN[match.groups()[1]]
    )


class TickerDTO(BaseModel):
    abbreviation: str

    def __str__(self):
        return self.abbreviation


class CandleDTO(BaseModel):
    ticker: TickerDTO
    timestamp: PositiveInt
    open: PositiveFloat
    high: PositiveFloat
    low: PositiveFloat
    close: PositiveFloat
    volume: PositiveFloat
    timeframe: TimeframeDTO
    source: str

    def redis_index_key(self):
        return f"candle:{self.ticker.abbreviation}:{self.timeframe}"

    def redis_key(self):
        return f"{self.redis_index_key()}:{self.timestamp}"
    