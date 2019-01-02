package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type cmd struct {
	f    func([]string)
	opts []option
}

func newCmd(f func([]string), opts ...option) *cmd {
	return &cmd{f, opts}
}

type option func([]string) ([]string, error)

type _check func([]string) error

func (cmd *cmd) Run(args []string) {
	var err error
	for _, opt := range cmd.opts {
		args, err = opt(args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	cmd.f(args)
}

func defaultArgs(defArgs ...string) option {
	return func(args []string) ([]string, error) {
		switch len(args) {
		case 0:
			return defArgs, nil
		case len(defArgs):
			return args, nil
		default:
			return nil, errors.New("unsupported argument count")
		}
	}
}

func check(c _check) option {
	return func(args []string) ([]string, error) {
		return args, c(args)
	}
}

func count(n int) _check {
	return func(args []string) error {
		if len(args) != n {
			return fmt.Errorf("expected exactly %v arguments", n)
		}
		return nil
	}
}

var noArgs = count(0)

func validUrl(args []string) error {
	_, err := url.Parse(fmt.Sprintf("http://%v:%v", args[0], args[1]))
	return err
}

func validPort(args []string) error {
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	if i < 0 {
		return errors.New("invalid port")
	}
	return nil
}
