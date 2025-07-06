import datetime
import re

class InvalidTimeframeException(Exception):
    pass


class Timeframe:
    RE_TIMEFRAME = re.compile(r'^([0-9]+)([smhdwMqy])$')

    def __init__(self, timeframe):
        match = Timeframe.RE_TIMEFRAME.match(timeframe)
        if not match:
            raise InvalidTimeframeException(f"Invalid timeframe string: '{timeframe}'")
        
        self.timeframe = timeframe
        self.multiplier = int(match.groups()[0])
        self.timespan = match.groups()[1]


    def __str__(self):
        return self.timeframe


    def to_seconds(self):
        if self.timespan in 'Mqy':
            raise NotImplementedError("The number of seconds for months, quarters, or years cannot be calculated.")

        mapping = {
            "s": 1,
            "m": 60,
            "h": 60 * 60,
            "d": 24 * 60 * 60,
            "w": 7 * 24 * 60 * 60
        }
        return self.multiplier * mapping[self.timespan]

    
    def to_timedelta(self):
        return datetime.timedelta(seconds=self.to_seconds())

