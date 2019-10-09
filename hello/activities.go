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

func stravaAPIGetCall(url string, params map[string]uint64, accessToken string) map[string]interface{} {
	log.Println("API call URL ", url)

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
	log.Println("API call URL ", request.URL.String())

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

	//data, err := json.Unmarshal(body)
	log.Printf("activity body: %s\n", body) // dumps the whole resonse

	//var parsed map[string]interface{}
	//err = json.Unmarshal(body, &parsed)
	//if err != nil {
	//	//	return auth, err
	//	log.Fatalln("Unable to read the activities get: ", err)
	//}

	//return parsed
	return nil
}

func getActivities(startDate uint64, endDate uint64, accessToken string) {
	activitiyListParams := map[string]uint64{
		"before": endDate,
		"after":  startDate,
		//"page": 1,
		"per_page": 200,
	}
	jsonResponse := stravaAPIGetCall(stravaListActivitiesPath, activitiyListParams, accessToken)

	log.Println(jsonResponse)
}

func outputActivityStartStop(id uint64, accessToken string) {
	url := stravaGetActivityPath + strconv.FormatUint(id, 10)
	log.Println("activity URL ", url)

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)

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

	//data, err := json.Unmarshal(body)
	//log.Printf("activity body: %s\n", body) // dumps the whole resonse

	var parsed map[string]interface{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		//	return auth, err
		log.Fatalln("Unable to read the activities get: ", err)
	}

	// update the token struct with the new refresh token from Strava OAuth request
	name := parsed["name"].(string)

	log.Println("activity name: ", name)
}
