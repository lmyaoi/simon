package cmd

import (
	"simon/conf"
	"strconv"
)

var es = make([]*el, 0)
var L *List

var (
	list = &el{}
	man = &el{}
	host = &el{}
	setHost = &el{}
	join = &el{}
	setVPort = &el{}
	setIval = &el{}
	status = &el{}
	save = &el{}
	exit = &el{}
)

func init() {
	conf.Initialize()
}

// list
func init() {
	list.Run = newCmd(_list, check(noArgs))
	list.Name = "list"
	list.Desc = "prints a list of all the supported commands"
	list.Man = "todo: man\n"
}

// list
func init() {
	man.Run = newCmd(_man, check(count(1)))
	man.Name = "man"
	man.Desc = "prints the man page of a specified command"
	man.Man = "todo: man\n"
}

// host
func init() {
	host.Run = newCmd(_host, defaultArgs(strconv.Itoa(conf.Get().HostingPort)))
	host.Name = "host"
	host.Desc = "hosts a connection (optional: the port at which to listen)"
	host.Man = "todo: man\n"
}

// setHost
func init() {
	setHost.Run = newCmd(_setHost, check(count(2)), check(validUrl))
	setHost.Name = "sethost"
	setHost.Desc = "sets the ip address and port of the host to connect to"
	setHost.Man = "todo: man\n"
}

// join
func init() {
	join.Run = newCmd(_join, defaultArgs(conf.Get().HostAddr, strconv.Itoa(conf.Get().HostPort)), check(validUrl))
	join.Name = "join"
	join.Desc = "joins a connection (optional: specify the ip address and port of the host)"
	join.Man = "todo: man\n"
}

// setVPort
func init() {
	setVPort.Run = newCmd(_setvPort, check(count(1)), check(validPort))
	setVPort.Name = "setvport"
	setVPort.Desc = "sets the port at which to host the vlc playback server"
	setVPort.Man = "todo: man\n"
}

// setIval
func init() {
	setIval.Run = newCmd(_setIval, check(count(1)), check(validIval))
	setIval.Name = "setival"
	setIval.Desc = "sets the interval at which to poll the host"
	setIval.Man = "todo: man\n"
}

// status
func init() {
	status.Run = newCmd(_status, check(noArgs))
	status.Name = "status"
	status.Desc = "prints the current config being used"
	status.Man = "todo: man\n"
}

// save
func init() {
	save.Run = newCmd(_save, check(noArgs))
	save.Name = "save"
	save.Desc = "saves the current config to the config file"
	save.Man = "todo: man\n"
}

// exit
func init() {
	exit.Run = newCmd(_exit, check(noArgs))
	exit.Name = "exit"
	exit.Desc = "exits the program cleanly"
	exit.Man = "todo: man\n"
}

func init() {
	es = append(es,
		list,
		man,
		host,
		setHost,
		join,
		setVPort,
		setIval,
		status,
		save,
		exit,
		)
}

func Init() {
	L = newList(es...)
}
