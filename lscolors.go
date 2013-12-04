package main

import (
	"fmt"
	"os"
)

var lsColorSuffix = make(map[string]string)
var lsColorTypes = [...]string{
	"\033[",  // lc: Left of color sequence
	"m",      // rc: Right of color sequence
	"",       // ec: End color (replaces lc+no+rc)
	"0",      // rs: Reset to ordinary colors
	"",       // no: Normal
	"",       // fi: File: default
	"01;34",  // di: Directory: bright blue
	"01;36",  // ln: Symlink: bright cyan
	"33",     // pi: Pipe: yellow/brown
	"01;35",  // so: Socket: bright magenta
	"01;33",  // bd: Block device: bright yellow
	"01;33",  // cd: Char device: bright yellow
	"",       // mi: Missing file: undefined
	"",       // or: Orphaned symlink: undefined
	"01;32",  // ex: Executable: bright green
	"01;35",  // do: Door: bright magenta
	"37;41",  // su: setuid: white on red
	"30;43",  // sg: setgid: black on yellow
	"37;44",  // st: sticky: black on blue
	"34;42",  // ow: other-writable: blue on green
	"30;42",  // tw: ow w/ sticky: black on green
	"30;41",  // ca: black on red
	"",       // mh: disabled by default
	"\033[K", // cl: clear to end of line
}

func isColored(t indicator) bool {
	return lsColorTypes[t] != ""
}

func color(name string, in indicator) string {
	if in == TypeFile {
		for i := 0; i < len(name); i++ {
			if name[i] == '.' {
				if v, ok := lsColorSuffix[name[i:]]; ok {
					return v
				}
			}
		}
	}
	return lsColorTypes[in]
}

var indicatorNamesMap = map[string]indicator{
	"lc": TypeLeft,
	"rc": TypeRight,
	"ec": TypeEnd,
	"rs": TypeReset,
	"no": TypeNorm,
	"fi": TypeFile,
	"di": TypeDir,
	"ln": TypeLink,
	"pi": TypeFifo,
	"so": TypeSock,
	"bd": TypeBlk,
	"cd": TypeChr,
	"mi": TypeMissing,
	"or": TypeOrphan,
	"ex": TypeExec,
	"do": TypeDoor,
	"su": TypeSetuid,
	"sg": TypeSetgid,
	"st": TypeSticky,
	"ow": TypeOtherWritable,
	"tw": TypeStickyOtherWritable,
	"ca": TypeCap,
	"mh": TypeMultihardlink,
	"cl": TypeClrToEol,
}

func parseLSColor() {
	lsc := os.Getenv("LS_COLORS")
	var eq bool
	var kb, ke int
	for i := 0; i < len(lsc); i++ {
		b := lsc[i]
		if b == '=' {
			ke = i
			eq = true
		} else if eq && b == ':' {
			if lsc[kb] == '*' {
				lsColorSuffix[lsc[kb+1:ke]] = lsc[ke+1 : i]
			} else {
				k := lsc[kb:ke]
				if in, ok := indicatorNamesMap[k]; ok {
					lsColorTypes[in] = lsc[ke+1 : i]
				} else {
					fmt.Printf("Unrecognized key: %s\n", k)
				}
			}
			kb = i + 1
			i += 2
			eq = false
		}
	}
}