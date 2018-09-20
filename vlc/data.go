package vlc

import (
	"fmt"
	"time"
)

type Status struct {

}

var Latest Status

func (*Status) String() string {
	now := time.Now()
	return fmt.Sprintf("[%v:%v:%v] <INSERT DATA HERE>", now.Hour(), now.Minute(), now.Second())
}
