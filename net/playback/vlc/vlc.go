package vlc

import (
	"fmt"
	"net/url"
	"simon/conf"
	"simon/net/playback"
	"strconv"
)

func Start() playback.Server {
	run()
	port := strconv.Itoa(conf.Get().VlcPort())
	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port))
	server := New(addr)
	if err := server.Connect(); err != nil {
		panic(err)
	}
	return server
}
