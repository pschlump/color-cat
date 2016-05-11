package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/pschlump/Go-FTL/server/sizlib"
	"github.com/pschlump/MiscLib"
)

var Debug = flag.Bool("debug", false, "Debug flag")                                // 0
var LineNo = flag.Bool("number", false, "Add Line Numbers")                        // 2
var NonBlank = flag.Bool("nonblank", false, "Add Line Numbers to non-blank lines") // 3
var UnBuffer = flag.Bool("unbuffer", false, "Do not buffer output")                // 4
var ColorForced = flag.Bool("forcedcolor", false, "Turn on color even if piped")   // 5
var Squeeze = flag.Bool("squeeze", false, "Remove ajacent blank lines")            // 6
var Color = flag.String("color", "red", "Color to use")                            // 7
// var Cfg = flag.String("cfg", "../global_cfg.json", "Configuration file")           // 1
func init() {
	flag.BoolVar(Debug, "D", false, "Debug flag")                             // 0
	flag.BoolVar(LineNo, "n", false, "Add Line Numbers")                      // 2
	flag.BoolVar(NonBlank, "b", false, "Add Line Numbers to non-blank lines") // 3
	flag.BoolVar(UnBuffer, "u", false, "Do not buffer output")                // 4
	flag.BoolVar(ColorForced, "F", false, "Turn on color even if piped")      // 5
	flag.BoolVar(Squeeze, "s", false, "Remove ajacent blank lines")           // 6
	flag.StringVar(Color, "c", "red", "Color to use")                         // 7
	// flag.StringVar(Cfg, "C", "../global_cfg.json", "Configuration file")      // 1
}

var ColorOn string
var ColorReset string

func main() {

	flag.Parse()
	fns := flag.Args()

	if !MiscLib.StdOutPiped() {
		if *Debug {
			fmt.Printf("Is NOT Piped\n")
		}
		ColorReset = MiscLib.ColorReset
		switch *Color {
		case "red":
			ColorOn = MiscLib.ColorRed
		case "green":
			ColorOn = MiscLib.ColorGreen
		case "yellow":
			ColorOn = MiscLib.ColorYellow
		case "blue":
			ColorOn = MiscLib.ColorBlue
		case "black":
			ColorOn = MiscLib.ColorBlack
		case "cyan":
			ColorOn = MiscLib.ColorCyan
		case "magenta":
			ColorOn = MiscLib.ColorMagenta
		default:
			fmt.Fprintf(os.Stderr, "Error: color must be one of [ red | yellow | green | blue | black | cyan | magenta ]\n")
			flag.Usage()
			os.Exit(1)
		}
	} else {
		if *Debug {
			fmt.Printf("Is Piped\n")
		}
	}

	if *ColorForced {
		ColorOn = MiscLib.ColorRed
		ColorReset = MiscLib.ColorReset
	}

	if len(fns) == 0 {
		CatFile(os.Stdin)
	} else {
		for _, fn := range fns {
			if fn == "-" {
				CatFile(os.Stdin)
			} else {
				fi, err := sizlib.Fopen(fn, "r")
				if err != nil {
					// xyzzy
				} else {
					defer fi.Close()
					CatFile(fi)
				}
			}
		}
	}
}

func CatFile(fi *os.File) {
	fmt.Printf("%s", ColorOn)
	scanner := bufio.NewScanner(fi)
	line_no := 1
	nb := 0
	for scanner.Scan() {
		s := scanner.Text()
		isBlank := false
		output := true
		if *NonBlank || *Squeeze {
			isBlank = IsBlank(s)
		}

		if isBlank {
			nb++
		} else {
			nb = 0
		}

		if *Squeeze && nb >= 2 {
			output = false
		}

		if *Debug {
			fmt.Printf("nb=%d isBlank=%v >%s<\n", nb, isBlank, MakePrintable(s))
		}

		if *NonBlank {
			if isBlank {
				if output {
					fmt.Printf("     %s\n", s)
				}
			} else {
				if output {
					fmt.Printf("%3d: %s\n", line_no, s)
				}
				line_no++
			}
		} else if *LineNo {
			if output {
				fmt.Printf("%3d: %s\n", line_no, s)
			}
			line_no++
		} else {
			if output {
				fmt.Printf("%s\n", s)
			}
		}

		if *UnBuffer {
			// xyzzy - hm... os.Stdout.Flush()
			// os.Flush(os.Stdout)
		}
	}
	fmt.Printf("%s", ColorReset)
}

func MakePrintable(s string) (rv string) {
	var t []byte
	for _, c := range s {
		if c >= ' ' && c < 128 {
			t = append(t, byte(c))
		} else if c > 128 {
			t = append(t, []byte("M-")...)
			t = append(t, byte(c&0x7f))
		} else {
			t = append(t, '^')
			t = append(t, byte((c+'A')&0x7f))
		}
	}
	rv = string(t)
	return
}

func IsBlank(s string) bool {
	for _, x := range s {
		if x == ' ' || x == '\t' || x == '\f' || x == '\r' || x == '\n' {
		} else {
			return false
		}
	}
	return true
}
