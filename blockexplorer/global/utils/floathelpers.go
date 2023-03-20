package utils

import (
	"fmt"
	"math/rand"
	"strconv"
)

func RangeRandomF64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RangeRandomBTC(min, max float64) float64 {
	return BtcRoundFloat(RangeRandomF64(min, max))
}

func BtcRoundFloat(f float64) float64 {
	fStr := fmt.Sprintf("%.8f", f)
	f, _ = strconv.ParseFloat(fStr, 64)
	return f
}
