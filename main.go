package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ad        = regexp.MustCompile("<archdesc.*archdesc>")
	datePtn   = regexp.MustCompile("<date>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} -[0-9]{4}</date>")
	idPtn     = regexp.MustCompile("id=\"aspace_.{32}\"")
	parentPtn = regexp.MustCompile("parent=\"aspace_.{32}\"")
	subDirs   = []string{"archives", "fales", "tamwag", "vlp"}
	dump      bool
)

func init() {
	flag.BoolVar(&dump, "dump", false, "")
}

func main() {
	flag.Parse()
	if dump {
		DumpEADs()
	}

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

			//fmt.Printf("comparing %s with %s", dir1Path, dir2path)
			originalBytes, err := GetFileBytes(dir1Path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			originalBytes, err = RedactedParentAttr(originalBytes)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			newBytes, err := GetFileBytes(dir2path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			newBytes, err = RedactEAD(newBytes)
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

func GetFileBytes(path string) ([]byte, error) {
	eadBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	eadBytes = bytes.ReplaceAll(eadBytes, []byte("\n"), []byte(""))
	return eadBytes, nil
}

func RedactEAD(eadBytes []byte) ([]byte, error) {
	var err error
	eadBytes, err = RedactCreateDate(eadBytes)
	if err != nil {
		return nil, err
	}

	eadBytes, err = RedactedIDAttr(eadBytes)
	if err != nil {
		return nil, err
	}

	eadBytes, err = RedactedParentAttr(eadBytes)
	if err != nil {
		return nil, err
	}

	return eadBytes, nil
}

func RedactCreateDate(eadBytes []byte) ([]byte, error) {

	matches := datePtn.FindAllSubmatchIndex(eadBytes, 1)
	if len(matches) != 1 {
		return nil, fmt.Errorf("Could not find creation date in file")
	}

	match := matches[0]

	for i := match[0] + 6; i < match[1]-5; i++ {
		eadBytes[i] = 88
	}

	return eadBytes, nil
}

func RedactedIDAttr(eadBytes []byte) ([]byte, error) {
	ids := idPtn.FindAllSubmatchIndex(eadBytes, -1)
	if len(ids) < 1 {
		return nil, fmt.Errorf("Could not find any id attrs")
	}

	for _, id := range ids {
		for i := id[0] + 11; i < id[1]-1; i++ {
			eadBytes[i] = 88
		}
	}

	return eadBytes, nil
}

func RedactedParentAttr(eadBytes []byte) ([]byte, error) {
	ids := parentPtn.FindAllSubmatchIndex(eadBytes, -1)
	if len(ids) < 1 {
		return nil, fmt.Errorf("Could not find any parent attrs")
	}

	for _, id := range ids {
		for i := id[0] + 15; i < id[1]-1; i++ {
			eadBytes[i] = 88
		}
	}

	return eadBytes, nil
}

func DumpEADs() {
	fmt.Println("Dumping Redacted EAD")

	fileBytes, err := GetFileBytes(os.Args[2])
	if err != nil {
		panic(err)
	}

	fileBytes, err = RedactEAD(fileBytes)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fi, _ := os.Stat(os.Args[2])

	os.WriteFile(fi.Name()+"-redacted", fileBytes, 0644)

	os.Exit(0)
}
