// Package uses the Strava refresh token to obtain a new one.
// It reads the application secrets from ./api_client_secrets.json. Use ./api_client_secrets.template.json as a
// template. Fill it out and rename it to ./api_client_secrets.json.
// If ./tokens.json is available, it will read the refresh token from there and try to use it. If ./tokens.json
// is not available, then it will prompt the user to authentication the application via the web browser and
// obtain the refresh token from that.
package main

import (
	"log"
	"strconv"
	"time"
)

// outputActivityStartStop will take the id of a Strava activity and print out the times
// the activity started and stopped. It uses the accessToken to make a API call to Strava.
func outputActivityStartStop(id uint64, accessToken string) {
	JSONResponse := stravaAPIGetJSON(stravaGetActivityPath+strconv.FormatUint(id, 10), nil, accessToken)

	name := JSONResponse["name"].(string)

	log.Println("activity name: ", name)
}

// activityDistanceTotals takes an array of Strava activities (in the format returned by Strava)
// and returns the total distance traveled and the total commute distance traveled. Both are
// provided in kilometers.
func ridingDistanceTotals(allActivities []map[string]interface{}) (float64, float64) {
	commute := 0.0
	total := 0.0

	for j := 0; j < len(allActivities); j++ {
		//log.Println("Activity Name: ", allActivities[j]["name"].(string))
		if allActivities[j]["type"] == "Ride" || allActivities[j]["type"] == "EBikeRide" {
			distance := allActivities[j]["distance"].(float64) / 1000 // convert m to km
			total += distance
			if allActivities[j]["commute"].(bool) == true {
				commute += distance
			}
		}
	}
	return total, commute
}

// getActivities returns an array of Strava activities given a date range and accessToken.
// The dates are provided as time since epoc.
func getRidingActivities(startDate uint64, endDate uint64, accessToken string) []map[string]interface{} {
	var allActivities []map[string]interface{}

	for i := 1; ; i++ { // strava pages start at 1
		activitiyListParams := map[string]uint64{
			"before":   endDate,
			"after":    startDate,
			"page":     uint64(i),
			"per_page": 200,
		}
		arrayJSONResponse := stravaAPIGetArray(stravaListActivitiesPath, activitiyListParams, accessToken)
		if len(arrayJSONResponse) == 0 { // empty response, so no more data
			break
		}
		allActivities = append(allActivities, arrayJSONResponse...)
		//log.Println("Page: ", i)
		//log.Println("API call response page: "+strconv.Itoa(i)+": ", arrayJSONResponse)
		//log.Println("\n\nOne response: ", arrayJSONResponse[1])
	}
	return allActivities
}

// main execution function.
func main() {
	var auth tokens

	sec, err := loadSecrets()
	if err != nil {
		log.Fatalln(err)
	}
	refreshToken, accessToken, err := loadTokens(sec)
	auth.RefreshToken = refreshToken
	auth.AccessToken = accessToken

	if err != nil {
		log.Fatalln(err)
	}

	auth, err = stravaOAuthCall(sec, "refresh_token", auth)
	err = storeTokens(auth)
	if err != nil {
		log.Fatalln(err)
	}

	//outputActivityStartStop(2685947039, auth.AccessToken)

	//"before": 1577865599, //December 31, 2019 11:59:59 PM GMT-08:00
	//"after": 1546329601, //January 1, 2019 12:00:01 AM GMT-08:
	year := "2019" // later to be replaced with user input
	var startTime time.Time
	var endTime time.Time
	startTime, err = time.Parse(time.RFC3339, year+"-01-01T12:00:01-08:00")
	endTime, err = time.Parse(time.RFC3339, year+"-12-31T11:59:59-08:00")

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Start Time: ", startTime)
	log.Println("End Time: ", endTime)

	allActivities := getRidingActivities(uint64(startTime.Unix()), uint64(endTime.Unix()), auth.AccessToken)
	total, commute := ridingDistanceTotals(allActivities)
	log.Printf("Total Distance (km): %.1f\n", total)
	log.Printf("Total Commute (km): %.1f, %.1f%%\n", commute, (commute/total)*100)
	log.Printf("Total Pleasure (km): %.1f, %.1f%%\n", total-commute, ((total-commute)/total)*100)
}
