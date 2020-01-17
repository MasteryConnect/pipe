package message

// PathInfo is the interface for providing the info for a path.
// This follows the os.FileInfo interface convention and is
// designed to bring the path and the FileInfo together.
type PathInfo interface {
	Path() string
}

// Path is a string that can be used as a PathInfo implementation
type Path string

// Path implements the PathInfo interface
func (fp Path) Path() string {
	return string(fp)
}
