package utils

import (
	"os"
	"path"
)

func DirExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Mkdir(path string) error {
	if !DirExist(path) {
		err := os.MkdirAll(path, 0755)
		return err
	}
	return nil
}

func SaveFile(data, p string) error {
	dir := path.Dir(p)
	Mkdir(dir)
	file, err := os.Create(p)
	defer file.Close()

	if err != nil {
		return err
	}
	file.WriteString(data)
	return nil
}
