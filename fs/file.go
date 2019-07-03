package fs

import (
	"bufio"
	"github.com/BakerHub/trivial/check"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MustClose(f io.Closer) {
	check.Check(f.Close())
}

func MustCopyFile(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src) // #nosec
	check.Check(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	check.Check(err)
}

func CopyFile(src string, dst string) error {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src) // #nosec
	if err != nil {
		return err
	}
	// Write data to dst
	return ioutil.WriteFile(dst, data, 0644)
}

func Exists(pathname string) bool {
	_, err := os.Stat(pathname)

	return err == nil
}

func EnsureDir(directory string) {
	if Exists(directory) {
		return
	}
	check.Check(os.MkdirAll(directory, 0750))
}

func ListDirectory(root string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	check.Check(err)

	return files
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path) // #nosec
	if err != nil {
		return nil, err
	}
	defer MustClose(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
