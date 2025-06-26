package localizeTz

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/jaredbarranco/fz-tz/internal/geoapify"
)

func LocalizeTzHandler(w http.ResponseWriter, r *http.Request){

	var req geoapify.GeoapifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
	}

	var geoRes, err = geoapify.GetLocationGuesses(req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(geoRes); err != nil {
		fmt.Println(err)
		http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
		return
	}
}
