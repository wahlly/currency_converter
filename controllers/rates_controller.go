package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wahlly/currency_converter/services"
)

type RatesController struct {
	Cache *services.RatesCache
}

func (rc *RatesController) FetchRates(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	symbols := strings.ToUpper(query.Get("symbols"))

	rates, err := services.FetchRates(symbols)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(rates)
}