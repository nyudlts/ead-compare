package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var ad = regexp.MustCompile("<archdesc.*archdesc>")
var datePtn = regexp.MustCompile("<date>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} -[0-9]{4}</date>")
var subDirs = []string{"akkasah", "archives", "cbh", "fales", "nyhs", "nyuad", "poly", "tamwag", "vlp"}

func main() {
	dir1 := os.Args[1]
	dir2 := os.Args[2]
	fmt.Println(dir1, dir2)

	for _, subDir := range subDirs {
		dir1Files, err := os.ReadDir(filepath.Join(dir1, subDir))
		if err != nil {
			panic(err)
		}

		for _, dir1File := range dir1Files {
			dir1Filename := dir1File.Name()
			dir1Path := filepath.Join(dir1, subDir, dir1Filename)
			dir2path := filepath.Join(dir2, subDir, dir1Filename)

			err := FileExists(dir2path)
			if err != nil {
				fmt.Println("file2 does not exist: ", dir2path)
				continue
			}

			fmt.Printf("comparing %s with %s", dir1Path, dir2path)
			originalBytes, err := GetEadBytesWithoutModDate(dir1Path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			newBytes, err := GetEadBytesWithoutModDate(dir2path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if bytes.Equal(originalBytes, newBytes) != true {
				fmt.Println(dir1Path, "has changed")
			}
		}
	}
}

func FileExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return err

	} else {
		return err
	}
}

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

func GetEadBytesWithRedactedCreateDate(path string) ([]byte, error) {
	eadBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	eadBytes = bytes.ReplaceAll(eadBytes, []byte("\n"), []byte(""))

	matches := datePtn.FindAllSubmatchIndex(eadBytes, 1)
	if len(matches) != 1 {
		return []byte{}, fmt.Errorf("Could not find creation date in: %s", path)
	}

	match := matches[0]

	for i := match[0] + 6; i < match[1]-5; i++ {
		eadBytes[i] = 88
	}

	return eadBytes, nil
}
