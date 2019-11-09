package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sort"
	"time"

	"github.com/vdobler/chart"
	"github.com/vdobler/chart/imgg"
)

func ticFormat(x float64) string {
	output := fmt.Sprintf("%d", int(x))
	return output
}

func graphResults(results map[int]stravaDistances) {
	// create the file to write to
	fileName := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	imgFile, err := os.Create("commute-" + fileName + ".png")
	if err != nil {
		ERROR.Panic(err)
	}
	defer imgFile.Close()

	// draw the base image and set its size
	i := image.NewRGBA(image.Rect(0, 0, 500, 500))             // RGBA image that is a 500x500 rectangle starting at 0,0
	bg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff}) // white background
	draw.Draw(i, i.Bounds(), bg, image.ZP, draw.Src)

	// set the chart style
	red := chart.Style{Symbol: 'o', LineColor: color.NRGBA{0xcc, 0x00, 0x00, 0xff},
		FillColor: color.NRGBA{0xff, 0x80, 0x80, 0xff},
		LineStyle: chart.SolidLine, LineWidth: 2}
	green := chart.Style{Symbol: '#', LineColor: color.NRGBA{0x00, 0xcc, 0x00, 0xff},
		FillColor: color.NRGBA{0x80, 0xff, 0x80, 0xff},
		LineStyle: chart.SolidLine, LineWidth: 2}

	var years []float64
	var commutes []float64
	var pleasure []float64
	firstYear := time.Now().Year()
	lastYear := epoch

	var keys []int
	for key := range results {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, v := range keys {
		years = append(years, float64(results[v].year))
		commutes = append(commutes, results[v].commute)
		pleasure = append(pleasure, results[v].pleasure)
		if firstYear > results[v].year {
			firstYear = results[v].year
		}
		if lastYear < results[v].year {
			lastYear = results[v].year
		}
	}

	// create the chart and add data
	barc := chart.BarChart{Title: "Strava Commutes and Pleasure Rides"}
	barc.Key.Hide = false
	barc.Key.Pos = "itl" // means to left
	barc.XRange.Fixed(float64(firstYear-1), float64(lastYear+1), 1)
	barc.XRange.Label = "Year"
	barc.XRange.TicSetting.Format = ticFormat
	barc.YRange.Label = "Distance (km)"
	barc.YRange.TicSetting.Format = ticFormat
	barc.ShowVal = 3 // show the value at top of the bar (above bar doesn't work for stacked graphs)

	barc.AddDataPair("Commutes", years, commutes, red)
	barc.AddDataPair("Pleasure", years, pleasure, green)

	// essentially create the image, and then plot it
	igr := imgg.AddTo(i, 0, 0, 500, 500, color.RGBA{0xff, 0xff, 0xff, 0xff}, nil, nil)
	barc.Stacked = true

	barc.Plot(igr)

	// encode it all as png format into the file
	err = png.Encode(imgFile, i)
	if err != nil {
		ERROR.Panicln(err)
	}
}
