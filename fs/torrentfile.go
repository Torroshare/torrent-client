package fs

import (
	"fmt"
	"os"
)

func Opendottorrent(fulladr string) ([]byte, error) {
	file, err := os.ReadFile(fulladr)
	if err != nil {
		return nil, fmt.Errorf("Error opening torrent file")
	}
	fmt.Printf(string(file))
	return file, nil
}
