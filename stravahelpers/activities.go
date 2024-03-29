// Package stravahelpers is a package of helper functions for interacting with the Strava APIs.
// It makes use of github.com/droppedbars/strava-commute-times/logger for logging. The main
// application should set up logger before making use of stravahelpers. The logging is useful to
// show what messaging is going back and forth, and it makes use of INFO, DEBUG and TRACE log levels.
// Errors are returned as errors to the main application.
// The first function to call is StravaAuthenticate. This will attempt to authenticate using OAuth
// with Strava and populate the authentication tokens. If it is not run, the tokens will be empty
// and attempts to interact with Strava APIs will fail.
package stravahelpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/droppedbars/strava-commute-times/logger"
)

const stravaBasePath = "https://www.strava.com/api/v3/"

// StravaGetActivityPath is the URL for strava's GET activities
const StravaGetActivityPath = stravaBasePath + "activities/"

// StravaListActivitiesPath is the URL to GET the list of all of an athletes activities
const StravaListActivitiesPath = stravaBasePath + "athlete/activities/"

// StravaListClubMembersPath is the URL to GET the list of members for a club. The club ID must
// be provided, url := sprintf(StravaListClubMembersPath, "12345")
const StravaListClubMembersPath = stravaBasePath + "clubs/%d/members"

// StravaListClubActivitiesPath is the URL to GET the list of activities for a club. The club ID must
// be provided, url := sprintf(StravaListClubActivitiesPath, "12345")
const StravaListClubActivitiesPath = stravaBasePath + "clubs/%d/activities"

// StravaGetAtheleteStatsPath is the URL to GET the athlete's stats. The athlete ID must
// be provided, url := sprintf(StravaGetAtheleteStatsPath, "12345")
const StravaGetAtheleteStatsPath = stravaBasePath + "athletes/%d/stats"

// StravaAPIGetResponse makes a call to a Strava GET API.
//  url is the URL to the API
//  params is a map of key/value parameters to provide to the API
// TODO: params should handle parameters that are not uint64
func StravaAPIGetResponse(url string, params map[string]uint64) ([]byte, error) {
	logger.DEBUG.Println("Base API call URL ", url)
	accessToken := auth.AccessToken

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
	logger.INFO.Println("Full API call URL ", request.URL.String())

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Unable to access the activities get: %s", err)
	}
	defer resp.Body.Close()

	// ensure a proper response. Anything other than 200 is an error (user or server)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status not 200: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %s", err)
	}

	logger.TRACE.Printf("activity body: %s\n", body) // dumps the whole resonse
	return body, nil
}

// StravaAPIGetJSON returns the Strava API response which is expected to be a json result.
// TODO: params should handle parameters that are not uint64
func StravaAPIGetJSON(url string, params map[string]uint64) (map[string]interface{}, error) {
	rawResponse, err := StravaAPIGetResponse(url, params)
	if err != nil {
		return nil, err
	}

	var parsed map[string]interface{}
	err = json.Unmarshal(rawResponse, &parsed)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the response: %s", err)
	}

	return parsed, nil
}

// StravaAPIGetArray returns the Strava API response which is expected to be an array of json results.
//  url is the API url, params is the key/value map of paramters, accessToken is the Strava access token.
// TODO: params should handle parameters that are not uint64
// TODO: need to ensure it gracefully handles API calls that do not return arrays of json
func StravaAPIGetArray(url string, params map[string]uint64) ([]map[string]interface{}, error) {
	rawResponse, err := StravaAPIGetResponse(url, params)
	if err != nil {
		return nil, err
	}

	var parsed []map[string]interface{}
	err = json.Unmarshal(rawResponse, &parsed)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the response: %s", err)
	}

	return parsed, nil
}
