package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"owlharbour-api/pkg/util"
	"strings"
)

func StatusCheck(coord [2]float64, polygon [][2]float64) bool {
	x := coord[0] // Latitude
	y := coord[1] // Longitude

	isInside := false

	numVertices := len(polygon)
	for i, j := 0, numVertices-1; i < numVertices; i, j = i+1, i {
		xi := polygon[i][0] // Latitude
		yi := polygon[i][1] // Longitude
		xj := polygon[j][0] // Latitude
		yj := polygon[j][1] // Longitude

		intersect := ((yi > y) != (yj > y)) &&
			(x < (xj-xi)*(y-yi)/(yj-yi)+xi)

		if intersect {
			isInside = !isInside
		}
	}

	return isInside
}

type CoordinateResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []CoordinateData `json:"data"`
}

type CoordinateData struct {
	Coordinate []float64 `json:"coordinate"`
	IsWater    bool      `json:"is_water"`
}

func IsWater(latitude, longitude float64) (bool, error) {
	APIHost := util.GetEnv("CEKLAUT_HOST", "")
	url := fmt.Sprintf("http://%s/ceklaut", APIHost)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`[%f, %f]`, latitude, longitude))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	fmt.Println(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var response CoordinateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	fmt.Println(response)

	return response.Data[0].IsWater, nil
}

func IsWaterRapidAPI(latitude, longitude float64) (bool, error) {
	rapidAPIHost := util.GetEnv("RAPIDAPI_ISITWATER_HOST", "")
	rapidAPIKey := util.GetEnv("RAPIDAPI_KEY", "")

	url := fmt.Sprintf("https://%s/?latitude=%.6f&longitude=%.6f", rapidAPIHost, latitude, longitude)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("X-RapidAPI-Host", rapidAPIHost)
	req.Header.Set("X-RapidAPI-Key", rapidAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	var data struct {
		Water bool `json:"water"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, err
	}

	return data.Water, nil
}
