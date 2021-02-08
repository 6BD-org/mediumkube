package utils

import "strconv"

type Unit int64

var (
	K Unit = 1024
	M Unit = 1024 * K
	G Unit = 1024 * M
)

func getUnit(unit byte) Unit {
	switch unit {
	case 'K':
		return K
	case 'M':
		return M
	case 'G':
		return G
	default:
		panic("Unsupported")
	}
}

func Convert(str string, to Unit) float32 {
	unit := getUnit(str[len(str)-1])
	sizeStr := str[:len(str)-1]

	size, err := strconv.ParseFloat(sizeStr, 32)
	CheckErr(err)
	return float32(size) * float32(unit) / float32(to)

}
