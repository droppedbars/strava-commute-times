package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const stravaBasePath = "https://www.strava.com/api/v3/"
const stravaGetActivityPath = stravaBasePath + "activities/"
const stravaListActivitiesPath = stravaBasePath + "athlete/activities/"

func stravaAPIGetResponse(url string, params map[string]uint64, accessToken string) []byte {
	log.Println("Base API call URL ", url)

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, strconv.FormatUint(value, 10))
	}
	request.URL.RawQuery = query.Encode()
	log.Println("Full API call URL ", request.URL.String())

	resp, err := client.Do(request)
	if err != nil {
		//return auth, err
		log.Fatalln("Unable to access the activities get: ", err)
	}
	defer resp.Body.Close()

	// ensure a proper response. Anything other than 200 is an error (user or server)
	if resp.StatusCode != 200 {
		log.Fatalf("HTTP Status not 200: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading body: ", err)
	}

	//log.Printf("activity body: %s\n", body) // dumps the whole resonse
	return body
}

func stravaAPIGetJSON(url string, params map[string]uint64, accessToken string) map[string]interface{} {
	rawResponse := stravaAPIGetResponse(url, params, accessToken)

	var parsed map[string]interface{}
	err := json.Unmarshal(rawResponse, &parsed)
	if err != nil {
		log.Fatalln("Unable to parse the response: ", err)
	}

	return parsed
}

func stravaAPIGetArray(url string, params map[string]uint64, accessToken string) []map[string]interface{} {
	rawResponse := stravaAPIGetResponse(url, params, accessToken)

	var parsed []map[string]interface{}
	err := json.Unmarshal(rawResponse, &parsed)
	if err != nil {
		log.Fatalln("Unable to parse the response: ", err)
	}

	return parsed
}
