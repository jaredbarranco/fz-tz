package geoapify

import (
  "fmt"
  "encoding/json"
  "os"
  "net/http"
  "io"
  "github.com/google/go-querystring/query"
  // "github.com/jaredbarranco/fz-tz/internal/config"
)
type LocationBias struct {
	Proximity     *GeoPoint
	Circle        *GeoCircle
	Rectangle     *GeoRectangle
	CountryCodes  []string
}

type LocationFilter struct {
	Circle        *GeoCircle
	Rectangle     *GeoRectangle
	CountryCodes  []string
	PlaceID       string
}

type GeoPoint struct {
	Lon float64
	Lat float64
}

type GeoCircle struct {
	Center GeoPoint
	RadiusMeters int
}

type GeoRectangle struct {
	Lon1 float64
	Lat1 float64
	Lon2 float64
	Lat2 float64
}
type GeoapifyRequest struct {
	// APIKey       string  `url:"apiKey" validate:"required"`
	Text         string  `url:"text,omitempty"` // For unstructured requests
	Name         string  `url:"name,omitempty"`
	HouseNumber  string  `url:"housenumber,omitempty"`
	Street       string  `url:"street,omitempty"`
	Postcode     string  `url:"postcode,omitempty"`
	City         string  `url:"city,omitempty"`
	State        string  `url:"state,omitempty"`
	Country      string  `url:"country,omitempty"`
	Type         string  `url:"type,omitempty"`   // e.g., country, state, city, etc.
	Lang         string  `url:"lang,omitempty"`   // 2-letter ISO 639-1
	Limit        int     `url:"limit,omitempty"`  // default is 5
	Filter       string  `url:"filter,omitempty"` // e.g., rect:-122.5,37.9,-122.4,38.0
	Bias         string  `url:"bias,omitempty"`   // e.g., proximity:-122.5,37.9
	ResponseFormat string `url:"format,omitempty"` // json, xml, geojson
}

type geoApifyAuthenticatedRequest struct {
  GeoapifyRequest
  APIKey       string  `url:"apiKey" validate:"required"`
}

type GeoapifyResponse struct {
	Features []struct {
		Properties struct {
			Timezone struct {
				Name string `json:"name"`
			} `json:"timezone"`
		} `json:"properties"`
	} `json:"features"`
}

func generateGeoapifyQuery(obj GeoapifyRequest) geoApifyAuthenticatedRequest {
  return geoApifyAuthenticatedRequest{
    GeoapifyRequest: obj,
    // APIKey: config.GeoApiKey,
    APIKey: os.Getenv("GEO_API_KEY"),

  }
}

func toQueryParams(obj interface{}) (string, error) {
	v, err := query.Values(obj)
	if err != nil {
		return "", err
	}
	return v.Encode(), nil
}

func GetLocationGuesses(obj GeoapifyRequest) (GeoapifyResponse, error) {
	url := "https://api.geoapify.com/v1/geocode/search?%s"
	queryParamFormattedStr, err := toQueryParams(generateGeoapifyQuery(obj))
	if err != nil {
		return GeoapifyResponse{}, err
	}

	fullUrl := fmt.Sprintf(url, queryParamFormattedStr)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return GeoapifyResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return GeoapifyResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GeoapifyResponse{}, err
	}

	var parsed GeoapifyResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return GeoapifyResponse{}, err
	}

	return parsed, nil
}
