package utils

import "k8s.io/klog/v2"

// CheckErr Simple error checking. Panic if err found
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

// WarnErr Not exit but log the error message
func WarnErr(err error) {
	if err != nil {
		klog.Error(err)
	}
}
