package main

import (
	"fmt"

	"github.com/markcheno/go-talib"
)

func main() {
	// Example data: close prices for a stock
	closePrices := []float64{100, 105, 110, 95, 120, 115, 130, 125, 140, 135}
	timePeriod := 5

	// Calculate linear regression slope
	slopes := talib.LinearRegSlope(closePrices, timePeriod)

	// Print the calculated slopes
	for i, slope := range slopes {
		fmt.Printf("Day %d: Slope = %.4f\n", i+1, slope)
	}
}
