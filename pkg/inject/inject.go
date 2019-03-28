package inject

import (
	"os"
	"fmt"
)

var fd *os.File

func OpenHid(id uint64) error {
	var err error
	fd, err = os.OpenFile(fmt.Sprintf("/dev/hidg%d", id), os.O_WRONLY, 0)
	return err
}

func InjectHid(msg []byte) error {
	_, err := fd.Write(msg)
	return err
}
