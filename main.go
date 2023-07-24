package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ad           = regexp.MustCompile("<archdesc.*archdesc>")
	datePtn      = regexp.MustCompile("<date>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} -[0-9]{4}</date>")
	idPtn        = regexp.MustCompile("id=\"aspace_.{32}\"")
	parentPtn    = regexp.MustCompile("parent=\"aspace_.{32}\"")
	subDirs      = []string{"akkasah", "archives", "cbh", "fales", "nyhs", "nyuad", "poly", "tamwag", "vlp"}
	dump         bool
	changedFiles = 0
	newFiles     = 0
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

	changeFile, _ := os.Create("changedFiles.txt")
	defer changeFile.Close()
	changeWriter := bufio.NewWriter(changeFile)

	newFile, _ := os.Create("newFiles.txt")
	defer newFile.Close()
	newWriter := bufio.NewWriter(newFile)

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
				newFiles++
				newWriter.WriteString(dir1Path + "\n")
				continue
			}

			//fmt.Printf("comparing %s with %s", dir1Path, dir2path)
			originalBytes, err := GetFileBytes(dir1Path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			originalBytes = RedactEAD(originalBytes)

			newBytes, err := GetFileBytes(dir2path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			newBytes = RedactEAD(newBytes)

			if bytes.Equal(originalBytes, newBytes) != true {
				changedFiles++
				changeWriter.WriteString(dir1Path + "\n")
			}
		}
	}

	newWriter.Flush()
	changeWriter.Flush()

	fmt.Println(changedFiles, " were changed")
	fmt.Println(newFiles, " were not in previous sample set")

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

func RedactEAD(eadBytes []byte) []byte {
	eadBytes = RedactCreateDate(eadBytes)
	eadBytes = RedactedIDAttr(eadBytes)
	eadBytes = RedactedParentAttr(eadBytes)
	return eadBytes
}

func RedactCreateDate(eadBytes []byte) []byte {

	matches := datePtn.FindAllSubmatchIndex(eadBytes, 1)
	if len(matches) > 0 {

		match := matches[0]

		for i := match[0] + 6; i < match[1]-6; i++ {
			eadBytes[i] = 88
		}
	}

	return eadBytes
}

func RedactedIDAttr(eadBytes []byte) []byte {
	ids := idPtn.FindAllSubmatchIndex(eadBytes, -1)
	if len(ids) > 0 {
		for _, id := range ids {
			for i := id[0] + 11; i < id[1]-1; i++ {
				eadBytes[i] = 88
			}
		}
	}

	return eadBytes
}

func RedactedParentAttr(eadBytes []byte) []byte {
	ids := parentPtn.FindAllSubmatchIndex(eadBytes, -1)
	if len(ids) > 0 {

		for _, id := range ids {
			for i := id[0] + 15; i < id[1]-1; i++ {
				eadBytes[i] = 88
			}
		}
	}

	return eadBytes
}

func DumpEADs() {
	fmt.Println("Dumping Redacted EAD")

	fileBytes, err := GetFileBytes(os.Args[2])
	if err != nil {
		panic(err)
	}

	fileBytes = RedactEAD(fileBytes)

	fi, _ := os.Stat(os.Args[2])

	os.WriteFile(fi.Name()+"-redacted", fileBytes, 0644)

	os.Exit(0)
}
