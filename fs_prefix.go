package layerfs

import (
	"io/fs"
	"os"
	"strings"
	"time"
)

func Prefix(prefix string, fs fs.FS) fs.FS {
	prefix = strings.Trim(prefix, "/")
	return prefixFS{prefix, fs}
}

// prefixFS is a filesystem that prefixes all paths with a given string.
type prefixFS struct {
	prefix string
	fs     fs.FS
}

// Open opens the named file for reading.
func (f prefixFS) Open(name string) (fs.File, error) {
	name = strings.Trim(name, "/")
	if name == "." || name == f.prefix {
		fl := file{
			fileInfo: fileInfo{
				name:    name,
				mode:    os.FileMode(0755),
				modTime: time.Now(),
				isDir:   true,
			},
		}
		return fl, nil
	}
	ps := f.prefix + "/"

	if len(name) > len(f.prefix) {
		n := strings.TrimPrefix(name, ps)
		if len(n) == len(name) {
			return nil, os.ErrNotExist
		}
		return f.fs.Open(n)
	}
	n := strings.TrimPrefix(f.prefix, name+"/")
	if len(n) == len(f.prefix) {
		return nil, os.ErrNotExist
	}
	fl := file{
		fileInfo: fileInfo{
			name:    n,
			mode:    os.FileMode(0755),
			modTime: time.Now(),
			isDir:   true,
		},
	}
	return fl, nil
}

// ReadDir reads the named directory if the path is part of the prefix it will
// return a single entry with next suffix parts.
func (f prefixFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = strings.Trim(name, "/")
	if name == "." {
		p := strings.SplitN(f.prefix, "/", 2)
		// Return single entry with first path:",
		entries := []fs.DirEntry{
			dirEntry{fileInfo{
				name:    p[0],
				mode:    os.FileMode(0755),
				modTime: time.Now(),
				isDir:   true,
			}},
		}
		return entries, nil
	}
	if name == f.prefix {
		return fs.ReadDir(f.fs, ".")
	}
	if len(name) > len(f.prefix) {
		n := strings.TrimPrefix(name, f.prefix+"/")
		if len(n) == len(name) {
			return nil, os.ErrNotExist
		}
		return fs.ReadDir(f.fs, n)
	}

	n := strings.TrimPrefix(f.prefix, name+"/")
	if len(n) == len(f.prefix) {
		return nil, os.ErrNotExist
	}

	p := strings.SplitN(n, "/", 2)
	entries := []fs.DirEntry{
		dirEntry{fileInfo{
			name:    p[0],
			mode:    os.FileMode(0755),
			modTime: time.Now(),
			isDir:   true,
		}},
	}
	return entries, nil
}
