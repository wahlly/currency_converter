package routes

import (
	"net/http"

	"github.com/wahlly/currency_converter/controllers"
)

func RegisterRoutes(mux *http.ServeMux, rc *controllers.RatesController) {
	mux.HandleFunc("GET /rates", rc.GetRates)
}