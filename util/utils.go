package util

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
)

func ReadJsonFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func WriteJsonFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, data, DefaultFilePerm); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func FileExists(fileName string) (os.FileInfo, bool) {

	fileInfo, err := os.Lstat(fileName)

	if fileInfo != nil || (err != nil && !os.IsNotExist(err)) {
		return fileInfo, true
	}
	return nil, false
}

func TouchDir(file string) error {
	if _, ok := FileExists(file); ok {
		return nil
	}
	err := os.MkdirAll(file, DefaultFilePerm)
	if err != nil {
		return err
	}
	return nil
}

func LoadServerTlsCnf(cPath, kPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cPath, kPath)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	return cfg, err
}