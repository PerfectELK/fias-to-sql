package filehandler

import "os"

type File struct {
	path string
	file *os.File
}

func NewFile(path string) File {
	return File{
		path: path,
	}
}

func (f *File) Open(flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(f.path, flag, perm)
	if err != nil {
		return nil, err
	}
	f.file = file
	return file, err
}

func (f *File) Close() error {
	return f.file.Close()
}

// test
