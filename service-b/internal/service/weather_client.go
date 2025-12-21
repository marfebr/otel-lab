package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

// GetWeatherAPICallWithURL busca clima por cidade usando WeatherAPI (padrão cloud-run)
func GetWeatherAPICallWithURL(city string, baseURL string) (ResponseTemps, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	api := viper.GetString("WEATHER_API")
	if api == "" {
		// MOCK: retorna dados fixos se não houver API key
		return ResponseTemps{
			TempC: 25.0,
			TempF: 77.0,
			TempK: 298.15,
		}, nil
	}
	client := http.Client{Transport: tr}

	encodedCity := url.QueryEscape(city)
	url := baseURL + "/current.json?key=" + api + "&q=" + encodedCity + "&aqi=no"

	resp, err := client.Get(url)
	if err != nil {
		return ResponseTemps{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResponseTemps{}, fmt.Errorf("weatherapi status: %s", resp.Status)
	}
	var weatherResponse WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return ResponseTemps{}, err
	}

	return ResponseTemps{
		TempC: weatherResponse.Current.TempC,
		TempF: weatherResponse.Current.TempF,
		TempK: (weatherResponse.Current.TempC + 273.15),
	}, nil
}

// GetWeatherAPICall mantém compatibilidade
func GetWeatherAPICall(city string) (ResponseTemps, error) {
	return GetWeatherAPICallWithURL(city, "https://api.weatherapi.com/v1")
}
