// +build !darwin

package cmd

func init() {
	List["setvpath"] = newCmd(_setvPath, check(count(1)), check(validFile))
}
