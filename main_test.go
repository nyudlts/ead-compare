package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDate(t *testing.T) {
	origPath := "test/mss_360_Orig.xml"
	alteredPath := "test/mss_360_Altered_Create_Date.xml"

	t.Run("Test Compare Redacted CreationDate", func(t *testing.T) {
		origBytes, err := GetFileBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		origBytes = RedactCreateDate(origBytes)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := GetFileBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}

		alteredBytes = RedactCreateDate(alteredBytes)
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

	t.Run("Test Compare Redacted ID attributes", func(t *testing.T) {
		alteredPath := "test/mss_360_Altered_IDs.xml"
		origBytes, err := GetFileBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		origBytes = RedactIDAttrs(origBytes)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := GetFileBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}

		alteredBytes = RedactIDAttrs(alteredBytes)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(origBytes, alteredBytes) != true {
			t.Error(fmt.Errorf("Original and Altered files were not the same after redacting id attrs"))
			DumpEAD("test.xml", origBytes, alteredBytes)
		}
	})

	t.Run("Test Compare Redacted Parent Attrs", func(t *testing.T) {
		alteredPath := "test/mss_360_Altered_Parent.xml"

		origBytes, err := GetFileBytes(origPath)
		if err != nil {
			t.Error(err)
		}

		origBytes = RedactParentAttrs(origBytes)
		if err != nil {
			t.Error(err)
		}

		alteredBytes, err := GetFileBytes(alteredPath)
		if err != nil {
			t.Error(err)
		}
		alteredBytes = RedactParentAttrs(alteredBytes)
		if err != nil {
			t.Error(err)
		}

		if bytes.Equal(origBytes, alteredBytes) != true {
			t.Error(fmt.Errorf("Original and Altered files were not the same after redacting id attrs"))

		}
	})
}
