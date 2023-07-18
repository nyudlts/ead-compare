package main

import (
	"fmt"
	"testing"
)

func TestCompareEAD(t *testing.T) {
	origPath := "test/mss_360_Orig.xml"
	origBytes, err := GetArchDescBytes(origPath)
	if err != nil {
		t.Error(err)
	}

	t.Run("testing compare changed headers", func(t *testing.T) {
		alteredPath := "test/mss_360_AlteredHeader.xml"
		alteredBytes, err := GetArchDescBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}
		match := (string(origBytes) == string(alteredBytes))
		if match != true {
			t.Error(fmt.Errorf("ArchDescs are not the same"))
		}
	})

	t.Run("testing compare changed archdescs", func(t *testing.T) {

		alteredPath := "test/mss_360_AlteredArchDesc.xml"
		alteredBytes, err := GetArchDescBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}

		match := (string(origBytes) == string(alteredBytes))
		if match != false {
			t.Error(fmt.Errorf("ArchDescs are the same"))
		}
	})
}
