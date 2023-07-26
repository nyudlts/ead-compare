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
	datePtn      = regexp.MustCompile("<date>[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} -[0-9]{4}</date>")
	idPtn        = regexp.MustCompile("id=\"aspace_.{32}\"")
	parentPtn    = regexp.MustCompile("parent=\"aspace_.{32}\"")
	subDirs      = []string{"akkasah", "archives", "cbh", "fales", "nyhs", "nyuad", "poly", "tamwag", "vlp"}
	dump         bool
	changedFiles = 0
	newFiles     = 0
	removedFiles = 0
	removeWriter *bufio.Writer
	currentDir   string
	prevDir      string
)

func init() {
	flag.BoolVar(&dump, "dump", false, "")
	flag.StringVar(&currentDir, "current", "", "")
	flag.StringVar(&prevDir, "prev", "", "")
}

func main() {
	flag.Parse()

	logFile, _ := os.Create("ead-compare.log")
	defer logFile.Close()
	log.SetOutput(logFile)

	fmt.Println(os.Args[0], "\ncurrent sample set location:", currentDir, " previous sample set location:", prevDir)

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

	fmt.Println("Checking for ead files removed from previous sample set")
	FindRemoved(currentDir, prevDir)
	removeWriter.Flush()

	for _, subDir := range subDirs {
		fmt.Println("Comparing EADS from", subDir, "respository")
		dir1Files, err := os.ReadDir(filepath.Join(currentDir, subDir))
		if err != nil {
			panic(err)
		}

		for _, dir1File := range dir1Files {

			dir1Filename := dir1File.Name()
			dir1Path := filepath.Join(currentDir, subDir, dir1Filename)
			dir2path := filepath.Join(prevDir, subDir, dir1Filename)
			log.Println("[DEBUG] comparing", dir1Path, "to", dir2path)

			err := FileExists(dir2path)
			if err != nil {
				newFiles++
				log.Println("[INFO]", dir2path, "does not exist in previous sample set")
				newWriter.WriteString(dir1Path + "\n")
				continue
			}

			//get the redacted bytes of file in current set
			originalBytes, err := GetRedactedEADByteSlice(dir1Path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			//get the redacted bytes of file in current set
			newBytes, err := GetRedactedEADByteSlice(dir2path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			//check if the byte slices are different
			if bytes.Equal(originalBytes, newBytes) != true {
				changedFiles++
				changeWriter.WriteString(dir1Path + "\n")
				log.Println("[INFO]", dir1Path, "has been changed in current sampleset")
				//dump
				if dump {
					if err := DumpEAD(dir1File.Name(), originalBytes, newBytes); err != nil {
						panic(err)
					}
				}

			}
		}

		//flush the writers
		newWriter.Flush()
		changeWriter.Flush()
	}

	//flush the writers
	newWriter.Flush()
	changeWriter.Flush()

	log.Println("[INFO]", changedFiles, "were changed from previous sample set")
	fmt.Println(changedFiles, " were changed")

	log.Println("[INFO]", newFiles, "were added in current sample set")
	fmt.Println(newFiles, " EAD files were added to current sample set")

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

func GetRedactedEADByteSlice(path string) ([]byte, error) {
	eadBytes, err := GetFileBytes(path)
	if err != nil {
		return nil, err
	}
	eadBytes = RedactCreateDate(eadBytes)
	eadBytes = RedactIDAttrs(eadBytes)
	eadBytes = RedactParentAttrs(eadBytes)
	return eadBytes, nil
}

func GetFileBytes(path string) ([]byte, error) {
	eadBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	eadBytes = bytes.ReplaceAll(eadBytes, []byte("\n"), []byte(""))
	return eadBytes, nil
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

func RedactIDAttrs(eadBytes []byte) []byte {
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

func RedactParentAttrs(eadBytes []byte) []byte {
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

func DumpEAD(filename string, origEad []byte, newEad []byte) error {
	err := os.Mkdir("dump", 0777)
	currentDir := filepath.Join("dump", "current")
	err = os.Mkdir(currentDir, 0777)
	prevDir := filepath.Join("dump", "previous")
	err = os.Mkdir(prevDir, 0777)

	err = os.WriteFile(filepath.Join(currentDir, filename), origEad, 0777)
	err = os.WriteFile(filepath.Join(prevDir, filename), newEad, 0777)

	if err != nil {
		return err
	} else {
		return nil
	}
}
