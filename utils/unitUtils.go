package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

type Unit int64

const (
	splitter = "(^[1-9][0-9]*)([a-zA-Z]*)$"

	K  Unit = 1024
	KB Unit = 1000

	M  Unit = 1024 * K
	MB Unit = 1000 * KB

	G  Unit = 1024 * M
	GB Unit = 1000 * MB

	DEFAULT = G
)

var (
	unitMapping = make(map[string]Unit)
)

func getUnit(unit string) Unit {
	switch unit {
	case "K":
		return K
	case "M":
		return M
	case "G":
		return G
	default:
		panic("Unsupported")
	}
}

// GetMagnitudeAndUnit returns magnitude, unit value, error
func GetMagnitudeAndUnit(str string) (float64, Unit, error) {
	mag, unitStr, err := GetMagnitudeAndUnitStr(str)
	if err != nil {
		return 0, 0, err
	}
	unit, ok := unitMapping[unitStr]
	if !ok {
		return 0, 0, fmt.Errorf("Invalid unit")
	}
	return mag, unit, err
}

// GetMagnitudeAndUnitStr returns magnitude, unit string, error
func GetMagnitudeAndUnitStr(str string) (float64, string, error) {
	exp, err := regexp.Compile(splitter)
	if err != nil {
		return 0, "", err
	}
	match := exp.MatchString(str)
	if !match {
		return 0, "", fmt.Errorf("Invalid format")
	}
	submatches := exp.FindStringSubmatch(str)
	magnitude, err := strconv.ParseFloat(submatches[1], 64)
	if err != nil {
		return 0, "", err
	}
	unitStr := submatches[2]
	_, ok := unitMapping[unitStr]
	if !ok {
		return 0, "", fmt.Errorf("Invalid unit")
	}
	return magnitude, unitStr, nil
}

// Convert a storage size string to target unit
func Convert(str string, to Unit) float64 {
	mag, unit, err := GetMagnitudeAndUnit(str)
	CheckErr(err)
	return mag * float64(unit) / float64(to)
}

func init() {
	unitMapping["K"] = K
	unitMapping["Ki"] = K
	unitMapping["KiB"] = K
	unitMapping["KB"] = KB

	unitMapping["M"] = M
	unitMapping["Mi"] = M
	unitMapping["MiB"] = M
	unitMapping["MB"] = MB

	unitMapping["G"] = G
	unitMapping["Gi"] = G
	unitMapping["GiB"] = G
	unitMapping["GB"] = GB

	unitMapping[""] = DEFAULT

}
