package main

import (
	"fmt"
	"os"
)

const usage = `Usage: lsp -[aAFcrtS] [file ...]
  -F  Append file type indicator
  -a  Show all files
  -c  Use ctime
  -u  Do not show colors
  -S  Sort by size
  -r  Reverse sort
  -t  Sort by time
  -h  Show this help`

var args = struct {
	all      bool
	classify bool
	color    bool
	ctime    bool
	reverse  bool
	sorter   sortFunc
	rest     []string
}{
	color:  true,
	sorter: sortByVer,
	rest:   make([]string, 0, len(os.Args[1:])),
}

func parseArgs() {
	for i, s := range os.Args[1:] {
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			args.rest = append(args.rest, s)
			continue
		}
		if s[1] == '-' && len(s) == 2 { // "--" ends args
			args.rest = append(args.rest, os.Args[i+1:]...)
			break
		}
		for i := 1; i < len(s); i++ {
			f := s[i]
			switch f {
			case 'a':
				args.all = true
			case 'F':
				args.classify = true
			case 'c':
				args.ctime = true
			case 'r':
				args.reverse = true
			case 't':
				args.sorter = sortByTime
			case 'S':
				args.sorter = sortBySize
			case 'u':
				args.color = false
			case 'h':
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "unsupported argument '%c'\n", f)
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(1)
			}
		}
	}
	if len(args.rest) == 0 {
		args.rest = []string{"."}
	}
}
