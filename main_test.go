package main

import (
	"fmt"
	"testing"
)

func TestCompareEAD(t *testing.T) {
	t.Run("testing compare changed archdescs", func(t *testing.T) {
		origPath := "test/mss_360_Orig.xml"
		origBytes, err := GetArchDescBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		alteredPath := "test/mss_360_Altered.xml"
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
