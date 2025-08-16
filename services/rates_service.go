package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var _ = godotenv.Load()

type RatesService struct {
	Redis	*redis.Client
}

func NewRatesService(redisClient *redis.Client) *RatesService {
	return &RatesService{Redis: redisClient}
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

func (rs *RatesService) GetRates(ctx context.Context) (map[string]float64, error) {
	val, err := rs.Redis.Get(ctx, "rates").Result()
	if err != nil {
		return nil, err
	}

	var rates map[string]float64
	if err := json.Unmarshal([]byte(val), &rates); err != nil {
		return nil, err
	}
	return rates, nil
}

func HandleRateUpdates(ctx context.Context, redisClient *redis.Client) {
	go func ()  {
		for {
			rates, err := FetchRates()
			if err == nil {
				//round each rate value to 2 dp
				for currency, rate := range rates{
					rates[currency] = math.Round(rate * 100)/100
				}
				// cache.SetRate(rates)
				rates_js, _ := json.Marshal(rates)
				if err := redisClient.Set(ctx, "rates", rates_js, 15*time.Minute).Err(); err != nil {
					log.Fatalf("unable to update rates: %s", err)
					//notify and send logs
				}
			}

			//run every 15 minutes, as the xchange rate site is updated every 15 minutes
			time.Sleep(time.Minute*15)
		}	
	}()
}