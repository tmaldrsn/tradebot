from src.candle.models import CandleDTO

def detect_swing_points(candles: list[CandleDTO]):
    results = []

    for i in range(1, len(candles) - 1):
        prev, curr, next_ = candles[i - 1], candles[i], candles[i + 1]
        is_swing_high = curr.high > prev.high and curr.high > next_.high
        is_swing_low = curr.low < prev.low and curr.low < next_.low

        if is_swing_high or is_swing_low:
            results.append({
                "timestamp": curr.timestamp,
                "ticker": str(curr.ticker),
                "timeframe": str(curr.timeframe),
                "type": "high" if is_swing_high else "low"
            })

    return results