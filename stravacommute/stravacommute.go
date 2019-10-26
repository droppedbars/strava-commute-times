// Package uses the Strava refresh token to obtain a new one.
// It reads the application secrets from ./api_client_secrets.json. Use ./api_client_secrets.template.json as a
// template. Fill it out and rename it to ./api_client_secrets.json.
// If ./tokens.json is available, it will read the refresh token from there and try to use it. If ./tokens.json
// is not available, then it will prompt the user to authentication the application via the web browser and
// obtain the refresh token from that.
package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

const annualCommuteKm = 5875 // assumes 25km/day, 5 days a week, 5 weeks of no riding per year
const hoursInYear = 24 * 365 // simplistic, ignores leap years
const epoch = 2009           // when strava started, so there should never be data before this

var flagYear1 = flag.Int("startYear", time.Now().Year(), "First year to run the commute numbers for. Defaults to current year.")
var flagYear2 = flag.Int("endYear", time.Now().Year(), "Last year to run the commute numbers for. Defaults to current year.")

// outputActivityStartStop will take the id of a Strava activity and print out the times
// the activity started and stopped. It uses the accessToken to make a API call to Strava.
func outputActivityStartStop(id uint64, accessToken string) {
	JSONResponse, err := stravaAPIGetJSON(stravaGetActivityPath+strconv.FormatUint(id, 10), nil, accessToken)
	if err != nil {
		ERROR.Fatal(err)
	}

	name := JSONResponse["name"].(string)

	fmt.Println("activity name: ", name)
}

// activityDistanceTotals takes an array of Strava activities (in the format returned by Strava)
// and returns the total distance traveled and the total commute distance traveled. Both are
// provided in kilometers.
func ridingDistanceTotals(allActivities []map[string]interface{}) (float64, float64) {
	commute := 0.0
	total := 0.0

	for j := 0; j < len(allActivities); j++ {
		TRACE.Println("Activity Name: ", allActivities[j]["name"].(string))
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
		arrayJSONResponse, err := stravaAPIGetArray(stravaListActivitiesPath, activitiyListParams, accessToken)
		if err != nil {
			ERROR.Fatal(err)
		}
		if len(arrayJSONResponse) == 0 { // empty response, so no more data
			break
		}
		allActivities = append(allActivities, arrayJSONResponse...)
		DEBUG.Println("Page: ", i)
		TRACE.Println("API call response page: "+strconv.Itoa(i)+": ", arrayJSONResponse)
		TRACE.Println("\n\nOne response: ", arrayJSONResponse[1])
	}
	return allActivities
}

// main execution function.
func main() {
	setLogging(true, debugLevel)

	flag.Parse()
	var year1 int
	var year2 int

	year1 = *flagYear1
	year2 = *flagYear2
	if year1 > year2 {
		year2 = *flagYear1
		year1 = *flagYear2
	}
	currentYear := time.Now().Year()
	if year2 > currentYear {
		year2 = currentYear
	}
	if year1 < epoch {
		year1 = epoch
	}

	var auth tokens

	sec, err := loadSecrets()
	if err != nil {
		ERROR.Fatalln(err)
	}
	refreshToken, accessToken, err := loadTokens(sec)
	auth.RefreshToken = refreshToken
	auth.AccessToken = accessToken

	if err != nil {
		ERROR.Fatalln(err)
	}

	auth, err = stravaOAuthCall(sec, "refresh_token", auth)
	err = storeTokens(auth)
	if err != nil {
		ERROR.Fatalln(err)
	}

	//outputActivityStartStop(2685947039, auth.AccessToken)

	for i := year1; i <= year2; i++ {
		year := strconv.Itoa(i)

		var startTime time.Time
		var endTime time.Time
		startTime, err = time.Parse(time.RFC3339, year+"-01-01T12:00:01-08:00")
		endTime, err = time.Parse(time.RFC3339, year+"-12-31T11:59:59-08:00")

		if err != nil {
			ERROR.Fatalln(err)
		}

		INFO.Println("Commute time range start: ", startTime)
		INFO.Println("Commute time range end: ", endTime)

		allActivities := getRidingActivities(uint64(startTime.Unix()), uint64(endTime.Unix()), auth.AccessToken)
		total, commute := ridingDistanceTotals(allActivities)
		percentageOfYear := 1.0
		fullYear := true
		if endTime.Unix() > time.Now().Unix() {
			percentageOfYear = time.Since(startTime).Hours() / hoursInYear
			fullYear = false
		}
		fmt.Println(year)
		fmt.Printf("Total Distance (km): %.1f\n", total)
		if !fullYear {
			fmt.Printf("  Estimated end of year distance (km): %.1f\n", total/percentageOfYear)
		}
		fmt.Printf("Total Commute (km): %.1f, %.1f%%\n", commute, (commute/total)*100)
		fmt.Printf("  Percentage of commute by bike: %.1f%%\n", (commute/annualCommuteKm)*100)
		if !fullYear {
			fmt.Printf("  Estimated percentage of commute by bike for year: %.1f%%\n", (commute/annualCommuteKm/percentageOfYear)*100)
		}
		fmt.Printf("Total Pleasure (km): %.1f, %.1f%%\n", total-commute, ((total-commute)/total)*100)
	}
}
