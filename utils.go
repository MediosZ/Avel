package avel

import (
	"io"
	"encoding/json"
)

func Json(i interface{}) string{
	jsondata, err := json.Marshal(i)
	if err != nil{
		return "something was wrong"
	} else{
		return string(jsondata)
	}
}

func Decode(raw_reader io.Reader) *map[string]interface{} {
	respbody := make(map[string]interface{})
	json.NewDecoder(raw_reader).Decode(&respbody)
	return &respbody
}

