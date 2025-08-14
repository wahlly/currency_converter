package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/wahlly/currency_converter/services"
)

type RatesController struct {
	Cache *services.RatesCache
}

func (rc *RatesController) GetRates(res http.ResponseWriter, req *http.Request) {
	rates := rc.Cache.GetRates()

	_ = json.NewEncoder(res).Encode(rates)
}