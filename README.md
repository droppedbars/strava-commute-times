# strava-commute-times
Test project to try Strava APIs, and some different things in Go

## What it does
This application will go through your bike rides in Strava for the current year, or a range of years, and produce some basic stats: your total commute distance, pleasure distance, and total distance for the year. It will also produce the percentages and if it is looking at the current year, provide a forecast of anticipated total distance.
It also provides for how much of your commute distance each year is done by bicycle (currently hardcoded to a assumption of 5875 km per year)

## How to set up
1. Log into Strava and go to https://www.strava.com/settings/api to set up your own API Application
1. Get the Client ID and Secret
1. Copy *api_client_secrets.json.template* to *api_client_secrets.json*
1. Enter the Client ID and Secrete from step 2 into *api_client_secrets.json*
1. Run the application and follow the instructions: enter the provided URL into your web browser (you may need to log into strava) and click Authorize, copy and paste the resulting URL into the input in the application
1. The application will output your Strava ride information for the current year and generate a bar chart for the current year

## Notes for the Application
* The application keeps your Client ID and Secrete in an unencrypted json file. So be aware of that.
* The application keeps the access token and renewal token in an unencrypted json file. If you delete the file then you will need to run the set up steps again
* Flags -firstYear and -lastYear can be used to get the application to produce ride stats for a range of years and a bar chart for the range of years. The application will use the earlier year of the two as the first year and the later as the last year.
* The bar chart is saved in the same directory as the application and is named commute-YYYY-MM-DD.png and will overwrite a file if it already exists for that date.
* When calculating portions of a year, the application makes the simple assumption that there are 24x365 hours in the year. It makes no attempts to determine if its a leap year.
* A log file is written to stravacommute.log. It will always overwrite the file on start. Log level is set to debug and cannot be changed outside of code (ie, if you run the executable you cannot change it).
* The application will error if you try to provide a year before 2009 (the year of Strava's release).

## Why?
This application was written as an exercise to use Go. It is not the best way to interact with Strava (a web app that handles the authorization by the user would be more appropriate). It exercised a few skills, basic Go, multiple files in a package, Godoc, various data structures, objects, logging, commandline flags, using third-party libraries (for creating the graphs), and basic Goroutines.
