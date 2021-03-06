package main

import "syscall"

type fileMode uint64

type fileInfo struct {
	name     sufIndexed
	size     int64
	mode     fileMode
	time     int64
	linkok   bool
	linkname sufIndexed
	linkmode fileMode
}

// get info about file/directory name
func ls(name string) (fileList, error) {
	fi, err := stat(name)
	if err != nil {
		return nil, err
	}
	if fi.mode&syscall.S_IFMT == syscall.S_IFDIR {
		return readdir(name)
	}
	return []*fileInfo{fi}, nil
}

// stat returns a fileInfo describing the named file
func stat(name string) (*fileInfo, error) {
	var stat syscall.Stat_t
	err := syscall.Lstat(name, &stat)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}

	fi := &fileInfo{
		name:   *newSufIndexed(basename(name)),
		size:   int64(stat.Size),
		mode:   fileMode(stat.Mode),
		time:   gettime(&stat),
		linkok: true,
	}

	if fi.mode&syscall.S_IFMT == syscall.S_IFLNK {
		ln, err := readlink(name)
		if err != nil {
			fi.linkok = false
			return fi, nil
		}
		fi.linkname = *newSufIndexed(ln)
		err = syscall.Stat(name, &stat)
		if err != nil {
			fi.linkok = false
			return fi, nil
		}
		fi.linkmode = fileMode(stat.Mode)
	}

	return fi, nil
}

// basename the leading directory name from path name
func basename(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			return name[i+1:]
		}
	}
	return name
}

func cleanRight(path []byte) []byte {
	for i := len(path); i > 0; i-- {
		if path[i-1] != '/' {
			return path[:i]
		}
	}
	return path
}
