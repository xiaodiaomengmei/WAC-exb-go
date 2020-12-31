package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// load json file and convert  json string to map
func LoadJsonAndToMap(path string) (map[string]string, error) {
	byteValue, err_readfile := ioutil.ReadFile(path)
	if err_readfile != nil {
		log.Printf("Read config_ssid.json with error: %+v\n", err_readfile)
	}
	m := make(map[string]string)
	err_Unmarshal := json.Unmarshal(byteValue, &m)
	if err_Unmarshal != nil {
		log.Printf("Unmarshal config_ssid.json with error: %+v\n", err_Unmarshal)
		return nil, err_Unmarshal
	}
	return m, nil
}

// Convert map json string
func MapToJson(m map[string]string) (string, error) {
	jsonByte, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}

	return string(jsonByte), nil
}

// Convert map json []byte
func MapToByte(m map[string]string) ([]byte ,error) {
	jsonByte, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Printf("Marshal to config_ssid.json with error: %+v\n", err)
	}
	return jsonByte, nil
}