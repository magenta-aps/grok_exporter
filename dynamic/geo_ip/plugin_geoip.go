package main

import (
	"encoding/json"
    "fmt"
    "reflect"
	"io/ioutil"
	"net/http"
	"text/template/parse"
	"github.com/fstab/grok_exporter/plugins"
)

func newGeoIPFunc() plugins.FunctionWithValidator {
	return plugins.FunctionWithValidator{
		Function: geoip,
		StaticValidator: validate,
	}
}

// JSON Unmarshal data-structure
// Derived from: https://freegeoip.app/json/
type GeoIP struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
    Timezone    string  `json:"time_zone"`
	Latitude    float32 `json:"latitude"`
	Lonitude    float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

func toString(a interface{}) (string, error) {
	val := reflect.ValueOf(a)
	switch val.Kind() {
        case reflect.String:
            return val.String(), nil
    }
	return "", fmt.Errorf("%T: unknown type", a)
}

func geoip(maybe_address interface{}) (string, error) {
    // Ensure the passed argument is a string
    // The argument can be a domain name or an IP address (v4 or v6)
    address, err := toString(maybe_address)
	if err != nil {
        return "", err
	}

    // Use freegeoip.app to get a JSON object for GeoIP
    response, err := http.Get("https://freegeoip.app/json/" + address)
	if err != nil {
        return "", err
	}
	defer response.Body.Close()

    // Read entire body
    body, err := ioutil.ReadAll(response.Body)
	if err != nil {
        return "", err
	}

	// Unmarshal the JSON body into our GeoIP struct
    var geo GeoIP
    err = json.Unmarshal(body, &geo)
	if err != nil {
        return "", err
	}

    return geo.CountryCode, nil
}

func validate(cmd *parse.CommandNode) error {
    return nil
}

func Generate() (string, plugins.FunctionWithValidator) { return "geoip", newGeoIPFunc() }
