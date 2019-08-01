package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var clientID = flag.Int("clientID", -1, "Client ID found at https://www.strava.com/settings/api")
var clientSecret = flag.String("clientSecret", "", "Client Secret found at https://www.strava.com/settings/api")
var refreshToken = flag.String("refreshToken", "", "Refresh token provided by Strava")

func main() {
	flag.Parse()

	fmt.Println("clientId: ", *clientID)
	fmt.Println("clientSecret: ", *clientSecret)
	fmt.Println("refreshToken: ", *refreshToken)

	formData := url.Values{
		"client_id":     {strconv.Itoa(*clientID)},
		"client_secret": {*clientSecret},
		"refresh_token": {*refreshToken},
		"grant_type":    {"refresh_token"},
	}

	resp, err := http.PostForm("https://www.strava.com/oauth/token", formData)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", string(body))
}
