/*
The Ichimoku Cloud is based on five lines, each with its formula:

Tenkan-sen (Conversion Line): (9-period high + 9-period low) / 2
Kijun-sen (Base Line): (26-period high + 26-period low) / 2
Senkou Span A (Leading Span A): (Conversion Line + Base Line) / 2, plotted 26 periods ahead
Senkou Span B (Leading Span B): (52-period high + 52-period low) / 2, plotted 26 periods ahead
Chikou Span (Lagging Span): Current closing price, plotted 26 periods behind
*/

package indicators

import "github.com/sudeepbatra/alpha-hft/common"

type IchimokuCloud struct {
	Tenkan  []float64
	Kijun   []float64
	SenkouA []float64
	SenkouB []float64
	Chikou  []float64
}

// The periods are 9, 26, 52 for Tenkan, Kijun, Senkou Span B respectively.
func CalculateIchimokuCloud(high, low, close []float64, tenkanPeriod, kijunPeriod, senkouBPeriod int) *IchimokuCloud {
	ic := &IchimokuCloud{}

	length := len(close)
	ic.Tenkan = make([]float64, length)
	ic.Kijun = make([]float64, length)
	ic.SenkouA = make([]float64, length)
	ic.SenkouB = make([]float64, length)
	ic.Chikou = make([]float64, length)

	for i := senkouBPeriod; i <= len(high); i++ {
		tenkanSen := (common.MaxFloat64(high[i-tenkanPeriod:i]) +
			common.MinFloat64(low[i-tenkanPeriod:i])) / 2

		kijunSen := (common.MaxFloat64(high[i-kijunPeriod-1:i]) +
			common.MinFloat64(low[i-kijunPeriod:i])) / 2

		senkouSpanA := (tenkanSen + kijunSen) / 2

		senkouSpanB := (common.MaxFloat64(high[i-senkouBPeriod:i]) +
			common.MinFloat64(low[i-senkouBPeriod:i])) / 2

		chikouSpan := close[i-1]

		ic.Tenkan = append(ic.Tenkan, tenkanSen)
		ic.Kijun = append(ic.Kijun, kijunSen)
		ic.SenkouA = append(ic.SenkouA, senkouSpanA)
		ic.SenkouB = append(ic.SenkouB, senkouSpanB)
		ic.Chikou = append(ic.Chikou, chikouSpan)
	}

	return ic
}
