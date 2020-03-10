package mapping

import "strconv"

func intToBoolean(i int) (op bool) {
	if i == 1 {
		op = true
	}
	return
}

func stringToFloat(s string) (op float64) {
	op, _ = strconv.ParseFloat(s, 32)
	return
}

func stringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
