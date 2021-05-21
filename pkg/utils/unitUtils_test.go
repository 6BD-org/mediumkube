package utils

import (
	"log"
	"testing"
)

func TestGetMagAndUnit(t *testing.T) {
	str := "100Gi"
	mag, unit, err := GetMagnitudeAndUnit(str)
	if err != nil {
		t.Fail()
	}
	if mag != 100 {
		t.Fail()
	}
	if unit != G {
		t.Fail()
	}

	str = "100"
	mag, unit, err = GetMagnitudeAndUnit(str)
	if err != nil {
		t.Fail()
	}
	if mag != 100 {
		t.Fail()
	}
	if unit != G {
		t.Fail()
	}

	str = "100GB"
	mag, unit, err = GetMagnitudeAndUnit(str)
	if err != nil {
		t.Fail()
	}
	if mag != 100 {
		t.Fail()
	}
	if unit != GB {
		t.Fail()
	}

	str = "100GC"
	mag, unit, err = GetMagnitudeAndUnit(str)
	if err == nil {
		t.Fail()
	}

	str = "03G"
	mag, unit, err = GetMagnitudeAndUnit(str)
	if err == nil {
		t.Fail()
	}

	str = "2G"
	cvted := Convert(str, M)
	log.Println(cvted)

}
