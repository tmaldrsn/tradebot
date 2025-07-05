import aiohttp
import os
import time

async def fetch_candles(ticker, timeframe):
    api_key = os.getenv("POLYGON_API_KEY")
    start = "2025-06-01"
    end = "2025-06-02"

    url = f"https://api.polygon.io/v2/aggs/ticker/X:{ticker}/range/1/{timeframe}/{start}/{end}?adjusted=true&sort=asc&limit=120&apiKey={api_key}"

    async with aiohttp.ClientSession() as session:
        async with session.get(url) as resp:
            data = await resp.json()

    candles = []
    for result in data.get("results", []):
        candles.append({
            "ticker": ticker,
            "timestamp": result["t"] // 1000,
            "open": result["o"],
            "high": result["h"],
            "low": result["l"],
            "close": result["c"],
            "volume": result["v"],
            "timeframe": timeframe,
            "source": "polygon"
        })
    return candles