package cmd

import (
	"simon/pref"
)

var List = make(map[string]Runner)

// fetch preferences
func init() {
	pref.Initialize()
}

// populate List
func init() {
	List["help"] = newCmd(_help, check(noArgs))
	List["host"] = newCmd(_host, defaultArgs(pref.Get().HostingPortStr()), check(validPort))
	List["join"] = newCmd(_join, defaultArgs(pref.Get().HostStr()), check(validUrl))
	List["exit"] = newCmd(_exit, check(noArgs))
	List["status"] = newCmd(_status, check(noArgs))
	List["vlc-port"] = newCmd(_vlcPort, defaultArgs(pref.Get().VlcPortStr()), check(validPort))
	List["pref"] = newCmd(_pref, check(noArgs))
}

// add aliases
func init() {
	List["list"] = List["help"]
	List["quit"] = List["exit"]
}
