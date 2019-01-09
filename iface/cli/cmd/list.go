package cmd

import (
	"fmt"
	"sort"
	"strings"
)

type List struct {
	lookup        map[string]*el
	order         []string
	maxLengthName int
}

type el struct {
	Name string
	Run  Runner
	Desc string
	Man  string
}

func newList(es ...*el) *List {
	l := &List{lookup: make(map[string]*el, len(es)), order: make([]string, len(es))}
	for i, e := range es {
		l.lookup[e.Name] = e
		l.order[i] = e.Name
		if len(e.Name) > l.maxLengthName {
			l.maxLengthName = len(e.Name)
		}
	}
	sort.Strings(l.order)
	return l
}

func (l *List) init() {
	*l = *newList()
}

func (l *List) Get(key string) (Runner, bool) {
	e, ok := l.lookup[key]
	if ok {
		return e.Run, ok
	}
	return nil, ok
}

func (l *List) Man(key string) (string, bool) {
	v, ok := l.lookup[key]
	if ok {
		return v.Man, ok
	}
	return "", ok
}

func (l *List) String() string {
	builder := &strings.Builder{}
	format := fmt.Sprint("\033", "[1m", `%-`, l.maxLengthName+4, `s`, "\033", "[0m", `%s`, "\n")
	for _, name := range l.order {
		_, _ = fmt.Fprintf(builder, format, name, l.lookup[name].Desc)
	}
	return builder.String()
}
