package main

import (
	"encoding/json"
	"log"
)

func StructToJsonString(obj interface{}) string {
	jBytes, err := json.Marshal(obj)
	if err != nil {
		log.Println("marshal error")
		return ""
	}
	return string(jBytes)
}
