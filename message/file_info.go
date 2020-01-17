package message

import "os"

// FileInfo represents the info including path of a file or directory.
// This is to pull together our PathInfo and the os.FileInfo.
type FileInfo struct {
	os.FileInfo
	PathInfo
}

// String implements fmt.Stringer and includes a trailing '/' for directories
func (fi FileInfo) String() string {
	if fi.FileInfo.IsDir() {
		return fi.Path() + string(os.PathSeparator)
	}
	return fi.Path()
}
