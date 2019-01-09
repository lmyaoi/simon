// +build !windows

package conf

import (
	"fmt"
	"os"
)

var path = fmt.Sprintf("%v%c%v", os.Getenv("HOME"), os.PathSeparator, name)