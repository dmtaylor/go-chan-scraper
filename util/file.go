package util

import (
	"io"
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

func SaveFile(fullFilename string, body io.ReadCloser) error {
	outfile, err := os.Create(fullFilename)
	if err != nil {
		return err
	}
	defer outfile.Close()
	_, err = io.Copy(outfile, body)
	if err != nil {
		return err
	}
	return nil

}
