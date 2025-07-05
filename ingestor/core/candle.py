class Candle:
    def __init__(self, ticker, timestamp, open, high, low, close, volume, timeframe, source):
        self.ticker = ticker
        self.timestamp = timestamp
        self.open = open
        self.high = high
        self.low = low
        self.close = close
        self.volume = volume
        self.timeframe = timeframe
        self.source = source

    def to_dict(self):
        return {
            "ticker": self.ticker,
            "timestamp": self.timestamp,
            "open": self.open,
            "high": self.high,
            "low": self.low,
            "close": self.close,
            "volume": self.volume,
            "timeframe": self.timeframe,
            "source": self.source,
        }
    