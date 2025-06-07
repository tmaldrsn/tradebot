package utils

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/polygon-io/client-go/rest/models"
)

type PolygonInterval struct {
	Multiplier int
	Timespan   models.Timespan
}

func ParseInterval(intervalStr string) (PolygonInterval, error) {
	// Match intervals like "1m", "5h", "2d", etc.
	re := regexp.MustCompile(`^(\d+)([mhdwMy])$`)
	matches := re.FindStringSubmatch(intervalStr)
	if len(matches) != 3 {
		return PolygonInterval{}, fmt.Errorf("invalid interval format: %s", intervalStr)
	}

	multiplier, err := strconv.Atoi(matches[1])
	if err != nil {
		return PolygonInterval{}, fmt.Errorf("invalid number in interval: %w", err)
	}

	var timespan models.Timespan
	switch matches[2] {
	case "m":
		timespan = models.Minute
	case "h":
		timespan = models.Hour
	case "d":
		timespan = models.Day
	case "w":
		timespan = models.Week
	case "M":
		timespan = models.Month
	case "y", "Y":
		timespan = models.Year
	default:
		return PolygonInterval{}, fmt.Errorf("unknown timespan suffix: %s", matches[2])
	}

	return PolygonInterval{
		Multiplier: multiplier,
		Timespan:   timespan,
	}, nil
}
