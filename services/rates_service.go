package services

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

type RatesCache struct {
	mutex	sync.RWMutex
	Rates	map[string]float64
}

func NewRatesCache(ttl time.Duration) *RatesCache {
	return &RatesCache{Rates: make(map[string]float64)}
}

var (
	XCHANGE_RATE_BASE_URL = os.Getenv("XCHANGE_RATE_BASE_URL")
	XCHANGE_RATE_API_KEY = os.Getenv("XCHANGE_RATE_API_KEY")
)

type ClientApiRes struct {
	Success bool `json:"success"`
	Timestamp any `json:"timestamp"`
	Base string `json:"base"`
	Date string `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func FetchRates() (map[string]float64, error) {
	url := fmt.Sprintf(
		"%s/latest?access_key=%s&symbols=%s",
		XCHANGE_RATE_BASE_URL,
		XCHANGE_RATE_API_KEY,
		"NGN,USD,EUR,AUD,CAD,XAU,XAG,AED",
	)

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
	// fmt.Println(data)
	if !data.Success {
		return nil, fmt.Errorf("failed to fetch rates")
	}

	return data.Rates, nil
}

func (cache *RatesCache) GetRates() (map[string]float64) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	return cache.Rates
}

func (cache *RatesCache) SetRate(rates map[string]float64) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.Rates = rates
}

func HandleRateUpdates(cache *RatesCache) {
	go func ()  {
		for {
			rates, err := FetchRates()
			if err == nil {
				//round each rate value to 2 dp
				for currency, rate := range rates{
					rates[currency] = math.Round(rate * 100)/100
				}
				cache.SetRate(rates)
			}

			//run every 15 minutes, as the xchange rate site is updated every 15 minutes
			time.Sleep(time.Minute*15)
		}	
	}()
}