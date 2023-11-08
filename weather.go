package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

type Location struct {
	Latitude                  float64 `json:"latitude"`
	LookupSource              string  `json:"lookupSource"`
	Longitude                 float64 `json:"longitude"`
	LocalityLanguageRequested string  `json:"localityLanguageRequested"`
	Continent                 string  `json:"continent"`
	ContinentCode             string  `json:"continentCode"`
	CountryName               string  `json:"countryName"`
	CountryCode               string  `json:"countryCode"`
	PrincipalSubdivision      string  `json:"principalSubdivision"`
	PrincipalSubdivisionCode  string  `json:"principalSubdivisionCode"`
	City                      string  `json:"city"`
	Locality                  string  `json:"locality"`
	Postcode                  string  `json:"postcode"`
	PlusCode                  string  `json:"plusCode"`
	LocalityInfo              struct {
		Administrative []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			IsoName     string `json:"isoName,omitempty"`
			Order       int    `json:"order"`
			AdminLevel  int    `json:"adminLevel"`
			IsoCode     string `json:"isoCode,omitempty"`
			WikidataID  string `json:"wikidataId"`
			GeonameID   int    `json:"geonameId"`
		} `json:"administrative"`
		Informative []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			IsoName     string `json:"isoName,omitempty"`
			Order       int    `json:"order"`
			IsoCode     string `json:"isoCode,omitempty"`
			WikidataID  string `json:"wikidataId,omitempty"`
			GeonameID   int    `json:"geonameId,omitempty"`
		} `json:"informative"`
	} `json:"localityInfo"`
}

func (loc *Location) initLocation() {
	urlLocation := "https://api.bigdatacloud.net/data/reverse-geocode-client"
	reqLoc, _ := http.NewRequest("GET", urlLocation, nil)
	reqLoc.Header.Add("Accept", "application/json")
	resLoc, err := http.DefaultClient.Do(reqLoc)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resLoc.Body.Close()

	bodyLoc, err := io.ReadAll(resLoc.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(bodyLoc, loc); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON (Location)")
		return
	}
}

type LocalWeather struct {
	localDateEu          [7]string  // The date in the EU format DD-MM-YYYY for 7 days
	weather              [7]string  // The weather for 7 days
	location             Location   // The Location from IP Adress
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	GenerationtimeMs     float64    `json:"generationtime_ms"`
	UtcOffsetSeconds     int        `json:"utc_offset_seconds"`
	Timezone             string     `json:"timezone"`
	TimezoneAbbreviation string     `json:"timezone_abbreviation"`
	Elevation            float64    `json:"elevation"`
	DailyUnits           DailyUnits `json:"daily_units"`
	Daily                Daily      `json:"daily"`
}
type DailyUnits struct {
	Time                         string `json:"time"`
	Temperature2MMin             string `json:"temperature_2m_min"`
	Temperature2MMax             string `json:"temperature_2m_max"`
	PrecipitationProbabilityMean string `json:"precipitation_probability_mean"`
	WeatherCode                  string `json:"weather_code"`
}
type Daily struct {
	Time                         []string  `json:"time"`
	Temperature2MMin             []float64 `json:"temperature_2m_min"`
	Temperature2MMax             []float64 `json:"temperature_2m_max"`
	PrecipitationProbabilityMean []int     `json:"precipitation_probability_mean"`
	WeatherCode                  []int     `json:"weather_code"`
}

// Initialize local weather from IP Adress
func (lw *LocalWeather) InitLocalWeather() {
	lw.location.initLocation()
	urlWeather := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?&latitude=%f&longitude=%f&localityLanguage=default&daily=temperature_2m_min,temperature_2m_max,precipitation_probability_mean,weather_code", lw.location.Latitude, lw.location.Longitude)

	req, _ := http.NewRequest("GET", urlWeather, nil)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(body, lw); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON (LocalWeather)")
		return
	}

	// Get weather and date in EU format for 7 days
	for day := 0; day < 7; day++ {
		lw.getWeatherString(day)
		lw.getLocalDateEu(day)
	}
}

func (lw *LocalWeather) getWeatherString(day int) {
	switch lw.Daily.WeatherCode[day] {
	case 0:
		lw.weather[day] = "klarer Himmel"
	case 1, 2, 3:
		lw.weather[day] = "hauptsächlich klar, teilweise bewölkt"
	case 45, 48:
		lw.weather[day] = "Nebel und sich ablagernder Raureifnebel"
	case 51, 53, 55:
		lw.weather[day] = "Nieselregen: Leichte, mäßige und dichte Intensität"
	case 56, 57:
		lw.weather[day] = "Gefrierender Nieselregen: Leichte und dichte Intensität"
	case 61, 63, 65:
		lw.weather[day] = "Regen: Leichte, mäßige und starke Intensität"
	case 66, 67:
		lw.weather[day] = "gefrierender Regen: Leichte und dichte Intensität"
	case 71, 73, 75:
		lw.weather[day] = "Schnee: Leichte, mäßige und starke Intensität"
	case 77:
		lw.weather[day] = "Schneekörner"
	case 80, 81, 82:
		lw.weather[day] = "Regenschauer: Leicht, mäßig und heftig"
	case 85, 86:
		lw.weather[day] = "leichte und heftige Schneeschauer"
	case 95:
		lw.weather[day] = "Gewitter: Leicht oder mäßig"
	case 96, 99:
		lw.weather[day] = "Gewitter mit leichtem und schwerem Hagel"
	default:
		lw.weather[day] = "unbekannt"
	}
}

func (lw *LocalWeather) getLocalDateEu(day int) {
	localTime := lw.Daily.Time[day]
	localTs := strings.Split(localTime, "-")
	slices.Reverse(localTs)
	lw.localDateEu[day] = fmt.Sprintf("%s.%s.%s", localTs[0], localTs[1], localTs[2])
}
