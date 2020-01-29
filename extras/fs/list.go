package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/masteryconnect/pipe/message"
)

// List is a pipe/line producer/transformer to list the files/folders in a folder
type List struct {
	Root string

	Recursive    bool
	ShowHidden   bool
	IncludeDirs  bool
	ExcludeFiles bool
}

// T is the Tfunc for a pipe/line.
func (l List) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		l.run(message.String(m), out, errs)
	}
}

// P is the producer
func (l List) P(out chan<- interface{}, errs chan<- error) {
	l.run(l.Root, out, errs)
}

func (l List) run(root string, out chan<- interface{}, errs chan<- error) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// exclude hidden files and folders
		hidden, err := isHidden(info.Name())
		if !l.ShowHidden && err == nil && hidden {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			// don't output the root
			if path == root {
				return nil
			}

			// see if we should output dirs
			if l.IncludeDirs {
				out <- message.FileInfo{PathInfo: message.Path(path), FileInfo: info}
			}

			// see if we should go into the dir
			if !l.Recursive {
				return filepath.SkipDir
			}

			return nil
		}

		if !l.ExcludeFiles {
			out <- message.FileInfo{PathInfo: message.Path(path), FileInfo: info}
		}

		return nil
	})

	if err != nil {
		errs <- err
	}
}

func isHidden(filename string) (bool, error) {
	if runtime.GOOS != "windows" {
		// unix/linux file or directory that starts with . is hidden
		if len(filename) > 1 && filename[0:1] == "." {
			return true, nil
		}
	} else {
		return false, fmt.Errorf("unable to check if file is hidden under this OS")
	}
	return false, nil
}
