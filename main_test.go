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
		origBytes, err := GetFileBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		origBytes, err = RedactCreateDate(origBytes)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := GetFileBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err = RedactCreateDate(alteredBytes)
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
		origBytes, err := GetFileBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		origBytes, err = RedactedIDAttr(origBytes)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := GetFileBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}
		alteredBytes, err = RedactedIDAttr(alteredBytes)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(origBytes, alteredBytes) != true {
			t.Error(fmt.Errorf("Original and Altered files were not the same after redacting id attrs"))

		}
	})

	t.Run("Test Getting Redacted Parent Attrs", func(t *testing.T) {
		alteredPath := "test/mss_360_Altered_Parent.xml"

		origBytes, err := GetFileBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		origBytes, err = RedactedParentAttr(origBytes)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := GetFileBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}
		alteredBytes, err = RedactedParentAttr(alteredBytes)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(origBytes, alteredBytes) != true {
			t.Error(fmt.Errorf("Original and Altered files were not the same after redacting id attrs"))

		}
	})
}
