package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wahlly/currency_converter/cmd/api/services"
)

type RatesController struct {
	RateService *services.RatesService
}

func (rc *RatesController) GetRates(res http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	rates, err := rc.RateService.GetRates(ctx)
	if err != nil{
		fmt.Println(err)
	}

	_ = json.NewEncoder(res).Encode(rates)
}