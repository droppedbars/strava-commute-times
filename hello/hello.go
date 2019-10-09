// Package uses the Strava refresh token to obtain a new one.
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
	AuthCode     string
	RefreshToken string
	AccessToken  string
}

// loadTokens loads the authentication tokens by trying the tokens.json first. If that fails, then it will
// provide the user with a URL to enter in the web browser, and ask for the resulting URL back,
// then parses out the authorization code and makes an OAuth call to get a valid refresh and
// access token.
// Returns refreshToken, accessToken, error
func loadTokens(sec secrets) (string, string, error) {
	var obj tokens
	var refreshToken string
	var accessToken string

	if sec.ClientID == 0 || sec.ClientSecret == "" {
		return refreshToken, accessToken, fmt.Errorf("loadTokens must have non-nil secrets")
	}

	fileInfo, err := os.Stat(tokenJSONFileName)
	if (err == nil) && !(fileInfo.IsDir()) { // file exists and is not a directory, so read the auth tokens
		data, err := ioutil.ReadFile(tokenJSONFileName)
		if err != nil {
			return refreshToken, accessToken, err
		}

		log.Println("auth tokens raw data from file: ", data)

		// unmarshall it
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return refreshToken, accessToken, err
		}
	} else { // the auth tokens are missing, so we need to get them from the user
		fmt.Printf("Enter the following into your web browser: \n")
		fmt.Printf("   http://www.strava.com/oauth/authorize?client_id=%d&response_type=code&redirect_uri=http://localhost/exchange_token&approval_prompt=force&scope=activity:read_all\n", sec.ClientID)

		fmt.Printf("\nCopy and paste the URL from the browser: ")
		// need to get them to enter the response URL
		var responseURLString string
		fmt.Scanln(&responseURLString)

		log.Println("User entered URL: ", responseURLString)

		// parse out the code from Strava
		responseURL, err := url.Parse(responseURLString)
		if err != nil {
			return refreshToken, accessToken, err
		}
		paramMap, err := url.ParseQuery(responseURL.RawQuery)
		if err != nil {
			return refreshToken, accessToken, err
		}
		code, codeExists := paramMap["code"]
		if !codeExists {
			return refreshToken, accessToken, fmt.Errorf("The code key could not be found in the supplied URL: %s", responseURLString)
		}
		obj.AuthCode = code[0]
		log.Println("Auth code is: ", obj.AuthCode)
		// make a call to OAuth to authenticate and get the refresh token

		obj, err = stravaOAuthCall(sec, "authorization_code", obj)
	}

	refreshToken = obj.RefreshToken
	accessToken = obj.AccessToken
	log.Println("refreshToken: ", refreshToken)
	log.Println("accessToken: ", accessToken)

	return refreshToken, accessToken, nil
}

// loadSecrets loads the Strava client id, secret and refresh token from the json file
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

// storeTokens receives a token struct and stores them in a json file.
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

// stravaOAuthCall calls the Strava's OAuth APIs. Grant type can be either "refresh_token"
// or it can be "authorization_code". The values will be set appropriately when
// making the call to Strava
func stravaOAuthCall(sec secrets, grantType string, auth tokens) (tokens, error) {
	var formData map[string][]string
	if grantType == "refresh_token" {
		formData = url.Values{
			"client_id":     {strconv.Itoa(sec.ClientID)},
			"client_secret": {sec.ClientSecret},
			"refresh_token": {auth.RefreshToken},
			"grant_type":    {"refresh_token"},
		}
	} else if grantType == "authorization_code" {
		formData = url.Values{
			"client_id":     {strconv.Itoa(sec.ClientID)},
			"client_secret": {sec.ClientSecret},
			"code":          {auth.AuthCode},
			"grant_type":    {"authorization_code"},
		}
	} else { // unexpected grant_type, so fail
		return auth, fmt.Errorf("unexpected grant_type")
	}

	// execute an HTTP POST to Strava OAuth to get new tokens
	resp, err := http.PostForm(stravaOAuthPath, formData)
	if err != nil {
		return auth, err
	}
	defer resp.Body.Close()

	// read and parse out the auth tokens from Strava
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return auth, err
	}

	// ensure a proper response. Anything other than 200 is an error (user or server)
	if resp.StatusCode != 200 {
		return auth, fmt.Errorf("HTTP Status not 200: %d - %s", resp.StatusCode, resp.Status)
	}
	log.Printf("OAuth http response: %s\n", string(body))

	var parsed map[string]interface{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return auth, err
	}

	// update the token struct with the new refresh token from Strava OAuth request
	auth.RefreshToken = parsed["refresh_token"].(string)
	auth.AccessToken = parsed["access_token"].(string)

	log.Println("parsed body from OAuth call: ", auth)

	return auth, nil
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

	outputActivityStartStop(2685947039, auth.AccessToken)

	//"before": 1577865599, //December 31, 2019 11:59:59 PM GMT-08:00
	//"after": 1546329601, //January 1, 2019 12:00:01 AM GMT-08:00
	getActivities(1546329601, 1577865599, auth.AccessToken)
}
