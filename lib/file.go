package lib

import (
	"log"
	"os"
)

func IsFileExist(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Println(err)
		return false
	}
	return stat.IsDir() == false
}
