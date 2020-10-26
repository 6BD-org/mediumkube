package utils

// Simple error checking. Panic if err found
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}
