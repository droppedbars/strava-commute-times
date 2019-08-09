package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type tokens struct { // not, struct members must be capitized, or they're not visible outside the struct
	ClientID     int
	ClientSecret string
	RefreshToken string
}

var clientID = flag.Int("clientID", -1, "Client ID found at https://www.strava.com/settings/api")
var clientSecret = flag.String("clientSecret", "", "Client Secret found at https://www.strava.com/settings/api")
var refreshToken = flag.String("refreshToken", "", "Refresh token provided by Strava")

func main() {
	if len(os.Args) > 1 { // if arguments provided we'll use those to create the tokens file
		flag.Parse()

		var obj tokens
		obj.ClientID = *clientID
		obj.ClientSecret = *clientSecret
		obj.RefreshToken = *refreshToken

		data, err := json.Marshal(obj)
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println("write to json: ", data)

		ioutil.WriteFile("./tokens.json", data, 0644)
		if err != nil {
			log.Fatal(err)
		}

	} else { // read the values from the json file instead
		data, err := ioutil.ReadFile("./tokens.json")
		if err != nil {
			fmt.Print(err)
		}

		fmt.Println("data: ", data)

		// json data
		var obj tokens

		// unmarshall it
		err = json.Unmarshal(data, &obj)
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println("json: ", obj)

		*clientID = obj.ClientID
		*clientSecret = obj.ClientSecret
		*refreshToken = obj.RefreshToken

	}

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
