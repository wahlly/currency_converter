package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type RatesCache struct {
	mutex	sync.RWMutex
	Base	string
	Rates	map[string]float64
	Fetched	time.Time
	TTL	time.Duration
}

func NewRatesCache(ttl time.Duration) *RatesCache {
	return &RatesCache{Rates: make(map[string]float64), TTL: ttl}
}

var (
	XCHANGE_RATE_BASE_URL = os.Getenv("XCHANGE_RATE_BASE_URL")
	XCHANGE_RATE_API_KEY = os.Getenv("XCHANGE_RATE_API_KEY")
)

type ClientApiRes struct {
	Success bool `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Base string `json:"base"`
	Date string `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func FetchRates(base string, symbols []string) (map[string]float64, error) {
	url := fmt.Sprintf(
		"%s/latest?access_key=%s&base=%s",
		XCHANGE_RATE_BASE_URL,
		XCHANGE_RATE_API_KEY,
		base,
	)
	if len(symbols) > 0 {
		url += "&symbols=" + strings.Join(symbols, ",")
	}

	client := http.Client{Timeout: 10*time.Second}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data ClientApiRes
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	if !data.Success {
		return nil, fmt.Errorf("failed to fetch rates")
	}

	return data.Rates, nil
}