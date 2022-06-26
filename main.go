package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

// Global variables
var data string
var c int = 0
var f os.File

// A function to write stock data into .txt file
func GetStockData(symbol string, daysInterval int, interval datetime.Interval) {
	params := &chart.Params{
		Symbol: symbol,
		Start: &datetime.Datetime{
			Day:   time.Now().Day() - daysInterval,
			Month: int(time.Now().Month()),
			Year:  time.Now().Year(),
		},
		End: &datetime.Datetime{
			Day:   time.Now().Day(),
			Month: int(time.Now().Month()),
			Year:  2022,
		},
		Interval:   interval,
		IncludeExt: false,
	}

	var prevVol int = 0.0
	var prevPrice float64 = 0.0

	iter := chart.Get(params)
	for iter.Next() {
		data = fmt.Sprintf("Date: %v \t Volume: %v \t Close: %v\n", time.Unix(int64(iter.Bar().Timestamp), 0).Local(), iter.Bar().Volume, iter.Bar().Close)
		f, err2 := os.OpenFile(params.Symbol+".txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)

		if _, err2 = f.WriteString(data); err2 != nil {
			panic(err2)
		}
		temp, _ := strconv.ParseFloat(iter.Bar().Close.String(), 32)
		checkVolumePriceRelation(prevVol, iter.Bar().Volume, float32(prevPrice), float32(temp))
		prevVol = iter.Bar().Volume
		prevPrice = iter.Bar().Close.InexactFloat64()

	}

	defer f.Close()

	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
}

// A function to do some analysis on relationship between volumes and price changes.
func checkVolumePriceRelation(previousVol int, currVol int, previousPrice float32, currPrice float32) {
	changeInVol := currVol - previousVol
	changeInPrice := currPrice - previousPrice
	switch {
	case currVol == 0 || currPrice == 0:
	default:
		fmt.Printf("Change in volume in %%: %v\n", ((float32(changeInVol) / float32(currVol)) * 100))
		fmt.Printf("Change in price in %%: %v\n", (changeInPrice/currPrice)*100)
	}

}

func main() {
	GetStockData("AAPL", 10, datetime.OneHour)
}
