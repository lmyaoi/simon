// Code generated by "stringer -type=Signal"; DO NOT EDIT.

package ticker

import "strconv"

const _Signal_name = "OnOffKill"

var _Signal_index = [...]uint8{0, 2, 5, 9}

func (i Signal) String() string {
	if i < 0 || i >= Signal(len(_Signal_index)-1) {
		return "Signal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Signal_name[_Signal_index[i]:_Signal_index[i+1]]
}