// +build !darwin

package cmd

var setVPath = &el{}

func init() {
	setVPath.Run = newCmd(_setvPath, check(count(1)), check(validFile))
	setVPath.Name = "setvpath"
	setVPath.Desc = "sets the path of the vlc executable"
	setVPath.Man = "todo: man\n"
}

func init() {
	es = append(es, setVPath)
}
