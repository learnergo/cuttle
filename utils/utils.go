package utils

import (
	"io"
	"io/ioutil"
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

//仅支持单层文件夹复制
func CopyDir(src string, dest string) error {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		if !fi.IsDir() {
			err = CopyFile(src+"/"+fi.Name(), dest+"/"+fi.Name())
			if err != nil {
				return nil
			}
		}
	}
	return err
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	dir := path.Dir(dst)
	Mkdir(dir)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
