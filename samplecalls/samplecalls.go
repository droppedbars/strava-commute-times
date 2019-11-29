// package main makes a series of sample calls to Strava's APIS to show how to use some of them and
// what data results come back.
package main

import (
	"fmt"
	"strconv"

	"github.com/droppedbars/strava-commute-times/logger"
	"github.com/droppedbars/strava-commute-times/stravahelpers"
)

// getClubMembers iterates through the list of club members and prints out the API response
func getClubMembers(clubID uint64) {
	for i := 1; ; i++ {
		params := map[string]uint64{
			"page":     uint64(i),
			"per_page": 200,
		}
		path := fmt.Sprintf(stravahelpers.StravaListClubMembersPath, clubID)
		arrayJSONResponse, err := stravahelpers.StravaAPIGetArray(path, params)
		if err != nil {
			logger.ERROR.Println(err)
		}
		if len(arrayJSONResponse) == 0 { // empty response, so no more data
			break
		}
		logger.DEBUG.Println("Page: ", i)
		logger.TRACE.Println("API call response page: "+strconv.Itoa(i)+": ", arrayJSONResponse)
		fmt.Println("First Club Member response: ", arrayJSONResponse[1])
	}
}

// getActivity will return the API response for a strava activity. If the activity does not belong to
// the authenticated user, or the API client doesn't have activities:read_all scope and the activity
// is private, then the API response will be a 404 - Not Found.
func getActivity(activityID uint64) {
	params := map[string]uint64{}
	path := stravahelpers.StravaGetActivityPath + strconv.FormatUint(activityID, 10)
	arrayJSONResponse, err := stravahelpers.StravaAPIGetJSON(path, params)
	if err != nil {
		logger.ERROR.Println(err)
	}
	fmt.Println("Activity response: ", arrayJSONResponse)
}

// getAthleteStats will return a Strava athlete's stats. If the athlete is not the
// the authenticated user then the API response will be a 403 - Forbidden.
func getAthleteStats(athleteID uint64) {
	params := map[string]uint64{}
	path := fmt.Sprintf(stravahelpers.StravaGetAtheleteStatsPath, athleteID)
	arrayJSONResponse, err := stravahelpers.StravaAPIGetJSON(path, params)
	if err != nil {
		logger.ERROR.Println(err)
	}
	fmt.Println("AthleteStats Response: ", arrayJSONResponse)
}

// main execution function.
func main() {
	logger.SetLogging(false, logger.TraceLevel)

	err := stravahelpers.StravaAuthenticate()
	if err != nil {
		logger.ERROR.Fatalln(err)
	}

	getClubMembers(465748)
	getActivity(12345)
	getAthleteStats(541441)
}
