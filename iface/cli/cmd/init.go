package cmd

import (
	. "simon/conf"
	"sort"
	. "strconv"
)

var List = make(map[string]Runner)
var listKeys []string

// fetch config
func init() {
	Initialize()
}

// populate List
func init() {
	List["list"] = newCmd(_list, check(noArgs))
	List["host"] = newCmd(_host, defaultArgs(Itoa(Get().HostingPort)), check(validPort))
	List["sethost"] = newCmd(_setHost, check(count(2)), check(validUrl))
	List["join"] = newCmd(_join, defaultArgs(Get().HostAddr, Itoa(Get().HostPort)), check(validUrl))
	List["exit"] = newCmd(_exit, check(noArgs))
	List["status"] = newCmd(_status, check(noArgs))
	List["setvport"] = newCmd(_setvPort, check(count(1)), check(validPort))
	List["setival"] = newCmd(_setIval, check(count(1)), check(validIval))
	List["save"] = newCmd(_save, check(noArgs))
}

func getListKeys() []string {
	if listKeys == nil {
		listKeys = make([]string, len(List))
		i := 0
		for k := range List  {
			listKeys[i] = k
			i++
		}
		sort.Strings(listKeys)
	}
	return listKeys
}