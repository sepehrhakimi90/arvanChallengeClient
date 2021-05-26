package utils

import (
	"log"
)

func LogError(file, function string, err error) {
	log.Printf("%s::%s::%s::%s", "ERROR", file, function, err)
}