package utils

import (
	"encoding/json"
	"io"
	"os"
)

func GetSSID(ssidname string) string{
	f, _ := os.Open("config_ssid.json")
	defer f.Close()
	var data map[string]string
	data = make(map[string]string)
	r := io.Reader(f)
	json.NewDecoder(r).Decode(&data)
	return data[ssidname]
}
