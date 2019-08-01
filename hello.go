package main

import (
	"flag"
	"fmt"
)

var clientID = flag.Int("clientID", -1, "Client ID found at https://www.strava.com/settings/api")
var clientSecret = flag.String("clientSecret", "", "Client Secret found at https://www.strava.com/settings/api")
var refreshToken = flag.String("refreshToken", "", "Refresh token provided by Strava")

func main() {
	flag.Parse()

	fmt.Println("clientId: ", *clientID)
	fmt.Println("clientSecret: ", *clientSecret)
	fmt.Println("refreshToken: ", *refreshToken)
}
