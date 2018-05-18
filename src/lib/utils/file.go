package utils

import (
	"os"
	"github.com/astaxie/beego"
	"io/ioutil"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetFileContext(filepath string) ([]byte, error) {
	if CheckFileIsExist(filepath) {
		file, err := os.Open(filepath)
		if err != nil {
			beego.Error(err)
			return make([]byte, 0), err
		}

		data, err := ioutil.ReadAll(file)
		return data, err
	}
	return make([]byte, 0), nil
}

func SaveFile(filepath string, data []byte) {
	var file *os.File = nil
	if CheckFileIsExist(filepath) {
		file, _ = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	} else {
		file, _ = os.Create(filepath)
	}

	if file != nil {
		file.Write(data)
		file.Sync()
	}
}