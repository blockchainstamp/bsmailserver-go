package util

import (
	"crypto"
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
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
func ParseEmailAddress(email string) (string, string, error) {
	data := strings.Split(email, "@")
	if len(data) != 2 {
		return "", "", fmt.Errorf("invalid email address")
	}

	return data[0], data[1], nil
}

func LoadDkimKey(path string) (crypto.Signer, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}

	switch strings.ToUpper(block.Type) {
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return k.(crypto.Signer), nil
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EDDSA PRIVATE KEY":
		if len(block.Bytes) != ed25519.PrivateKeySize {
			return nil, fmt.Errorf("invalid Ed25519 private key size")
		}
		return ed25519.PrivateKey(block.Bytes), nil
	default:
		return nil, CfgInvalidDkim
	}
}
