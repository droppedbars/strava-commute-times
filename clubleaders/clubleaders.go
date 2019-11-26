package main

import (
	"fmt"
	"strconv"

	"github.com/droppedbars/strava-commute-times/logger"
	"github.com/droppedbars/strava-commute-times/stravahelpers"
)

func getClubMembers(clubID uint64) {
	i := 1
	params := map[string]uint64{
		"page":     uint64(i),
		"per_page": 200,
	}
	path := fmt.Sprintf(stravahelpers.StravaListClubMembersPath, clubID)
	arrayJSONResponse, err := stravahelpers.StravaAPIGetArray(path, params)
	if err != nil {
		logger.ERROR.Fatal(err)
	}
	//if len(arrayJSONResponse) == 0 { // empty response, so no more data
	//	break
	//}
	logger.DEBUG.Println("Page: ", i)
	logger.TRACE.Println("API call response page: "+strconv.Itoa(i)+": ", arrayJSONResponse)
	logger.TRACE.Println("\n\nOne response: ", arrayJSONResponse[1])
}

func getActivity() {
	i := 1
	params := map[string]uint64{
		//	"page":     uint64(i),
		//	"per_page": 200,
	}
	//path := fmt.Sprintf(stravahelpers.StravaListClubMembersPath, clubID)
	path := stravahelpers.StravaGetActivityPath + "2891808686"
	arrayJSONResponse, err := stravahelpers.StravaAPIGetJSON(path, params)
	if err != nil {
		logger.ERROR.Fatal(err)
	}
	//if len(arrayJSONResponse) == 0 { // empty response, so no more data
	//	break
	//}
	logger.DEBUG.Println("Page: ", i)
	//logger.TRACE.Println("API call response page: "+strconv.Itoa(i)+": ", arrayJSONResponse)
	logger.TRACE.Println("\n\nOne response: ", arrayJSONResponse)
}

func getAthleteStats() {
	i := 1
	params := map[string]uint64{
		//	"page":     uint64(i),
		//	"per_page": 200,
	}
	//path := fmt.Sprintf(stravahelpers.StravaListClubMembersPath, clubID)
	path := "https://www.strava.com/api/v3/athletes/541441/stats"
	arrayJSONResponse, err := stravahelpers.StravaAPIGetJSON(path, params)
	if err != nil {
		logger.ERROR.Fatal(err)
	}
	//if len(arrayJSONResponse) == 0 { // empty response, so no more data
	//	break
	//}
	logger.DEBUG.Println("Page: ", i)
	//logger.TRACE.Println("API call response page: "+strconv.Itoa(i)+": ", arrayJSONResponse)
	logger.TRACE.Println("\n\nOne response: ", arrayJSONResponse)
}

// main execution function.
func main() {
	logger.SetLogging(false, logger.TraceLevel)

	//year := 2019
	clubID := uint64(465748) // VicHillVelo

	err := stravahelpers.StravaAuthenticate()
	if err != nil {
		logger.ERROR.Fatalln(err)
	}

	getClubMembers(clubID)
	//getActivity()
	getAthleteStats()
}
