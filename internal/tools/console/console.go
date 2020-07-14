package console

import (
	"fmt"
)

// Log ...
func Log(data interface{}) (err error) {
	_, err = fmt.Printf("%+v\n", data)
	return
}
