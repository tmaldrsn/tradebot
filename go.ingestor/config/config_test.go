package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadValidConfig(t *testing.T) {
	yaml := `sources:
  - name: polygon
    type: rest
    tickers:
      - ticker: AAPL
        timeframe: 1m
        polling_interval: 1m`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(yaml))
	assert.NoError(t, err)
	tmpfile.Close()

	cfg, err := LoadConfig(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cfg.Sources))
	assert.Equal(t, "polygon", cfg.Sources[0].Name)
	assert.Equal(t, "AAPL", cfg.Sources[0].Tickers[0].Ticker)
}

func TestDuplicateSourceName(t *testing.T) {
	yaml := `sources:
  - name: polygon
    type: rest
    tickers:
      - ticker: BTC
        timeframe: 1m
        polling_interval: 1m
  - name: polygon
    type: rest
    tickers:
      - ticker: ETH
        timeframe: 1m
        polling_interval: 1m`
	tmpfile, _ := os.CreateTemp("", "config-*.yaml")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(yaml))
	tmpfile.Close()

	_, err := LoadConfig(tmpfile.Name())
	assert.Error(t, err) // if validation is done during LoadConfig
}
