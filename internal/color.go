package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	isatty "github.com/mattn/go-isatty"
)

const (
	// Rd red text style
	Rd = "31"
	// Grn green text style
	Grn = "32"
	// Yel yellow text style
	Yel = "33"
	// Blu blue text style
	Blu = "34"
	// Cyn cyan text style
	Cyn = "36"

	// RdBg red background style
	RdBg = "41"
	// GrnBg green background style
	GrnBg = "42"
	// YelBg yellow background style
	YelBg = "43"
	// BluBg blue background style
	BluBg = "44"
	// MgnBg magenta background style

	// B bold emphasis style
	B = "1"
	// U underline emphasis style
	U = "4"
)

var (
	noColor    bool
	once       sync.Once
	checkColor = func() {
		noColor = (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
	}

	Red    = outer(Rd)
	Green  = outer(Grn)
	Yellow = outer(Yel)
	Blue   = outer(Blu)
	Cyan   = outer(Cyn)

	RedBg    = outer(RdBg)
	GreenBg  = outer(GrnBg)
	YellowBg = outer(YelBg)
	BlueBg   = outer(BluBg)

	Bold      = outer(B)
	Underline = outer(U)
)

type inner func(msg interface{}, styles ...string) string

func outer(n string) inner {
	return func(msg interface{}, styles ...string) string {
		if noColor {
			return fmt.Sprintf("%v", msg)
		}

		b := new(bytes.Buffer)
		b.WriteString("\x1b[")
		b.WriteString(n)
		for _, s := range styles {
			b.WriteString(";")
			b.WriteString(s)
		}
		b.WriteString("m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), msg)
	}
}

func init() {
	once.Do(checkColor)
}

// represent keywords to highlight
func NewError(path string, err error) {
	log.SetOutput(os.Stderr)
	log.Printf(Underline("Receive error while processing [%s] : %v \n"), Red(path), err)
	log.SetOutput(os.Stdout)
}
