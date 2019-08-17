// An application that uses the Strava refresh token to obtain a new one.
// It reads the application secrets from ./api_client_secrets.json. Use ./api_client_secrets.template.json as a
// template. Fill it out and rename it to ./api_client_secrets.json.
// If ./tokens.json is available, it will read the refresh token from there and try to use it. If ./tokens.json
// is not available, then it will prompt the user to authentication the application via the web browser and
// obtain the refresh token from that.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const tokenJSONFileName = "./tokens.json"
const secretsJSONFileName = "./api_client_secrets.json"
const stravaOAuthPath = "https://www.strava.com/oauth/token"

// struct that contains the secrets for the API application
type secrets struct {
	ClientID     int
	ClientSecret string
}

// struct that defines the necessary authorization tokens for Strava
type tokens struct {
	RefreshToken string
}

func loadTokens(sec secrets) (tokens, error) {
	var obj tokens

	fileInfo, err := os.Stat(tokenJSONFileName)
	if (err == nil) && !(fileInfo.IsDir()) { // file exists and is not a directory, so read the auth tokens
		data, err := ioutil.ReadFile(tokenJSONFileName)
		if err != nil {
			return obj, err
		}

		log.Println("auth tokens raw data from file: ", data)

		// unmarshall it
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return obj, err
		}

		log.Println("auth tokens json: ", obj)
	} else { // the auth tokens are missing, so we need to get them from the user
		fmt.Printf("Enter the following into your web browser: \n")
		fmt.Printf("   http://www.strava.com/oauth/authorize?client_id=%d&response_type=code&redirect_uri=http://localhost/exchange_token&approval_prompt=force&scope=activity:read_all\n", sec.ClientID)

		fmt.Printf("\nCopy and paste the URL from the browser: ")
		// need to get them to enter the response URL
		var responseURL string
		fmt.Scanln(&responseURL)

		log.Println("User entered URL: ", responseUrl)
		// parse out the code
		// make a call to OAuth to authenticate and get the refresh token

		obj.RefreshToken = "temp value"

		return obj, fmt.Errorf("temp error, code path not completed")
	}

	log.Println("refreshToken: ", obj.RefreshToken)

	return obj, nil
}

// Loads the Strava client id, secret and refresh token either from command line flags, or the json file
// and return them in a tokens struct.
func loadSecrets() (secrets, error) {
	var obj secrets

	fileInfo, err := os.Stat(secretsJSONFileName)
	if err != nil || fileInfo.IsDir() {
		return obj, err
	}

	data, err := ioutil.ReadFile(secretsJSONFileName)
	if err != nil {
		return obj, err
	}

	log.Println("secrets raw data from file: ", data)

	// unmarshall it
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}

	log.Println("secrets json: ", obj)

	log.Println("clientId: ", obj.ClientID)
	log.Println("clientSecret: ", obj.ClientSecret)

	return obj, nil
}

// receives a token struct and stores them in a json file.
func storeTokens(auth tokens) error {
	data, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	log.Println("data to write to json: ", data)

	ioutil.WriteFile(tokenJSONFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// the main execution function.
func main() {
	sec, err := loadSecrets()
	if err != nil {
		log.Fatalln(err)
	}
	auth, err := loadTokens(sec)
	if err != nil {
		log.Fatalln(err)
	}

	// create the POST body for Strava OAuth
	formData := url.Values{
		"client_id":     {strconv.Itoa(sec.ClientID)},
		"client_secret": {sec.ClientSecret},
		"refresh_token": {auth.RefreshToken},
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
	log.Printf("OAuth http response: %s\n", string(body))

	var parsed map[string]interface{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		log.Fatalln(err)
	}

	// update the token struct with the new refresh token from Strava OAuth request
	auth.RefreshToken = parsed["refresh_token"].(string)

	log.Println("parsed body from OAuth call: ", auth)

	err = storeTokens(auth)
	if err != nil {
		log.Fatalln(err)
	}
}
