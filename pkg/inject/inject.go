package inject

import (
	"os"
	"fmt"
)

var fd *os.File

func OpenHid(id uint64) error {
	var err error
	fd, err = os.Open(fmt.Sprintf("/dev/hidg%d", id))
	return err
}

func InjectHid(msg string) error {
	_, err := fd.WriteString(msg)
	return err
}
