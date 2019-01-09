package conf

import (
	"fmt"
	"os"
)

var path = fmt.Sprintf("%v%c%v", os.Getenv("HOMEPATH"), os.PathSeparator, name)