from dataclasses import dataclass

from src.candle.models import CandleDTO


@dataclass
class Candle:
    timestamp: int
    open: float
    high: float
    low: float
    close: float
    volume: float

    @classmethod
    def from_dto(cls, dto: CandleDTO) -> "Candle": 
        dto = dto.validate_olhc_relationships()
        return cls(
            timestamp=dto.timestamp,
            open_=dto.open,
            high=dto.high,
            low=dto.low,
            close=dto.close,
            volume=dto.volume
        )
