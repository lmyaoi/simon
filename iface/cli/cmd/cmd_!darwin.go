// +build !darwin

package cmd

import "simon/conf"

func _setvPath(args []string) {
	conf.Get().SetVlcPath(args[0])
}
