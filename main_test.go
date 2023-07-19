package main

import (
	"bytes"
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
		match := bytes.Equal(origBytes, alteredBytes)
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

		match := bytes.Equal(origBytes, alteredBytes)
		if match != false {
			t.Error(fmt.Errorf("ArchDescs are the same"))
		}
	})
}

func TestDate(t *testing.T) {
	origPath := "test/mss_360_Orig.xml"

	t.Run("Test Ignore Modified CreationDate", func(t *testing.T) {
		origBytes, err := GetEadBytesWithoutModDate(origPath)
		if err != nil {
			t.Error(err)
		}

		alteredPath := "test/mss_360_Altered_CreationDate.xml"
		alteredBytes, err := GetEadBytesWithoutModDate(alteredPath)
		if err != nil {
			t.Error(err)
		}

		match := bytes.Equal(origBytes, alteredBytes)
		if match == false {
			t.Error(fmt.Errorf("Creation Date was not ignored"))
		}

	})
}
