package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
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
	removedFiles = 0
	removeWriter *bufio.Writer
)

func init() {
	flag.BoolVar(&dump, "dump", false, "")
}

func main() {
	flag.Parse()

	if dump {
		DumpEADs()
	}

	logFile, _ := os.Create("ead-compare.log")
	defer logFile.Close()
	log.SetOutput(logFile)

	dir1 := os.Args[1]
	dir2 := os.Args[2]

	fmt.Println(os.Args[0], "\ncurrent sample set location:", dir1, " previous sample set location:")

	log.Println("[INFO] creating output file: changedFiles.txt")
	changeFile, _ := os.Create("changedFiles.txt")
	defer changeFile.Close()
	changeWriter := bufio.NewWriter(changeFile)

	log.Println("[INFO] creating outputFile: newFiles.txt")
	newFile, _ := os.Create("newFiles.txt")
	defer newFile.Close()
	newWriter := bufio.NewWriter(newFile)

	log.Println("[INFO] creating outputFile: removedFiles.txt")
	removeFile, _ := os.Create("removedFiles.txt")
	defer removeFile.Close()
	removeWriter = bufio.NewWriter(removeFile)

	FindRemoved(dir1, dir2)
	removeWriter.Flush()

	for _, subDir := range subDirs {
		fmt.Println("Comparing EADS from", subDir, "respository")
		dir1Files, err := os.ReadDir(filepath.Join(dir1, subDir))
		if err != nil {
			panic(err)
		}

		for _, dir1File := range dir1Files {

			dir1Filename := dir1File.Name()
			dir1Path := filepath.Join(dir1, subDir, dir1Filename)
			dir2path := filepath.Join(dir2, subDir, dir1Filename)
			log.Println("[DEBUG] comparing", dir1Path, "to", dir2path)

			err := FileExists(dir2path)
			if err != nil {
				newFiles++
				log.Println("[INFO]", dir2path, "does not exist in previous sample set")
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
				log.Println("[INFO]", dir1Path, "has been changed in current sampleset")
			}
		}
		newWriter.Flush()
		changeWriter.Flush()
	}

	newWriter.Flush()
	changeWriter.Flush()

	log.Println("[INFO]", changedFiles, "were changed from previous sample set")
	fmt.Println(changedFiles, " were changed")
	log.Println("[INFO]", changedFiles, "were added in current sample set")
	fmt.Println(newFiles, " were not in revious sample set")
	log.Println("[INFO]", removedFiles, "were removed in current sample set")
	fmt.Println(removedFiles, " were removed in current sample set")

}

func FindRemoved(currentDir string, prevDir string) {
	for _, subDir := range subDirs {
		prevSubdirPath := filepath.Join(prevDir, subDir)
		prevFiles, err := os.ReadDir(prevSubdirPath)
		if err != nil {
			panic(err)
		}

		currentSubDirPath := filepath.Join(currentDir, subDir)
		for _, prevFile := range prevFiles {

			if _, err := os.Stat(filepath.Join(currentSubDirPath, prevFile.Name())); err == nil {
				//do nothing
			} else if errors.Is(err, os.ErrNotExist) {
				prevFilePath := filepath.Join(subDir, prevFile.Name())
				removedFiles++
				removeWriter.WriteString(prevFilePath)
			} else {
				panic(err)
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
