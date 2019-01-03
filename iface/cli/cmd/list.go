package cmd

import "sort"

type list struct {
	mapping      map[string]Runner
	l            []string
	cmdMaxLength int
}

func newList() *list {
	return &list{make(map[string]Runner), make([]string, 0), 0}
}

func (l *list) set(key string, r Runner) {
	l.mapping[key] = r
	l.l = append(l.l, key)
	if len(key) > l.cmdMaxLength {
		l.cmdMaxLength = len(key)
	}
}

func (l *list) Get(key string) (Runner, bool) {
	r, ok := l.mapping[key]
	return r, ok
}

func (l *list) sortList() {
	sort.Strings(l.l)
}

func (l *list) list() []string {
	return l.l
}
