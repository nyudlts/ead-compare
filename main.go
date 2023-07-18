package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
)

var ad = regexp.MustCompile("<archdesc.*archdesc>")

func main() {}

func GetArchDescBytes(path string) ([]byte, error) {
	eadBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	eadBytes = bytes.ReplaceAll(eadBytes, []byte("\n"), []byte(""))
	matches := ad.FindAll(eadBytes, -1)
	if len(matches) == 1 {
		return matches[0], nil
	}

	return []byte{}, fmt.Errorf("wtf")
}
