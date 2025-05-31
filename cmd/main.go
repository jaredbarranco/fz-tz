package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	// "context"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jaredbarranco/fz-tz/internal/config"
	"github.com/jaredbarranco/fz-tz/internal/tz"
)
func main() {
	r:= chi.NewRouter()
	cfg := config.LoadConfig()
	fmt.Printf("App: %s\n", cfg.AppName)
	fmt.Printf("Environment: %s\n", cfg.AppEnv)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	//Hello World Endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("hi"))
	})

	r.Get("/offset", func(w http.ResponseWriter, r *http.Request){
		timeIn:= r.URL.Query().Get("time")
		tzIn:= r.URL.Query().Get("tz")

		tzData:= tz.GetTzOffset(timeIn, tzIn)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tzData)

	})
	http.ListenAndServe(":3333", r)
}
