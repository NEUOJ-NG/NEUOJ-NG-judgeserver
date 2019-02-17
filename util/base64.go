package util

import (
	"encoding/base64"
	"io/ioutil"
)

// get file encoded by base64
func GetFileBase64(path string) (string, error) {
	if content, err := ioutil.ReadFile(path); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(content), nil
	}
}
