package conf

import (
	"fmt"
	"os"
)

var filepath = fmt.Sprintf("%v%c%v", os.Getenv("HOMEPATH"), os.PathSeparator, name)