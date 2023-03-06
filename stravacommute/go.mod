module github.com/droppedbars/strava-commute-times/stravacommute

go 1.18

require (
	github.com/droppedbars/strava-commute-times/logger v0.0.0-00010101000000-000000000000
	github.com/droppedbars/strava-commute-times/stravahelpers v0.0.0-00010101000000-000000000000
	github.com/vdobler/chart v1.0.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/llgcode/draw2d v0.0.0-20180825133448-f52c8a71aff0 // indirect
	golang.org/x/image v0.5.0 // indirect
)

replace github.com/droppedbars/strava-commute-times/logger => ../logger

replace github.com/droppedbars/strava-commute-times/stravahelpers => ../stravahelpers
