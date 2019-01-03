package cmd

import (
	. "simon/conf"
	. "strconv"
)

var List = newList()

// fetch config
func init() {
	Initialize()
}

// populate List
func init() {
	List.set("list", newCmd(_list, check(noArgs)))
	List.set("host", newCmd(_host, defaultArgs(Itoa(Get().HostingPort)), check(validPort)))
	List.set("sethost", newCmd(_setHost, check(count(2)), check(validUrl)))
	List.set("join", newCmd(_join, defaultArgs(Get().HostAddr, Itoa(Get().HostPort)), check(validUrl)))
	List.set("exit", newCmd(_exit, check(noArgs)))
	List.set("status", newCmd(_status, check(noArgs)))
	List.set("setvport", newCmd(_setvPort, check(count(1)), check(validPort)))
	List.set("setival", newCmd(_setIval, check(count(1)), check(validIval)))
	List.set("save", newCmd(_save, check(noArgs)))
}

func init() {
	List.sortList()
}
