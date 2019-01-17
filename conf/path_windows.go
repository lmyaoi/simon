package conf

import (
	"fmt"
	"os"
)

var filepath = fmt.Sprintf("%v%c%v", os.Getenv("USERPROFILE"), os.PathSeparator, name)