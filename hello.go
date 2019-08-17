// An application that if provided a Strava client id, client secret, and refresh token via
// the command line will use the refresh token to get a new refresh token, and save that and
// the client id, ad cliet secret to a json file. If the client id, secret and refresh token
// are not provided in the command line, then the application will attempt to read them
// from the json file.
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const tokenJSONFileName = "./tokens.json"
const stravaOAuthPath = "https://www.strava.com/oauth/token"

// struct that defines the necessary authorization tokens for Strava
type tokens struct {
	// note: struct members must be capitized or they're not visible outside the struct
	ClientID     int
	ClientSecret string
	RefreshToken string
}

// Loads the Strava client id, secret and refresh token either from command line flags, or the json file
// and return them in a tokens struct.
func loadSecrets() (tokens, error) {
	var obj tokens

	// if arguments were supplied, those are used to look for the Strava secrets
	if len(os.Args) > 1 {
		// command line flags
		var clientID = flag.Int("clientID", -1, "Client ID found at https://www.strava.com/settings/api")
		var clientSecret = flag.String("clientSecret", "", "Client Secret found at https://www.strava.com/settings/api")
		var refreshToken = flag.String("refreshToken", "", "Refresh token provided by Strava")

		flag.Parse()

		obj.ClientID = *clientID
		obj.ClientSecret = *clientSecret
		obj.RefreshToken = *refreshToken

		data, err := json.Marshal(obj)
		if err != nil {
			return obj, err
		}

		log.Println("write to json: ", data)

		ioutil.WriteFile(tokenJSONFileName, data, 0644)
		if err != nil {
			return obj, err
		}

	} else { // read the values from the json file instead
		data, err := ioutil.ReadFile(tokenJSONFileName)
		if err != nil {
			return obj, err
		}

		log.Println("data: ", data)

		// unmarshall it
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return obj, err
		}

		log.Println("json: ", obj)
	}

	log.Println("clientId: ", obj.ClientID)
	log.Println("clientSecret: ", obj.ClientSecret)
	log.Println("refreshToken: ", obj.RefreshToken)

	return obj, nil
}

// receives a token struct and stores them in a json file.
func storeSecrets(obj tokens) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	log.Println("write to json: ", data)

	ioutil.WriteFile(tokenJSONFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// the main execution function.
func main() {
	obj, err := loadSecrets()
	if err != nil {
		log.Fatalln(err)
	}

	// create the POST body for Strava OAuth
	formData := url.Values{
		"client_id":     {strconv.Itoa(obj.ClientID)},
		"client_secret": {obj.ClientSecret},
		"refresh_token": {obj.RefreshToken},
		"grant_type":    {"refresh_token"},
	}

	// execute an HTTP POST to Strava OAuth to get new tokens
	resp, err := http.PostForm(stravaOAuthPath, formData)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// read and parse out the auth tokens from Strava
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	// ensure a proper response. Anything other than 200 is an error (user or server)
	if resp.StatusCode != 200 {
		log.Fatalf("HTTP Status not 200: %d - %s\n", resp.StatusCode, resp.Status)
	}
	log.Printf("http response: %s\n", string(body))

	var parsed map[string]interface{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		log.Fatalln(err)
	}

	// update the token struct with the new refresh token from Strava OAuth request
	obj.RefreshToken = parsed["refresh_token"].(string)

	log.Println("parsed body: ", obj)

	err = storeSecrets(obj)
	if err != nil {
		log.Fatalln(err)
	}
}
