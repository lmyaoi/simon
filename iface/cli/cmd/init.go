package cmd

import (
	. "simon/conf"
	. "strconv"
)

var List = make(map[string]Runner)

// fetch config
func init() {
	Initialize()
}

// populate List
func init() {
	List["help"] = newCmd(_help, check(noArgs))
	List["host"] = newCmd(_host, defaultArgs(Itoa(Get().HostingPort())), check(validPort))
	List["sethost"] = newCmd(_setHost, check(count(2)), check(validUrl))
	List["join"] = newCmd(_join, defaultArgs(Get().HostAddr(), Itoa(Get().HostPort())), check(validUrl))
	List["exit"] = newCmd(_exit, check(noArgs))
	List["status"] = newCmd(_status, check(noArgs))
	List["setvport"] = newCmd(_setvPort, check(count(1)), check(validPort))
	List["setvpath"] = newCmd(_setvPath, check(count(1)), check(validPath))
	List["setival"] = newCmd(_setIval, check(count(1)), check(validIval))
	List["save"] = newCmd(_save, check(noArgs))
}

// add aliases
func init() {
	List["list"] = List["help"]
	List["quit"] = List["exit"]
}
