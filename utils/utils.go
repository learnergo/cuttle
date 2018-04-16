package utils

import (
	"os"
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
		err := os.MkdirAll(path, os.ModePerm)
		return err
	}
	return nil
}

func SaveFile(data, path string) error {
	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return err
	}
	file.WriteString(data)
	return nil
}
