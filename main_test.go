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
	alteredPath := "test/mss_360_Altered_Create_Date.xml"

	t.Run("Test Ignore Modified CreationDate", func(t *testing.T) {
		origBytes, err := RedactCreateDate(origPath)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := RedactCreateDate(alteredPath)
		if err != nil {
			t.Error(err)
		}

		match := bytes.Equal(origBytes, alteredBytes)
		if match == false {
			t.Error(fmt.Errorf("Creation Date was not ignored"))
		}

	})
}

func TestRedactAttr(t *testing.T) {
	origPath := "test/mss_360_Orig.xml"

	t.Run("Test Getting Redacted ID attributes", func(t *testing.T) {
		alteredPath := "test/mss_360_Altered_IDs.xml"
		origBytesRedax, err := RedactedIDAttr(origPath)
		if err != nil {
			t.Error(err)
		}
		alteredBytesRedax, err := RedactedIDAttr(alteredPath)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(origBytesRedax, alteredBytesRedax) != true {
			t.Error(fmt.Errorf("Original and Altered files were not the same after redacting id attrs"))

		}
	})

	t.Run("Test Getting Redacted Parent Attrs", func(t *testing.T) {
		alteredPath := "test/mss_360_Altered_Parent.xml"

		origBytesRedax, err := RedactedIDAttr(origPath)
		if err != nil {
			t.Error(err)
		}
		alteredBytesRedax, err := RedactedParentAttr(alteredPath)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(origBytesRedax, alteredBytesRedax) != true {
			t.Error(fmt.Errorf("Original and Altered files were not the same after redacting parent attrs"))
		}
	})
}
