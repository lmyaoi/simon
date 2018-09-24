package log

import (
	"fmt"
	"io"
	"time"
)

func Now() string {
	return time.Now().Format("[15:04:05]")
}

func now() string {
	return Now() + " "
}

func Printf(format string, a ...interface{}) {
	fmt.Print(now())
	fmt.Printf(format, a...)
}

func Println(a ...interface{}) {
	fmt.Print(now())
	fmt.Println(a...)
}

func Print(a ...interface{}) {
	fmt.Print(now())
	fmt.Print(a...)
}

func Fprintf(w io.Writer, format string, a ...interface{}) {
	fmt.Fprintf(w, now())
	fmt.Fprintf(w, format, a...)
}

func Fprintln(w io.Writer, a ...interface{}) {
	fmt.Fprintf(w, now())
	fmt.Fprintln(w, a...)
}

func Fprint(w io.Writer, a ...interface{}) {
	fmt.Fprintf(w, now())
	fmt.Fprint(w, a...)
}

func Sprintf(format string, a ...interface{}) string {
	return now() + fmt.Sprintf(format, a...)
}

func Sprintln(a ...interface{}) string {
	return now() + fmt.Sprintln(a...)
}

func Sprint(a ...interface{}) string {
	return now() + fmt.Sprint(a...)
}