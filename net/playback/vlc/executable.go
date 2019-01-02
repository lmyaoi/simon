// +build !darwin

package vlc

import "simon/pref"

func Executable() string {
	return pref.Get().VlcPath()
}
