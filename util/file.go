package util

import (
	"log"
	"os"
)

func DirExists(directory string) (bool, error) {
	stat, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("does not exist inner")
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}
