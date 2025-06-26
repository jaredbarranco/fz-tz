package main

import (
	"os"
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
	// "context"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jaredbarranco/fz-tz/internal/config"
	"github.com/jaredbarranco/fz-tz/internal/tz"
	"github.com/jaredbarranco/fz-tz/internal/localizeTz"
	"github.com/jaredbarranco/fz-tz/internal/geoapify"
	"github.com/go-chi/docgen"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	flag.Parse()
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
	
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("pong"))
	})

	r.Get("/offset", func(w http.ResponseWriter, r *http.Request){
		timeIn:= r.URL.Query().Get("time")
		tzIn:= r.URL.Query().Get("tz")

		tzData:= tz.GetTzOffset(timeIn, tzIn)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tzData)

	})

	r.Get("/location/offset", func(w http.ResponseWriter, r *http.Request){
		timeIn:= r.URL.Query().Get("time")

		//location data

		locQuery:=  geoapify.GeoapifyRequest{
			City: r.URL.Query().Get("city"),
			State: r.URL.Query().Get("state"),
			Country: r.URL.Query().Get("country"),
		}
		locData, err:= geoapify.GetLocationGuesses(locQuery)
		if err != nil {
			fmt.Print(err)
			panic("Error getting location guesses")
		}
		
		tzIn:= locData.Features[0].Properties.Timezone.Name
		tzData:= tz.GetTzOffset(timeIn, tzIn)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tzData)
	})

	r.Post("/localize", localizeTz.LocalizeTzHandler)

	if *routes {
		doc := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/jaredbarranco/fz-tz/api",
			Intro:       "Welcome to the fz-tz generated docs.",
		})

		err := os.WriteFile("./api/api-docs.md", []byte(doc), 0644)
		if err != nil {
			fmt.Println("Error writing docs to file:", err)
		} else {
			fmt.Println("Docs written to api-docs.md")
		}

		return
	}
	http.ListenAndServe(":3333", r)
}
